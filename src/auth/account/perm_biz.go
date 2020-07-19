package account

var expectBizResourceType = []ResourceType{
	//{
	//	ResourceTypeID:       BizHostInstance,
	//	ResourceTypeName:     "主机",
	//	ParentResourceTypeID: "",
	//	Share:                true,
	//	Actions: []Action{
	//		{
	//			ActionID:          Create,
	//			ActionName:        "新建",
	//			IsFunctional:	   true,
	//			IsRelatedResource: false,
	//		},
	//		{
	//			ActionID:          Edit,
	//			ActionName:        "编辑",
	//			IsFunctional:	   false,
	//			IsRelatedResource: true,
	//		},
	//		{
	//			ActionID:          Delete,
	//			ActionName:        "删除/归还",
	//			IsFunctional:	   false,
	//			IsRelatedResource: true,
	//		},
	//	},
	//},
	{
		ResourceTypeID:       BizRBizInstance,
		ResourceTypeName:     "业务",
		ParentResourceTypeID: "",
		Share:                true,
		Actions: []Action{
			{
				ActionID:          Create,
				ActionName:        "新建",
				IsFunctional:	   true,
				IsRelatedResource: false,
			},
			{
				ActionID:          Edit,
				ActionName:        "编辑",
				IsFunctional:	   false,
				IsRelatedResource: true,
			},
			{
				ActionID:          Delete,
				ActionName:        "删除",
				IsFunctional:	   false,
				IsRelatedResource: true,
			},
			{
				ActionID:          Get,
				ActionName:        "查询",
				IsFunctional:	   false,
				IsRelatedResource: true,
			},
		},
	},
}
