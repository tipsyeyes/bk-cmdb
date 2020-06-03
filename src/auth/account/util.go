package account

var ResourceTypeMap = make(map[ResourceTypeID]ResourceType)

func init() {
	for _, bizResourceType := range expectBizResourceType {
		ResourceTypeMap[bizResourceType.ResourceTypeID] = bizResourceType
	}
	for _, sysResourceType := range expectSystemResourceType {
		ResourceTypeMap[sysResourceType.ResourceTypeID] = sysResourceType
	}
}

// IsRelatedToResourceID check whether authorization on this resourceType need resourceID
func IsRelatedToResourceID(resourceTypeID ResourceTypeID) bool {
	resourceType, exist := ResourceTypeMap[resourceTypeID]
	if exist == false {
		return false
	}
	for _, action := range resourceType.Actions {
		if action.IsRelatedResource == true {
			return true
		}
	}
	return false
}
