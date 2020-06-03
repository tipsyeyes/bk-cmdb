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

	return nil
}