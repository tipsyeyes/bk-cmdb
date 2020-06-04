package account

import (
	"configdatabase/src/auth/meta"
	"context"
	"fmt"
	"net/http"
)

func (ac *AccountCenter) Init(ctx context.Context, configs meta.InitConfig) error {
	if err := ac.initAuthResources(ctx, configs); err != nil {
		return fmt.Errorf("initial auth resources failed, err: %v", err)
	}

	return nil
}

func (ac *AccountCenter) initAuthResources(ctx context.Context, configs meta.InitConfig) error {
	header := http.Header{}
	if err := ac.authClient.RegisterSystem(ctx, header, expectSystem); err != nil && err != ErrDuplicated {
		return err
	}

	// register system &project resource/action
	if err := ac.authClient.UpsertResourceTypeBatch(ctx, header, SystemIDCMDB, ScopeTypeIDSystem, expectSystemResourceType); err != nil {
		return err
	}
	if err := ac.authClient.UpsertResourceTypeBatch(ctx, header, SystemIDCMDB, ScopeTypeIDBiz, expectBizResourceType); err != nil {
		return err
	}

	// register resource instance


	return nil
}