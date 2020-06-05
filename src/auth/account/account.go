package account

import (
	"configdatabase/src/apimachinery/flowctrl"
	"configdatabase/src/apimachinery/rest"
	"configdatabase/src/apimachinery/util"
	"configdatabase/src/auth/meta"
	"configdatabase/src/common/auth"
	"configdatabase/src/common/blog"
	"configdatabase/src/common/metadata"
	"context"
	"errors"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
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
	}, nil
}

type AccountCenter struct {
	Config AuthConfig
	// http client instance
	client rest.ClientInterface
	// http header info
	header     http.Header
	authClient *authClient
}

func (ac *AccountCenter) Enabled() bool {
	return auth.IsAuthed()
}

func (ac *AccountCenter) Authorize(ctx context.Context, a *meta.AuthAttribute) (decision meta.Decision, err error) {
	return meta.Decision{Authorized: true}, nil
}

func (ac *AccountCenter) AuthorizeBatch(ctx context.Context, user meta.UserInfo, resources ...meta.ResourceAttribute) (decisions []meta.Decision, err error) {
	return nil, nil
}

func (ac *AccountCenter) GetAnyAuthorizedBusinessList(ctx context.Context, user meta.UserInfo) ([]int64, error) {
	return nil, nil
}

func (ac *AccountCenter) GetExactAuthorizedBusinessList(ctx context.Context, user meta.UserInfo) ([]int64, error) {
	return nil, nil
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

// Register a resource instance
func (ac *AccountCenter) RegisterResource(ctx context.Context, rs ...meta.ResourceAttribute) error {
	return nil
}

// Register a resource instance
func (ac *AccountCenter) DryRunRegisterResource(ctx context.Context, rs ...meta.ResourceAttribute) (*RegisterInfo, error) {
	return nil, nil
}

// Deregister a resource instance
func (ac *AccountCenter) DeregisterResource(ctx context.Context, rs ...meta.ResourceAttribute) error {
	return nil
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