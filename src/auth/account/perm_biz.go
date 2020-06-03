package account

var expectBizResourceType = []ResourceType{
	{
		ResourceTypeID:       BizHostInstance,
		ResourceTypeName:     "主机",
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
				ActionID:          Delete,
				ActionName:        "删除/归还",
				IsRelatedResource: true,
			},
		},
	},
}
