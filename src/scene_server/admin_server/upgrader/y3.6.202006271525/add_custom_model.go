package y3_6_202006271525

import (
	"configdatabase/src/common"
	"configdatabase/src/common/blog"
	"configdatabase/src/common/condition"
	"configdatabase/src/common/metadata"
	mCommon "configdatabase/src/scene_server/admin_server/common"
	"configdatabase/src/scene_server/admin_server/upgrader"
	"configdatabase/src/storage/dal"
	"context"
	"errors"
)

func addCustomModel(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	// 添加自定义配置模型对象，关联到逻辑资源组
	ownerID := conf.OwnerID
	objID := common.BKInnerObjIDCONFIG
	t := metadata.Now()
	tableName := common.BKTableNameObjDes
	obj := &metadata.Object{ObjCls: "bk_logic", ObjectID: objID, ObjectName: "自定义配置模型", IsPre: true, ObjIcon: "icon-cc-port", Position: ``}
	obj.CreateTime = &t
	obj.LastTime = &t
	obj.IsPaused = false
	obj.Creator = common.CCSystemOperatorUserName
	obj.OwnerID = ownerID
	obj.Description = ""
	obj.Modifier = ""
	if _, _, err := upgrader.Upsert(ctx, db, tableName, obj, "id", []string{common.BKObjIDField, common.BKClassificationIDField, common.BKOwnerIDField}, []string{"id"}); err != nil {
		blog.Errorf("add data for %s table error: %s", tableName, err)
		return err
	}

	// 添加模型属性分组信息
	tableName = common.BKTableNamePropertyGroup
	groupRows := []*metadata.Group{
		&metadata.Group{ObjectID: objID, GroupID: mCommon.BaseInfo, GroupName: mCommon.BaseInfoName, GroupIndex: 1, OwnerID: ownerID, IsDefault: true},
	}
	for _, row := range groupRows {
		if _, _, err := upgrader.Upsert(ctx, db, tableName, row, "id", []string{common.BKObjIDField, "bk_group_id"}, []string{"id"}); err != nil {
			blog.Errorf("add data for %s table error: %s", tableName, err)
			return err
		}
	}

	// 添加模型属性字段信息
	tableName = common.BKTableNameObjAttDes
	attrRows := []*metadata.Attribute{
		&metadata.Attribute{ObjectID: objID, PropertyID: common.BKInstNameField, PropertyName: "配置名称", IsRequired: true, IsOnly: false, IsPre: true, IsEditable: true, PropertyGroup: mCommon.BaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_inst_value", PropertyName: "配置值", IsRequired: true, IsOnly: false, IsPre: true, IsEditable: true, PropertyGroup: mCommon.BaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
	}
	for _, r := range attrRows {
		r.OwnerID = ownerID
		r.IsPre = true
		if false != r.IsEditable {
			r.IsEditable = true
		}
		r.IsReadOnly = false
		r.CreateTime = &t
		r.LastTime = &t
		r.Creator = common.CCSystemOperatorUserName
		r.Description = ""
	}
	for _, row := range attrRows {
		_, _, err := upgrader.Upsert(ctx, db, tableName, row, "id", []string{common.BKObjIDField, common.BKPropertyIDField, common.BKOwnerIDField}, []string{})
		if nil != err {
			blog.Errorf("add data for %s table error: %s", tableName, err)
			return err
		}
	}

	// 添加模型校验
	attrs := make([]metadata.Attribute, 0)
	attrCond := condition.CreateCondition()
	attrCond.Field(common.BKObjIDField).In([]string{
		objID,
	})
	err := db.Table(tableName).Find(attrCond.ToMapStr()).All(ctx, &attrs)
	if nil != err {
		blog.Errorf("find data for %s table error: %s", tableName, err)
		return err
	} else if len(attrs) != 2 {
		blog.Errorf("find data for %s table error", tableName)
		return errors.New("find attrs failed")
	}
	propertyIDToProperty := make(map[string]metadata.Attribute)
	for _, attr := range attrs {
		propertyIDToProperty[attr.PropertyID] = attr
	}

	tableName = common.BKTableNameObjUnique
	objUnique := metadata.ObjectUnique{
		ObjID:     objID,
		MustCheck: true,
		Keys: []metadata.UniqueKey{
			{
				Kind: metadata.UniqueKeyKindProperty,
				ID:   uint64(propertyIDToProperty[common.BKInstNameField].ID),
			},
		},
		Ispre:    false,
		OwnerID:  ownerID,
		LastTime: t,
	}
	uid, err := db.NextSequence(ctx, tableName)
	if err != nil {
		return err
	}
	objUnique.ID = uid
	if err := db.Table(tableName).Insert(ctx, objUnique); err != nil {
		return err
	}

	return nil
}