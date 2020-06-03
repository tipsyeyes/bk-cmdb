package account

import (
	"configdatabase/src/apimachinery/rest"
	"configdatabase/src/common/util"
	"context"
	"errors"
	"fmt"
	"net/http"
)

const (
	AuthSupplierAccountHeaderKey = "HTTP_BK_SUPPLIER_ACCOUNT"
)

const (
	codeDuplicated = "1901409"
	codeNotFound   = "1901404"
)

// Error define
var (
	ErrDuplicated = errors.New("Duplicated Item")
	ErrNotFound   = errors.New("Not Found")
)

type authClient struct {
	Config AuthConfig
	// http client instance
	client rest.ClientInterface
	// http header info
	basicHeader http.Header
}

func (a *authClient) RegisterSystem(ctx context.Context, header http.Header, system System) error {
	util.CopyHeader(a.basicHeader, header)
	const url = "/iam/perm-model/systems"
	resp := struct {
		BaseResponse
		Data System `json:"data"`
	}{}

	err := a.client.Post().
		SubResource(url).
		WithContext(ctx).
		WithHeaders(header).
		Body(system).
		Do().Into(&resp)
	if err != nil {
		return err
	}

	if resp.Status != statusSuccess {
		if resp.Code == codeDuplicated {
			return ErrDuplicated
		}
		return fmt.Errorf("regist system info for [%s] failed, err: %v", system.SystemID, resp.ErrorString())
	}

	return nil
}

