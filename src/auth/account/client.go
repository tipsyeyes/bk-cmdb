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

// Register system to auth center
func (a *authClient) RegisterSystem(ctx context.Context, header http.Header, system System) error {
	util.CopyHeader(a.basicHeader, header)
	const url = "/iam/perm-model/system"
	resp := BaseResponse{}

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

// Create& update cmdb resource action
func (a *authClient) UpsertResourceTypeBatch(ctx context.Context, header http.Header, systemID, scopeType string, resources []ResourceType) error {
	util.CopyHeader(a.basicHeader, header)
	resp := BaseResponse{}

	err := a.client.Post().
		SubResourcef("/iam/perm-model/system/%s/scope-type/%s/resource-type/batch-upsert", systemID, scopeType).
		WithContext(ctx).
		WithHeaders(header).
		Body(struct {
			ResourceTypes []ResourceType `json:"resource_types"`
		}{resources}).
		Do().Into(&resp)
	if err != nil {
		return fmt.Errorf("regist resource %+v for [%s] failed, error: %v", resources, systemID, err)
	}
	if resp.Status != statusSuccess {
		return fmt.Errorf("regist resource %+v for [%s] failed, message: %s, code: %v", resources, systemID, resp.Message, resp.Code)
	}

	return nil
}

// Register resource instance
func (a *authClient) registerResource(ctx context.Context, header http.Header, info *RegisterInfo) error {
	// register resource with empty id will make crash
	for _, resource := range info.Resources {
		if resource.ResourceID == nil || len(resource.ResourceID) == 0 {
			return fmt.Errorf("resource id can't be empty, resource: %+v", resource)
		}
	}

	util.CopyHeader(a.basicHeader, header)
	resp := new(BaseResponse)
	err := a.client.Post().
		SubResourcef("/iam/perm/system/%s/resource/batch-register", a.Config.SystemID).
		WithContext(ctx).
		WithHeaders(header).
		Body(info).
		Do().Into(resp)

	if err != nil {
		return err
	}

	if resp.Status != statusSuccess {
		// 1901409 is for: resource already exist, can not created repeatedly
		if resp.Code == codeDuplicated {
			return ErrDuplicated
		}
		return &AuthError{Reason: fmt.Errorf("register resource failed, error code: %d, message: %s", resp.Code, resp.Message)}
	}

	return nil
}

func (a *authClient) deregisterResource(ctx context.Context, header http.Header, info *DeregisterInfo) error {
	util.CopyHeader(a.basicHeader, header)
	resp := new(BaseResponse)
	err := a.client.Delete().
		SubResourcef("/iam/perm/system/%s/resource/batch-delete", a.Config.SystemID).
		WithContext(ctx).
		WithHeaders(header).
		Body(info).
		Do().Into(resp)

	if err != nil {
		return err
	}

	if resp.Status != statusSuccess {
		if resp.Code == codeNotFound {
			return nil
		}
		return &AuthError{fmt.Errorf("deregister resource failed, error code: %d, message: %s", resp.Code, resp.Message)}
	}

	return nil
}

func (a *authClient) GetAuthorizedResources(ctx context.Context, header http.Header, body *ListAuthorizedResources) ([]AuthorizedResource, error) {
	util.CopyHeader(a.basicHeader, header)
	resp := ListAuthorizedResourcesResult{}

	err := a.client.Post().
		SubResourcef("/iam/perm/system/%s/authorized-resource/search", SystemIDCMDB).
		WithContext(ctx).
		WithHeaders(header).
		Body(body).
		Do().Into(&resp)
	if err != nil {
		return nil, fmt.Errorf("get authorized resource failed, err: %v", err)
	}
	if resp.Status != statusSuccess {
		return nil, fmt.Errorf("get authorized resource failed, err: %v", resp.ErrorString())
	}

	return resp.Data, nil
}

// 对特定资源批量鉴权
// 包含资源实例信息
func (a *authClient) verifyExactResourceBatch(ctx context.Context, header http.Header, batch *AuthBatch) ([]BatchStatus, error) {
	util.CopyHeader(a.basicHeader, header)
	resp := new(BatchResult)
	err := a.client.Post().
		SubResourcef("/iam/perm/systems/%s/resource-perm/batch-verify", a.Config.SystemID).
		WithContext(ctx).
		WithHeaders(header).
		Body(batch).
		Do().Into(resp)

	if err != nil {
		return nil, err
	}
	if resp.Status != statusSuccess {
		return nil, &AuthError{
			Reason:    fmt.Errorf("verify resource failed, error code: %d, message: %s", resp.Code, resp.Message),
		}
	}

	if len(batch.ResourceActions) != len(resp.Data) {
		return nil, fmt.Errorf("expect %d result, IAM returns %d result", len(batch.ResourceActions), len(resp.Data))
	}
	return resp.Data, nil
}

// 对任意资源批量鉴权
// 不包含资源实例信息，仅仅对功能权限
func (a *authClient) verifyAnyResourceBatch(ctx context.Context, header http.Header, batch *AuthBatch) ([]BatchStatus, error) {
	util.CopyHeader(a.basicHeader, header)
	resp := new(BatchResult)
	err := a.client.Post().
		// TODO 权限中心那边其实可以和 verifyExactResourceBatch用相同的接口
		SubResourcef("/iam/perm/systems/%s/any-resource-perm/batch-verify", a.Config.SystemID).
		WithContext(ctx).
		WithHeaders(header).
		Body(batch).
		Do().Into(resp)

	if err != nil {
		return nil, err
	}
	if resp.Status != statusSuccess {
		return nil, &AuthError{
			Reason:    fmt.Errorf("verify resource failed, error code: %d, message: %s", resp.Code, resp.Message),
		}
	}

	if len(batch.ResourceActions) != len(resp.Data) {
		return nil, fmt.Errorf("expect %d result, IAM returns %d result", len(batch.ResourceActions), len(resp.Data))
	}
	return resp.Data, nil
}