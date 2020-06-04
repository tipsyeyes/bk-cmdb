package account

// ResourceTypeID is resource's type in auth center.
type ResourceTypeID string

// ActionID is action's type in auth center.
type ActionID string

// 全局级资源类型
const (
	SysBusinessInstance   ResourceTypeID = "sys_project_instance"
)

const (
	SysSystemBase         ResourceTypeID = "sys_system_base"
	SysHostInstance       ResourceTypeID = "sys_host_instance"
	SysEventPushing       ResourceTypeID = "sys_event_pushing"
	SysModelGroup         ResourceTypeID = "sys_model_group"
	SysModel              ResourceTypeID = "sys_model"
	SysInstance           ResourceTypeID = "sys_instance"
	SysAssociationType    ResourceTypeID = "sys_association_type"
	SysAuditLog           ResourceTypeID = "sys_audit_log"
	SysOperationStatistic ResourceTypeID = "sys_operation_statistic"
)

// 项目级资源类型
const (
	BizHostInstance           ResourceTypeID = "proj_host_instance"
)

const (
	BizCustomQuery            ResourceTypeID = "proj_custom_query"
	BizProcessInstance        ResourceTypeID = "proj_process_instance"
	BizTopology               ResourceTypeID = "proj_topology"
	BizModelGroup             ResourceTypeID = "proj_model_group"
	BizModel                  ResourceTypeID = "proj_model"
	BizInstance               ResourceTypeID = "proj_instance"
	BizAuditLog               ResourceTypeID = "proj_audit_log"
	BizProcessServiceTemplate ResourceTypeID = "proj_process_service_template"
	BizProcessServiceCategory ResourceTypeID = "proj_process_service_category"
	BizProcessServiceInstance ResourceTypeID = "proj_process_service_instance"
	BizSetTemplate            ResourceTypeID = "proj_set_template"
)

var ResourceTypeIDMap = map[ResourceTypeID]string{
	SysSystemBase:       "系统基础",
	SysBusinessInstance: "项目",
	SysHostInstance:     "主机",
	SysEventPushing:     "事件推送",
	SysModelGroup:       "模型分组",
	SysModel:            "模型",
	SysInstance:         "实例",
	SysAssociationType:  "关联类型",
	SysAuditLog:         "操作审计",
	BizCustomQuery:      "动态分组",
	BizHostInstance:     "项目主机",
	BizProcessInstance:  "进程",
	// TODO: delete this when upgrade to v3.5.x
	BizTopology:               "项目拓扑",
	BizModelGroup:             "模型分组",
	BizModel:                  "模型",
	BizInstance:               "实例",
	BizAuditLog:               "操作审计",

	BizProcessServiceTemplate: "服务模板",
	BizProcessServiceCategory: "服务分类",
	BizProcessServiceInstance: "服务实例",
	BizSetTemplate:            "业务模板",
	SysOperationStatistic:     "运营统计",
}

// ActionID define
const (
	// Unknown action is a action that can not be recognized by the auth center.
	Unknown ActionID = "unknown"
	Get     ActionID = "get"
	Create  ActionID = "create"
	Edit    ActionID = "edit"
	Delete  ActionID = "delete"
	// Archive for business
	Archive ActionID = "archive"
)

const (
	// business model topology operation.
	ModelTopologyOperation ActionID = "model_topology_operation"
)

var ActionIDNameMap = map[ActionID]string{
	Unknown:                "未知操作",
	Get:                    "查询",
	Create:                 "新建",
	Edit:                   "编辑",
	Delete:                 "删除",
	Archive:                "归档",
	ModelTopologyOperation: "编辑项目层级",
}








