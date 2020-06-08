package account

import (
	"configdatabase/src/auth/meta"
	"fmt"
)

var NotEnoughLayer = fmt.Errorf("not enough layer")

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

const (
	UserCustom ResourceTypeID = "userCustom"
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
	// host action
	ModuleTransfer ActionID = "module_transfer"
	// business topology action
	HostTransfer ActionID = "host_transfer"
	// system base action, related to model topology
	ModelTopologyView ActionID = "model_topology_view"
	// business model topology operation.
	ModelTopologyOperation ActionID = "model_topology_operation"
	// assign host(s) to a business
	// located system/host/assignHostsToBusiness in auth center.
	AssignHostsToBusiness ActionID = "assign_hosts_to_business"
	BindModule            ActionID = "bind_module"
	AdminEntrance         ActionID = "admin_entrance"
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

func AdaptorAction(r *meta.ResourceAttribute) (ActionID, error) {
	if r.Basic.Type == meta.ModelAttributeGroup ||
		r.Basic.Type == meta.ModelUnique ||
		r.Basic.Type == meta.ModelAttribute {
		if r.Action == meta.Delete || r.Action == meta.Update || r.Action == meta.Create {
			return Edit, nil
		}
	}

	if r.Basic.Type == meta.Business {
		if r.Action == meta.Archive {
			return Archive, nil
		}

		// edit a business.
		if r.Action == meta.Create {
			return Create, nil
		}

		if r.Action == meta.Update {
			return Edit, nil
		}
	}

	// TODO: confirm this.
	if r.Action == meta.Execute && r.Basic.Type == meta.DynamicGrouping {
		return Get, nil
	}

	if r.Action == meta.Find || r.Action == meta.Delete || r.Action == meta.Create {
		if r.Basic.Type == meta.MainlineModel {
			return ModelTopologyOperation, nil
		}
	}

	if r.Action == meta.Find || r.Action == meta.Update {
		if r.Basic.Type == meta.ModelTopology {
			return ModelTopologyView, nil
		}
		if r.Basic.Type == meta.MainlineModelTopology {
			return ModelTopologyOperation, nil
		}

	}

	if r.Basic.Type == meta.Process {
		if r.Action == meta.BoundModuleToProcess || r.Action == meta.UnboundModuleToProcess {
			return Edit, nil
		}
	}

	if r.Basic.Type == meta.HostInstance {
		if r.Action == meta.MoveResPoolHostToBizIdleModule {
			return Edit, nil
		}

		if r.Action == meta.AddHostToResourcePool {
			return Create, nil
		}

		if r.Action == meta.MoveResPoolHostToBizIdleModule {
			return Edit, nil
		}
	}

	switch r.Action {
	case meta.Create, meta.CreateMany:
		return Create, nil

	case meta.Find, meta.FindMany:
		return Get, nil

	case meta.Delete, meta.DeleteMany:
		return Delete, nil

	case meta.Update, meta.UpdateMany:
		return Edit, nil

	case meta.MoveResPoolHostToBizIdleModule:
		if r.Basic.Type == meta.ModelInstance && r.Basic.Name == meta.Host {
			return Edit, nil
		}

	case meta.MoveHostToBizFaultModule,
		meta.MoveHostToBizIdleModule,
		meta.MoveHostToBizRecycleModule,
		meta.MoveHostToAnotherBizModule,
		meta.CleanHostInSetOrModule,
		meta.TransferHost,
		meta.MoveBizHostToModule:
		return Edit, nil

	case meta.MoveHostFromModuleToResPool:
		return Delete, nil

	case meta.MoveHostsToBusinessOrModule:
		return Edit, nil
	case meta.ModelTopologyView:
		return ModelTopologyView, nil
	case meta.ModelTopologyOperation:
		return ModelTopologyOperation, nil
	case meta.AdminEntrance:
		return AdminEntrance, nil
	}

	return Unknown, fmt.Errorf("unsupported action: %s", r.Action)
}

// Adaptor is a middleware wrapper which works for converting concepts
// between bk-cmdb and blueking auth center. Especially the policies
// in auth center.
func ConvertResourceType(resourceType meta.ResourceType, businessID int64) (*ResourceTypeID, error) {
	var iamResourceType ResourceTypeID
	switch resourceType {
	case meta.Business:
		iamResourceType = SysBusinessInstance

	case meta.Model,
		meta.ModelUnique,
		meta.ModelAttribute,
		meta.ModelAttributeGroup:
		if businessID > 0 {
			iamResourceType = BizModel
		} else {
			iamResourceType = SysModel
		}

	case meta.ModelModule, meta.ModelSet, meta.MainlineInstance, meta.MainlineInstanceTopology:
		iamResourceType = BizTopology

	case meta.MainlineModel, meta.ModelTopology:
		iamResourceType = SysSystemBase

	case meta.ModelClassification:
		if businessID > 0 {
			iamResourceType = BizModelGroup
		} else {
			iamResourceType = SysModelGroup
		}

	case meta.AssociationType:
		iamResourceType = SysAssociationType

	case meta.ModelAssociation:
		if businessID > 0 {
			iamResourceType = BizInstance
		} else {
			iamResourceType = SysInstance
		}
	case meta.ModelInstanceAssociation:
		if businessID > 0 {
			iamResourceType = BizInstance
		} else {
			iamResourceType = SysInstance
		}
	case meta.MainlineModelTopology:
		iamResourceType = SysSystemBase

	case meta.ModelInstance:
		if businessID <= 0 {
			iamResourceType = SysInstance
		} else {
			iamResourceType = BizInstance
		}

	case meta.Plat:
		iamResourceType = SysInstance
	case meta.HostInstance:
		if businessID <= 0 {
			iamResourceType = SysHostInstance
		} else {
			iamResourceType = BizHostInstance
		}

	case meta.HostFavorite:
		iamResourceType = BizHostInstance

	case meta.Process:
		iamResourceType = BizProcessInstance
	case meta.EventPushing:
		iamResourceType = SysEventPushing
	case meta.DynamicGrouping:
		iamResourceType = BizCustomQuery
	case meta.AuditLog:
		if businessID <= 0 {
			iamResourceType = SysAuditLog
		} else {
			iamResourceType = BizAuditLog
		}
	case meta.SystemBase:
		iamResourceType = SysSystemBase
	case meta.UserCustom:
		iamResourceType = UserCustom
	case meta.NetDataCollector:
		return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
	case meta.ProcessServiceTemplate:
		iamResourceType = BizProcessServiceTemplate
	case meta.ProcessServiceCategory:
		iamResourceType = BizProcessServiceCategory
	case meta.ProcessServiceInstance:
		iamResourceType = BizProcessServiceInstance
	case meta.BizTopology:
		iamResourceType = BizTopology
	case meta.SetTemplate:
		iamResourceType = BizSetTemplate

	case meta.OperationStatistic:
		iamResourceType = SysOperationStatistic
	default:
		return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
	}

	return &iamResourceType, nil
}

// ResourceTypeID is resource's type in auth center.
func adaptor(attribute *meta.ResourceAttribute) (*ResourceInfo, error) {
	var err error
	info := new(ResourceInfo)
	info.ResourceName = attribute.Basic.Name

	resourceTypeID, err := ConvertResourceType(attribute.Type, attribute.BusinessID)
	if err != nil {
		return nil, err
	}
	info.ResourceType = *resourceTypeID

	info.ResourceID, err = GenerateResourceID(info.ResourceType, attribute)
	if err != nil {
		return nil, err
	}

	return info, nil
}


