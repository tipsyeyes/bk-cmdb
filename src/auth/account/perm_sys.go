package account

import "strings"

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
				ActionID:          Create,
				ActionName:        "新建",
				IsRelatedResource: false,
			},
			{
				ActionID:          Edit,
				ActionName:        "编辑",
				IsRelatedResource: true,
			},
			{
				ActionID:          Archive,
				ActionName:        "归档",
				IsRelatedResource: true,
			},
			{
				ActionID:          Get,
				ActionName:        "查询",
				IsRelatedResource: true,
			},
		},
	},
}



