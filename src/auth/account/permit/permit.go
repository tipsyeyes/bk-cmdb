package permit

import (
	"configdatabase/src/auth/meta"
)

func IsReadAction(action meta.Action) bool {
	if action == meta.FindMany || action == meta.Find {
		return true
	}
	return false
}

// this function is used to check where this request attribute is permitted as default,
// so that it is not need to check permission status in auth center.
func ShouldSkipAuthorize(rsc *meta.ResourceAttribute) bool {

	switch {

	case rsc.Type == meta.AuditLog:
		return false
	case rsc.Type == meta.ResourceSync:
		return true
	// case rsc.Type == meta.DynamicGrouping && rsc.Action == meta.FindMany || rsc.Action == meta.Find:
	// 	return true

	case rsc.Type == meta.ModelClassification && IsReadAction(rsc.Action):
		return true

	case rsc.Type == meta.AssociationType && IsReadAction(rsc.Action):
		return true

	case rsc.Type == meta.Model && IsReadAction(rsc.Action):
		return true
	case rsc.Type == meta.ModelAttribute && IsReadAction(rsc.Action):
		return true
	case rsc.Type == meta.ModelUnique && IsReadAction(rsc.Action):
		return true
	case rsc.Type == meta.ModelAttributeGroup && IsReadAction(rsc.Action):
		return true

	case rsc.Type == meta.UserCustom:
		return true

	case rsc.Type == meta.ModelAssociation && IsReadAction(rsc.Action):
		return true

	// all the model instance association related operation is all authorized for now.
	case rsc.Type == meta.ModelInstanceAssociation && IsReadAction(rsc.Action):
		return true

	// case rsc.Type == meta.ModelInstance && (rsc.Action == meta.Find || rsc.Action == meta.FindMany):
	// 	return true

	// all the network data collector related operation is all authorized for now.
	case rsc.Type == meta.NetDataCollector:
		return true
	// host search operation, skip.
	case rsc.Type == meta.HostInstance && IsReadAction(rsc.Action):
		return true

	// topology instance resource types.
	case rsc.Type == meta.ModelInstanceTopology || rsc.Type == meta.MainlineInstanceTopology:
		return true
	case rsc.Type == meta.ProcessServiceInstance && IsReadAction(rsc.Action):
		return true
	case rsc.Type == meta.ProcessServiceTemplate && IsReadAction(rsc.Action):
		return true
	case rsc.Type == meta.SetTemplate && IsReadAction(rsc.Action):
		return true
	case rsc.Type == meta.ProcessServiceCategory && IsReadAction(rsc.Action):
		return true
	case rsc.Type == meta.ProcessTemplate && IsReadAction(rsc.Action):
		return true

	case rsc.Type == meta.MainlineInstance && IsReadAction(rsc.Action):
		return true
	case rsc.Type == meta.Process && IsReadAction(rsc.Action):
		return true
	case rsc.Type == meta.ModelSet && IsReadAction(rsc.Action):
		return true
	case rsc.Type == meta.ModelModule && IsReadAction(rsc.Action):
		return true
	case rsc.Type == meta.MainlineModelTopology && IsReadAction(rsc.Action):
		return true
	case rsc.Type == meta.OperationStatistic && IsReadAction(rsc.Action):
		return true
	case rsc.Type == meta.HostFavorite:
		return true
	case rsc.Type == meta.InstallBK:
		return true
	default:
		return false
	}

	return false
}
