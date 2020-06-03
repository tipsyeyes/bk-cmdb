package account

import "fmt"

// system constant
const (
	SystemIDCMDB   = "cc"
	SystemNameCMDB = "配置平台"
)

// ScopeTypeID constant
const (
	ScopeTypeIDSystem     = "system"
	ScopeTypeIDSystemName = "全局"

	ScopeTypeIDBiz     = "proj"
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
	Functional			bool 		`json:"is_functional"`
	IsRelatedResource 	bool    	`json:"is_related_resource"`
}

