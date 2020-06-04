package account

import (
	"strings"
)

var expectSystem = System{
	SystemID:           SystemIDCMDB,
	SystemName:         SystemNameCMDB,
	Desc:               "巨人运维配置平台",
	RelatedScopeTypes:	strings.Join([]string{ScopeTypeIDBiz, ScopeTypeIDSystem}, ";"),
	Managers:           "admin",
	Creator:            "admin",
	Updater:            "admin",
}

// 全局级资源信息
var expectSystemResourceType = []ResourceType{
	{
		ResourceTypeID:       SysBusinessInstance,
		ResourceTypeName:     ResourceTypeIDMap[SysBusinessInstance],
		ParentResourceTypeID: "",
		Share:                true,
		Actions: []Action{
			{
				ActionID:          Get,
				ActionName:        ActionIDNameMap[Get],
				IsFunctional:	   false,
				IsRelatedResource: true,
			},
			{
				ActionID:          Create,
				ActionName:        ActionIDNameMap[Create],
				IsFunctional:	   true,
				IsRelatedResource: false,
			},
			{
				ActionID:          Edit,
				ActionName:        ActionIDNameMap[Edit],
				IsFunctional:	   false,
				IsRelatedResource: true,
			},
			{
				ActionID:          Archive,
				ActionName:        ActionIDNameMap[Archive],
				IsFunctional:	   false,
				IsRelatedResource: true,
			},

		},
	},
}



