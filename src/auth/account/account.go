package account

import (
	"configdatabase/src/auth/account/permit"
	"configdatabase/src/common"
	commonutil "configdatabase/src/common/util"
	"configdatabase/src/apimachinery/flowctrl"
	"configdatabase/src/apimachinery/rest"
	"configdatabase/src/apimachinery/util"
	"configdatabase/src/auth/meta"
	"configdatabase/src/common/auth"
	"configdatabase/src/common/blog"
	"configdatabase/src/common/metadata"
	"context"
	"errors"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"strconv"
	"sync"
)

const (
	statusSuccess = 0
	authApiPrefix = "/api/v1"
)

const (
	authAppCodeHeaderKey   string = "X-CC-APP-CODE"
	authAppSecretHeaderKey string = "X-CC-APP-SECRET"
	cmdbUser               string = "user"
	cmdbUserID             string = "system"
)

type acDiscovery struct {
	// auth's servers address, must prefixed with http:// or https://
	servers []string
	index   int
	sync.Mutex
}

func (s *acDiscovery) GetServers() ([]string, error) {
	s.Lock()
	defer s.Unlock()

	num := len(s.servers)
	if num == 0 {
		return []string{}, errors.New("oops, there is no auth center server can be used")
	}

	if s.index < num-1 {
		s.index = s.index + 1
		return append(s.servers[s.index-1:], s.servers[:s.index-1]...), nil
	} else {
		s.index = 0
		return append(s.servers[num-1:], s.servers[:num-1]...), nil
	}
}

// NewAccountCenter create a instance to handle resources with AuthCenter.
func NewAccountCenter(tls *util.TLSClientConfig, cfg AuthConfig, reg prometheus.Registerer) (*AccountCenter, error) {
	blog.V(5).Infof("new auth center client with parameters tls: %+v, cfg: %+v", tls, cfg)
	if !auth.IsAuthed() {
		return new(AccountCenter), nil
	}
	client, err := util.NewClient(tls)
	if err != nil {
		return nil, err
	}

	c := &util.Capability{
		Client: client,
		Discover: &acDiscovery{
			servers: cfg.Address,
		},
		Throttle: flowctrl.NewRateLimiter(1000, 1000),
		Mock: util.MockInfo{
			Mocked: false,
		},
		Reg: reg,
	}

	header := http.Header{}
	header.Set("Content-Type", "application/json")
	header.Set("Accept", "application/json")
	header.Set(authAppCodeHeaderKey, cfg.AppCode)
	header.Set(authAppSecretHeaderKey, cfg.AppSecret)

	return &AccountCenter{
		Config: cfg,
		authClient: &authClient{
			client:      rest.NewRESTClient(c, authApiPrefix),
			Config:      cfg,
			basicHeader: header,
		},
		authFilter: &authFilter{},
	}, nil
}

type AccountCenter struct {
	Config AuthConfig
	// http client instance
	client rest.ClientInterface
	// http header info
	header     http.Header
	authClient *authClient

	// account center filter
	authFilter *authFilter
}

func (ac *AccountCenter) Enabled() bool {
	return auth.IsAuthed()
}

func (ac *AccountCenter) Authorize(ctx context.Context, a *meta.AuthAttribute) (decision meta.Decision, err error) {
	if !auth.IsAuthed() {
		return meta.Decision{Authorized: true}, nil
	} else if a == nil {
		return meta.Decision{Authorized: true}, nil
	}

	// super user return true
	if commonutil.ExtractRequestSuperFromContext(ctx) {
		return meta.Decision{Authorized: true}, nil
	}

	// filter out SkipAction, which set by api server to skip authorization
	noSkipResources := make([]meta.ResourceAttribute, 0)
	for _, resource := range a.Resources {
		if resource.Action == meta.SkipAction {
			continue
		}
		noSkipResources = append(noSkipResources, resource)
	}
	a.Resources = noSkipResources
	if len(noSkipResources) == 0 {
		blog.V(5).Infof("Authorize skip. auth attribute: %+v", a)
		return meta.Decision{Authorized: true}, nil
	}

	batchResult, err := ac.AuthorizeBatchEx(ctx, a.User, a.Resources...)
	if err != nil {
		blog.Errorf("AuthorizeBatch error. err:%s", err.Error())
		return meta.Decision{}, err
	}
	noAuth := make([]string, 0)
	for i, item := range batchResult {
		if !item.Authorized {
			noAuth = append(noAuth, fmt.Sprintf("resource [%v] permission deny by reason: %s", a.Resources[i].Type, item.Reason))
		}
	}

	if len(noAuth) > 0 {
		return meta.Decision{
			Authorized: false,
			Reason:     fmt.Sprintf("%v", noAuth),
		}, nil
	}

	return meta.Decision{Authorized: true}, nil
}

func (ac *AccountCenter) AuthorizeBatchEx(ctx context.Context, user meta.UserInfo, resources ...meta.ResourceAttribute) (decisions []meta.Decision, err error) {
	// api层不做鉴权处理
	return nil, nil
}

// AuthorizeBatch api层批量鉴权
// 包括各种功能权限，比如"业务的创建(cc.business.create)"
// 或者对某个资源的修改或删除权限，比如"业务的修改(cc.business.delete)"，附带鉴权的资源id信息
func (ac *AccountCenter) AuthorizeBatch(ctx context.Context, user meta.UserInfo, resources ...meta.ResourceAttribute) (decisions []meta.Decision, err error) {
	rid := commonutil.ExtractRequestIDFromContext(ctx)
	token := commonutil.ExtractRequestTokenFromContext(ctx)
	decisions = make([]meta.Decision, len(resources))
	if !auth.IsAuthed() {
		for i := range decisions {
			decisions[i].Authorized = true
		}
		return decisions, nil
	}

	header := http.Header{}
	header.Set(AuthSupplierAccountHeaderKey, user.SupplierAccount)
	header.Set(common.BKHTTPCCRequestID, rid)
	header.Set(common.BKHTTPAuthorization, token)

	// this two index array record the resources's action original index.
	// used for recover the order of decisions.
	sysInputIndexes := make([]int, 0)
	sysInputExactIndexes := make([]int, 0)
	bizInputIndexes := make(map[int64][]int)
	bizInputExactIndexes := make(map[int64][]int)

	sysInput := AuthBatch{
		Principal: Principal{
			Type: cmdbUser,
			ID:   user.UserName,
		},
		ScopeInfo: ScopeInfo{
			ScopeType: ScopeTypeIDSystem,
			ScopeID:   SystemIDCMDB,
		},
		ResourceActions: make([]ResourceAction, 0),
	}
	sysExactInput := sysInput

	businessesInputs := make(map[int64]AuthBatch)
	businessesExactInputs := make(map[int64]AuthBatch)
	for index, rsc := range resources {
		action, err := AdaptorAction(&rsc)
		if err != nil {
			blog.Errorf("auth batch, but adaptor action:%s failed, err: %v, rid: %s", rsc.Action, err, rid)
			return nil, err
		}

		// pick out skip resource at first.
		if permit.ShouldSkipAuthorize(&rsc) {
			// this resource should be skipped, do not need to verify in auth center.
			decisions[index].Authorized = true
			blog.V(5).Infof("skip authorization for resource: %+v, rid: %s", rsc, rid)
			continue
		}

		info, err := adaptor(&rsc)
		if err != nil {
			blog.Errorf("auth batch, but adaptor resource type:%s failed, err: %v, rid: %s", rsc.Basic.Type, err, rid)
			return nil, err
		}

		// modify special resource
		if rsc.Type == meta.MainlineModel || rsc.Type == meta.ModelTopology {
			blog.V(4).Infof("force convert scope type to global for resource type: %s, rid: %s", rsc.Type, rid)
			rsc.BusinessID = 0
		}

		if rsc.BusinessID > 0 {
			// this is a business resource.
			var tmpInputs map[int64]AuthBatch
			var tmpIndexes map[int64][]int
			if len(info.ResourceID) > 0 {
				tmpInputs = businessesExactInputs
				tmpIndexes = bizInputExactIndexes
			} else {
				tmpInputs = businessesInputs
				tmpIndexes = bizInputIndexes

			}

			if _, exist := tmpInputs[rsc.BusinessID]; !exist {
				tmpInputs[rsc.BusinessID] = AuthBatch{
					Principal: Principal{
						Type: cmdbUser,
						ID:   user.UserName,
					},
					ScopeInfo: ScopeInfo{
						ScopeType: ScopeTypeIDBiz,
						ScopeID:   strconv.FormatInt(rsc.BusinessID, 10),
					},
				}
				// initialize the business input indexes.
				tmpIndexes[rsc.BusinessID] = make([]int, 0)
			}

			a := tmpInputs[rsc.BusinessID]
			a.ResourceActions = append(a.ResourceActions, ResourceAction{
				ActionID:     action,
				ResourceType: info.ResourceType,
				ResourceID:   info.ResourceID,
			})
			tmpInputs[rsc.BusinessID] = a

			// record it's resource index
			indexes := tmpIndexes[rsc.BusinessID]
			indexes = append(indexes, index)
			tmpIndexes[rsc.BusinessID] = indexes
		} else {

			if len(info.ResourceID) > 0 {
				sysExactInput.ResourceActions = append(sysExactInput.ResourceActions, ResourceAction{
					ActionID:     action,
					ResourceType: info.ResourceType,
					ResourceID:   info.ResourceID,
				})

				// record it's system resource index
				sysInputExactIndexes = append(sysInputExactIndexes, index)
			} else {
				sysInput.ResourceActions = append(sysInput.ResourceActions, ResourceAction{
					ActionID:     action,
					ResourceType: info.ResourceType,
				})

				// record it's system resource index
				sysInputIndexes = append(sysInputIndexes, index)
			}

		}
	}

	// it's time to get the auth status from auth center now.
	// get biz resource auth status at first.
	// any business inputs
	for biz, rsc := range businessesInputs {
		// if resourceType that not related to resourceID, clear resourceID field
		for idx, resourceAction := range rsc.ResourceActions {
			if IsRelatedToResourceID(resourceAction.ResourceType) == false {
				rsc.ResourceActions[idx].ResourceID = make([]RscTypeAndID, 0)
			}
		}
		statuses, err := ac.authClient.verifyAnyResourceBatch(ctx, header, &rsc)
		if err != nil {
			return nil, fmt.Errorf("get any resource[%s/%s] auth status failed, err: %v", rsc.ScopeType, rsc.ScopeID, err)
		}

		if len(statuses) != len(rsc.ResourceActions) {
			return nil, fmt.Errorf("got mismatch any biz authorize response from auth center, want: %d, got: %d", len(rsc.ResourceActions), len(statuses))
		}

		// update the decisions
		for index, status := range statuses {
			if rsc.ResourceActions[index].ResourceType != status.ResourceType ||
				string(rsc.ResourceActions[index].ActionID) != status.ActionID {
				return nil, fmt.Errorf("got any business auth mismatch info from auth center, with resource type[%s:%s], action[%s:%s]",
					rsc.ResourceActions[index].ResourceType, status.ResourceType, rsc.ResourceActions[index].ActionID, status.ActionID)
			}
			decisions[bizInputIndexes[biz][index]].Authorized = status.IsPass
		}
	}

	// exact business inputs
	for biz, rsc := range businessesExactInputs {
		// if resourceType that not related to resourceID, clear resourceID field
		for idx, resourceAction := range rsc.ResourceActions {
			if IsRelatedToResourceID(resourceAction.ResourceType) == false {
				rsc.ResourceActions[idx].ResourceID = make([]RscTypeAndID, 0)
			}
		}
		statuses, err := ac.authClient.verifyExactResourceBatch(ctx, header, &rsc)
		if err != nil {
			return nil, fmt.Errorf("get exact resource[%s/%s] auth status failed, err: %v", rsc.ScopeType, rsc.ScopeID, err)
		}

		if len(statuses) != len(rsc.ResourceActions) {
			return nil, fmt.Errorf("got mismatch exact biz authorize response from auth center, want: %d, got: %d", len(rsc.ResourceActions), len(statuses))
		}

		// update the decisions
		for index, status := range statuses {
			if rsc.ResourceActions[index].ResourceType != status.ResourceType ||
				string(rsc.ResourceActions[index].ActionID) != status.ActionID {
				return nil, fmt.Errorf("got exact business auth mismatch info from auth center, with resource type[%s:%s], action[%s:%s]",
					rsc.ResourceActions[index].ResourceType, status.ResourceType, rsc.ResourceActions[index].ActionID, status.ActionID)
			}
			decisions[bizInputExactIndexes[biz][index]].Authorized = status.IsPass
		}
	}

	if len(sysInput.ResourceActions) != 0 {
		// if resourceType that not related to resourceID, clear resourceID field
		for idx, resourceAction := range sysInput.ResourceActions {
			if IsRelatedToResourceID(resourceAction.ResourceType) == false {
				sysInput.ResourceActions[idx].ResourceID = make([]RscTypeAndID, 0)
			}
		}
		// get system resource auth status secondly.
		statuses, err := ac.authClient.verifyAnyResourceBatch(ctx, header, &sysInput)
		if err != nil {
			return nil, fmt.Errorf("get any system resource[%s/%s] auth status failed, err: %v", sysInput.ScopeType, sysInput.ScopeID, err)
		}

		if len(statuses) != len(sysInput.ResourceActions) {
			return nil, fmt.Errorf("got mismatch any system authorize response from auth center, want: %d, got: %d", len(sysInput.ResourceActions), len(statuses))
		}

		// update the system auth decisions
		for index, status := range statuses {
			if sysInput.ResourceActions[index].ResourceType != status.ResourceType ||
				string(sysInput.ResourceActions[index].ActionID) != status.ActionID {
				return nil, fmt.Errorf("got any system auth mismatch info from auth center, with resource type[%s:%s], action[%s:%s]",
					sysInput.ResourceActions[index].ResourceType, status.ResourceType,
					sysInput.ResourceActions[index].ActionID, status.ActionID)
			}
			decisions[sysInputIndexes[index]].Authorized = status.IsPass
		}
	}

	if len(sysExactInput.ResourceActions) != 0 {
		// if resourceType that not related to resourceID, clear resourceID field
		for idx, resourceAction := range sysExactInput.ResourceActions {
			if IsRelatedToResourceID(resourceAction.ResourceType) == false {
				sysExactInput.ResourceActions[idx].ResourceID = make([]RscTypeAndID, 0)
			}
		}
		// get system resource auth status secondly.
		statuses, err := ac.authClient.verifyExactResourceBatch(ctx, header, &sysExactInput)
		if err != nil {
			return nil, fmt.Errorf("get exact system resource[%s/%s] auth status failed, err: %v", sysInput.ScopeType, sysInput.ScopeID, err)
		}

		if len(statuses) != len(sysExactInput.ResourceActions) {
			return nil, fmt.Errorf("got mismatch exact authorize response from auth center, want: %d, got: %d", len(sysExactInput.ResourceActions), len(statuses))
		}

		// update the system auth decisions
		for index, status := range statuses {
			if sysExactInput.ResourceActions[index].ResourceType != status.ResourceType ||
				string(sysExactInput.ResourceActions[index].ActionID) != status.ActionID {
				return nil, fmt.Errorf("got exact system auth mismatch info from auth center, with resource type[%s:%s], action[%s:%s]",
					sysExactInput.ResourceActions[index].ResourceType, status.ResourceType,
					sysExactInput.ResourceActions[index].ActionID, status.ActionID)
			}
			decisions[sysInputExactIndexes[index]].Authorized = status.IsPass
		}
	}

	return decisions, nil
}

func (ac *AccountCenter) GetAnyAuthorizedBusinessList(ctx context.Context, user meta.UserInfo) ([]int64, error) {
	return ac.GetExactAuthorizedBusinessList(ctx, user)
}

// GetExactAuthorizedBusinessList
// get a user's authorized read business list.
// 获取所有授权的业务列表
// ctx := util.NewContextFromHTTPHeader(header)
func (ac *AccountCenter) GetExactAuthorizedBusinessList(ctx context.Context, user meta.UserInfo) ([]int64, error) {
	rid := commonutil.ExtractRequestIDFromContext(ctx)
	token := commonutil.ExtractRequestTokenFromContext(ctx)

	if !auth.IsAuthed() {
		return make([]int64, 0), nil
	}

	option := &ListAuthorizedResources{
		Principal: Principal{
			Type: cmdbUser,
			ID:   user.UserName,
		},
		ScopeInfo: ScopeInfo{
			ScopeType: ScopeTypeIDSystem,
			ScopeID:   SystemIDCMDB,
		},
		TypeActions: []TypeAction{
			{
				ActionID:     Get,
				ResourceType: SysBusinessInstance,
			},
		},
		DataType: "array",
		Exact:    true,
	}

	header := http.Header{}
	header.Set(common.BKHTTPCCRequestID, rid)
	header.Set(common.BKHTTPAuthorization, token)
	appListRsc, err := ac.authClient.GetAuthorizedResources(ctx, header, option)
	if err != nil {
		return nil, err
	}

	businessIDs := make([]int64, 0)
	for _, appRsc := range appListRsc {
		for _, app := range appRsc.ResourceIDs {
			id, err := strconv.ParseInt(app.ResourceID, 10, 64)
			if err != nil {
				return businessIDs, err
			}
			businessIDs = append(businessIDs, id)
		}
	}

	return businessIDs, nil
}

func (ac *AccountCenter) ListAuthorizedResources(ctx context.Context, username string, bizID int64, resourceType meta.ResourceType, action meta.Action) ([]IamResource, error) {
	return nil, nil
}

func (ac *AccountCenter) AdminEntrance(ctx context.Context, user meta.UserInfo) ([]string, error) {
	return nil, nil
}

func (ac *AccountCenter) GetAuthorizedAuditList(ctx context.Context, user meta.UserInfo, businessID int64) ([]AuthorizedResource, error) {
	return nil, nil
}

func (ac *AccountCenter) GetNoAuthSkipUrl(ctx context.Context, header http.Header, p []metadata.Permission) (url string, err error) {
	return "", nil
}

func (ac *AccountCenter) GetUserGroupMembers(ctx context.Context, header http.Header, bizID int64, groups []string) ([]UserGroupMembers, error) {
	return nil, nil
}

const pageSize = 500

// Register a resource instance
func (ac *AccountCenter) RegisterResource(ctx context.Context, rs ...meta.ResourceAttribute) error {
	rid := commonutil.ExtractRequestIDFromContext(ctx)
	token := commonutil.ExtractRequestTokenFromContext(ctx)

	if !auth.IsAuthed() {
		blog.V(5).Infof("auth disabled, auth config: %+v, rid: %s", ac.Config, rid)
		return nil
	}

	if len(rs) == 0 {
		return errors.New("no resource to be registered")
	}

	registerInfo, err := ac.DryRunRegisterResource(ctx, rs...)
	if err != nil {
		return err
	}

	// 清除不需要关联资源ID类型的注册
	resourceEntities := make([]ResourceEntity, 0)
	for index, resource := range registerInfo.Resources {
		if IsRelatedToResourceID(resource.ResourceType) == true {
			resourceEntities = append(resourceEntities, registerInfo.Resources[index])
		}
	}
	if len(resourceEntities) == 0 {
		return nil
	}
	registerInfo.Resources = resourceEntities

	header := http.Header{}
	//header.Set(AuthSupplierAccountHeaderKey, rs[0].SupplierAccount)
	header.Set(common.BKHTTPCCRequestID, rid)
	header.Set(common.BKHTTPAuthorization, token)

	var firstErr error
	count := len(resourceEntities)
	for start := 0; start < count; start += pageSize {
		end := start + pageSize
		if end > count {
			end = count
		}
		entities := resourceEntities[start:end]
		registerInfo.Resources = entities
		if err := ac.authClient.registerResource(ctx, header, registerInfo); err != nil {
			if err != ErrDuplicated {
				firstErr = err
			}
		}
	}

	return firstErr
}

// Register a resource instance
func (ac *AccountCenter) DryRunRegisterResource(ctx context.Context, rs ...meta.ResourceAttribute) (*RegisterInfo, error) {
	rid := commonutil.ExtractRequestIDFromContext(ctx)
	user := commonutil.ExtractRequestUserFromContext(ctx)
	if len(user) == 0 {
		user = cmdbUserID
	}

	if !auth.IsAuthed() {
		blog.V(5).Infof("auth disabled, auth config: %+v, rid: %s", ac.Config, rid)
		return new(RegisterInfo), nil
	}

	info := RegisterInfo{}
	info.CreatorType = cmdbUser
	info.CreatorID = user
	info.Resources = make([]ResourceEntity, 0)
	for _, r := range rs {
		if !ac.authFilter.needRegisterResource(ctx, r) {
			continue
		}

		if len(r.Basic.Type) == 0 {
			return nil, errors.New("invalid resource attribute with empty object")
		}
		scope, err := ac.getScopeInfo(&r)
		if err != nil {
			return nil, err
		}

		rscInfo, err := adaptor(&r)
		if err != nil {
			return nil, fmt.Errorf("adaptor resource info failed, err: %v", err)
		}
		entity := ResourceEntity{
			ResourceType: rscInfo.ResourceType,
			ScopeInfo: ScopeInfo{
				ScopeType: scope.ScopeType,
				ScopeID:   scope.ScopeID,
			},
			ResourceName: rscInfo.ResourceName,
			ResourceID:   rscInfo.ResourceID,
		}

		// TODO replace register with batch create or update interface, currently is register one by one.
		info.Resources = append(info.Resources, entity)
	}
	return &info, nil
}

// Deregister a resource instance
func (ac *AccountCenter) DeregisterResource(ctx context.Context, rs ...meta.ResourceAttribute) error {
	rid := commonutil.ExtractRequestIDFromContext(ctx)
	token := commonutil.ExtractRequestTokenFromContext(ctx)

	if !auth.IsAuthed() {
		return nil
	}

	if len(rs) <= 0 {
		// not resource should be deregister
		return nil
	}
	info := DeregisterInfo{}

	header := http.Header{}
	header.Set(common.BKHTTPCCRequestID, rid)
	header.Set(common.BKHTTPAuthorization, token)
	for _, r := range rs {
		if !ac.authFilter.needRegisterResource(ctx, r) {
			continue
		}

		if len(r.Basic.Type) == 0 {
			return errors.New("invalid resource attribute with empty object")
		}

		scope, err := ac.getScopeInfo(&r)
		if err != nil {
			return err
		}

		rscInfo, err := adaptor(&r)
		if err != nil {
			return fmt.Errorf("adaptor resource info failed, err: %v", err)
		}

		entity := ResourceEntity{}
		entity.ScopeID = scope.ScopeID
		entity.ScopeType = scope.ScopeType
		entity.ResourceType = rscInfo.ResourceType
		entity.ResourceID = rscInfo.ResourceID
		entity.ResourceName = rscInfo.ResourceName

		// 不关联实例ID的资源类型不需要取消注册
		if IsRelatedToResourceID(entity.ResourceType) == false {
			continue
		}

		info.Resources = append(info.Resources, entity)

		header.Set(AuthSupplierAccountHeaderKey, r.SupplierAccount)
	}

	if len(info.Resources) == 0 {
		if blog.V(5) {
			blog.InfoJSON("no resource to be deregister for original resource: %s, rid: %s", rs, rid)
		}
		return nil
	}

	return ac.authClient.deregisterResource(ctx, header, &info)
}

// Deregister a resource instance with raw iam resource id
func (ac *AccountCenter) RawDeregisterResource(ctx context.Context, scope ScopeInfo, rs ...meta.BackendResource) error {
	return nil
}

// Update a resource instance info
func (ac *AccountCenter) UpdateResource(ctx context.Context, r *meta.ResourceAttribute) error {
	return nil
}

// Get a resource instance info
func (ac *AccountCenter) Get(ctx context.Context) error {
	return nil
}

// List iam resources instance by condition/convert level
func (ac *AccountCenter) ListResources(ctx context.Context, r *meta.ResourceAttribute) ([]meta.BackendResource, error) {
	return nil, nil
}
func (ac *AccountCenter) RawListResources(ctx context.Context, header http.Header, searchCondition SearchCondition) ([]meta.BackendResource, error) {
	return nil, nil
}

// List iam resources in page
func (ac *AccountCenter) ListPageResources(ctx context.Context, r *meta.ResourceAttribute, limit, offset int64) (PageBackendResource, error) {
	return PageBackendResource{}, nil
}
func (ac *AccountCenter) RawPageListResources(ctx context.Context, header http.Header, searchCondition SearchCondition, limit, offset int64) (PageBackendResource, error) {
	return PageBackendResource{}, nil
}

func (ac *AccountCenter) getScopeInfo(r *meta.ResourceAttribute) (*ScopeInfo, error) {
	s := new(ScopeInfo)
	// TODO: this operation may be wrong, because some api filters does not
	// fill the business id field, so these api should be normalized.
	if r.BusinessID > 0 {
		s.ScopeType = ScopeTypeIDBiz
		s.ScopeID = strconv.FormatInt(r.BusinessID, 10)
	} else {
		s.ScopeType = ScopeTypeIDSystem
		s.ScopeID = SystemIDCMDB
	}
	return s, nil
}