package account

import (
	"configdatabase/src/auth/meta"
	"configdatabase/src/common"
	"configdatabase/src/common/util"
	"context"
	"fmt"
	"net/http"
	"strconv"
)

func (ac *AccountCenter) Init(ctx context.Context, header http.Header, configs meta.InitConfig) error {
	if err := ac.initAuthResources(ctx, header, configs); err != nil {
		return fmt.Errorf("initial auth resources failed, err: %v", err)
	}

	return nil
}

func (ac *AccountCenter) initAuthResources(ctx context.Context, header http.Header, configs meta.InitConfig) error {
	h := http.Header{}
	h.Set(common.BKHTTPCCRequestID, util.GetHTTPCCRequestID(header))
	h.Set(common.BKHTTPAuthorization, header.Get(common.BKHTTPAuthorization))
	if err := ac.authClient.RegisterSystem(ctx, h, expectSystem); err != nil && err != ErrDuplicated {
		return err
	}

	// register system &project resource/action
	if err := ac.authClient.UpsertResourceTypeBatch(ctx, header, SystemIDCMDB, ScopeTypeIDSystem, expectSystemResourceType); err != nil {
		return err
	}
	if err := ac.authClient.UpsertResourceTypeBatch(ctx, header, SystemIDCMDB, ScopeTypeIDBiz, expectBizResourceType); err != nil {
		return err
	}

	/// register resource instance
	// init business inst
	for _, biz := range configs.Bizs {
		bkBiz := RegisterInfo{
			CreatorID:   "system",
			CreatorType: "user",
			Resources: []ResourceEntity{
				{
					ResourceType: SysBusinessInstance,
					ResourceID: []RscTypeAndID{
						{ResourceType: SysBusinessInstance, ResourceID: strconv.FormatInt(biz.BizID, 10)},
					},
					ResourceName: biz.BizName,
					ScopeInfo: ScopeInfo{
						ScopeType: ScopeTypeIDSystem,
						ScopeID:   SystemIDCMDB,
					},
				},
			},
		}

		if err := ac.authClient.registerResource(ctx, header, &bkBiz); err != nil && err != ErrDuplicated {
			return err
		}
	}

	// add by elias 07/19
	/// register resource instance
	// init rbiz inst
	for _, rBiz := range configs.RBizs {
		bkRBiz := RegisterInfo{
			CreatorID:   "system",
			CreatorType: "user",
			Resources: []ResourceEntity{
				{
					ResourceType: BizRBizInstance,
					ResourceID: []RscTypeAndID{
						{ResourceType: BizRBizInstance, ResourceID: strconv.FormatInt(rBiz.InstanceID, 10)},
					},
					ResourceName: rBiz.Name,
					ScopeInfo: ScopeInfo{
						ScopeType: ScopeTypeIDBiz,
						ScopeID: strconv.FormatInt(rBiz.BizID, 10),
					},
				},
			},
		}

		if err := ac.authClient.registerResource(ctx, header, &bkRBiz); err != nil && err != ErrDuplicated {
			return err
		}
	}

	return nil
}