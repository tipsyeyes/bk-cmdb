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
	h.Set(common.BKHTTPAUTHORIZATION, header.Get(common.BKHTTPAUTHORIZATION))
	if err := ac.authClient.RegisterSystem(ctx, h, expectSystem); err != nil && err != ErrDuplicated {
		return err
	}

	// register system &project resource/action
	if err := ac.authClient.UpsertResourceTypeBatch(ctx, header, SystemIDCMDB, ScopeTypeIDSystem, expectSystemResourceType); err != nil {
		return err
	}
	//if err := ac.authClient.UpsertResourceTypeBatch(ctx, header, SystemIDCMDB, ScopeTypeIDBiz, expectBizResourceType); err != nil {
	//	return err
	//}

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

	return nil
}