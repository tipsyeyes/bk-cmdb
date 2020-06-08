package account

import (
	"configdatabase/src/auth/meta"
	"fmt"
)

// system constant
const (
	SystemIDCMDB   = "cc"
	SystemNameCMDB = "配置平台"
)

// ScopeTypeID constant
const (
	ScopeTypeIDSystem     = "system"
	ScopeTypeIDSystemName = "全局"

	ScopeTypeIDBiz     = "project"
	ScopeTypeIDBizName = "项目"
)

type AuthConfig struct {
	// blueking's auth center addresses
	Address []string
	// app code is used for authorize used.
	AppCode string
	// app secret is used for authorized
	AppSecret string
	// the system id that cmdb used in auth center.
	SystemID string
	// enable sync auth data to iam
	EnableSync          bool
	SyncWorkerCount     int
	SyncIntervalMinutes int
}

// Auth error
type AuthError struct {
	Reason	error
}

func (a *AuthError) Error() string {
	return fmt.Sprintf("err: %s", a.Reason.Error())
}

// BaseResp common result struct
type BaseResponse struct {
	Status  int    		`json:"status"`
	Code    string    	`json:"code"`
	Message string 		`json:"message"`
	RequestID string 	`json:"request_id"`
}

func (br BaseResponse) ErrorString() string {
	return fmt.Sprintf("request id: %s, error code: %s, message: %s", br.RequestID, br.Code, br.Message)
}

type System struct {
	SystemID   string `json:"system_id,omitempty"`
	SystemName string `json:"system_name"`
	Desc       string `json:"desc"`
	// 可为空，在使用注册资源的方式时
	QueryInterface string `json:"query_interface"`
	// 关联的资源所属，有业务、全局、项目等
	RelatedScopeTypes string `json:"related_scope_types"`

	// 管理者，可通过权限中心产品页面修改模型相关信息
	Managers string `json:"managers"`
	// 更新者，可为system
	Updater string `json:"updater,omitempty"`
	// 创建者，可为system
	Creator string `json:"creator,omitempty"`
}

// 注册资源类型信息
type ResourceType struct {
	ResourceTypeID       ResourceTypeID `json:"resource_type"`
	ResourceTypeName     string         `json:"resource_type_name"`
	ParentResourceTypeID ResourceTypeID `json:"parent_resource_type"`
	Share                bool           `json:"is_share"`
	Actions              []Action       `json:"actions"`
}

type Action struct {
	ActionID          	ActionID 	`json:"action_id"`
	ActionName        	string   	`json:"action_name"`
	// 是否为功能操作性权限
	IsFunctional		bool 		`json:"is_functional"`
	IsRelatedResource 	bool    	`json:"is_related_resource"`
}

// 范围信息 system or project
// 根据具体资源类型区分，是全局级，还是项目级
// scope_type: system/ project
// scope_id: cc/ instance id
type ScopeInfo struct {
	ScopeType string `json:"scope_type,omitempty"`
	ScopeID   string `json:"scope_id,omitempty"`
}

// 注册资源实体信息 RegisterEntityInfo
type RegisterInfo struct {
	// 创建者信息，可忽略
	// type: system/user
	// id: user/custom username
	CreatorType string           `json:"creator_type"`
	CreatorID   string           `json:"creator_id"`

	Resources   []ResourceEntity `json:"resources,omitempty"`
}

type ResourceEntity struct {
	ResourceType ResourceTypeID 	`json:"resource_type"`
	ResourceName string         	`json:"resource_name,omitempty"`
	// 资源id
	ResourceID   []RscTypeAndID 	`json:"resource_id,omitempty"`

	// TODO: remove
	ScopeInfo
}

// 资源 id &type
type RscTypeAndID struct {
	ResourceType ResourceTypeID `json:"resource_type"`
	ResourceID   string         `json:"resource_id,omitempty"`
}

type ResourceInfo struct {
	ResourceType ResourceTypeID `json:"resource_type"`
	// this filed is not always used, it's decided by the api
	// that is used.
	ResourceEntity
}

// iam授权资源
type IamResource []RscTypeAndID

type AuthorizedResource struct {
	ActionID     ActionID       `json:"action_id"`
	ResourceType ResourceTypeID `json:"resource_type"`
	ResourceIDs  []IamResource  `json:"resource_ids"`
}

// 用户组(即cmdb业务角色)成员信息
type UserGroupMembers struct {
	ID int64 `json:"group_id"`
	// user's group name, should be one of follows:
	// bk_biz_maintainer, bk_biz_productor, bk_biz_test, bk_biz_developer, operator
	Name  string   `json:"group_code"`
	Users []string `json:"users"`
}

// Search iam condition
type SearchCondition struct {
	ScopeInfo
	ResourceType    ResourceTypeID `json:"resource_type"`
	ParentResources []RscTypeAndID `json:"parent_resources"`
}

// iam page resource
type PageBackendResource struct {
	Count   int64                  `json:"count"`
	Results []meta.BackendResource `json:"results"`
}

type DeregisterInfo struct {
	Resources []ResourceEntity `json:"resources"`
}