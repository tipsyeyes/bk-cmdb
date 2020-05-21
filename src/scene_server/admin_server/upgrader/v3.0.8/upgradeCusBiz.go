package v3v0v8

import (
	"context"
	"time"

	"configdatabase/src/common"
	"configdatabase/src/common/blog"
	"configdatabase/src/scene_server/admin_server/upgrader"
	"configdatabase/src/storage/dal"

)

// add by tes
// addCusApp add custom app
func addCusApp(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	// add by tes
	// 增加2个测试的业务
	cusAppNames := []string{"王者一号", "王者二号"}
	for _, appName := range cusAppNames {
		cusBiz := map[string]interface{}{}
		cusBiz[common.BKAppNameField] = appName
		cusBiz[common.BKMaintainersField] = admin
		cusBiz[common.BKTimeZoneField] = "Asia/Shanghai"
		cusBiz[common.BKLanguageField] = "1" //中文
		cusBiz[common.BKLifeCycleField] = common.DefaultAppLifeCycleNormal
		cusBiz[common.BKOwnerIDField] = conf.OwnerID
		cusBiz[common.BKSupplierIDField] = common.BKDefaultSupplierID
		cusBiz[common.BKDefaultField] = common.DefaultFlagDefaultValue
		cusBiz[common.CreateTimeField] = time.Now()
		cusBiz[common.LastTimeField] = time.Now()
		filled := fillEmptyFields(cusBiz, AppRow())
		bizID, _, err := upgrader.Upsert(ctx, db, "cc_ApplicationBase", cusBiz, common.BKAppIDField, []string{common.BKOwnerIDField, common.BKAppNameField, common.BKDefaultField}, append(filled, common.BKAppIDField))
		if err != nil {
			blog.Error("add cusBiz error ", err.Error())
			return err
		}

		// add cus app default set
		inputSetInfo := make(map[string]interface{})
		inputSetInfo[common.BKAppIDField] = bizID
		inputSetInfo[common.BKInstParentStr] = bizID
		inputSetInfo[common.BKSetNameField] = common.DefaultResSetName
		inputSetInfo[common.BKDefaultField] = common.DefaultResSetFlag
		inputSetInfo[common.BKOwnerIDField] = conf.OwnerID
		filled = fillEmptyFields(inputSetInfo, SetRow())
		setID, _, err := upgrader.Upsert(ctx, db, "cc_SetBase", inputSetInfo, common.BKSetIDField, []string{common.BKOwnerIDField, common.BKAppIDField, common.BKSetNameField}, append(filled, common.BKSetIDField))
		if err != nil {
			blog.Error("add defaultSet error ", err.Error())
			return err
		}

		// add cus app default module
		inputResModuleInfo := make(map[string]interface{})
		inputResModuleInfo[common.BKSetIDField] = setID
		inputResModuleInfo[common.BKInstParentStr] = setID
		inputResModuleInfo[common.BKAppIDField] = bizID
		inputResModuleInfo[common.BKModuleNameField] = common.DefaultResModuleName
		inputResModuleInfo[common.BKDefaultField] = common.DefaultResModuleFlag
		inputResModuleInfo[common.BKOwnerIDField] = conf.OwnerID
		filled = fillEmptyFields(inputResModuleInfo, ModuleRow())
		_, _, err = upgrader.Upsert(ctx, db, "cc_ModuleBase", inputResModuleInfo, common.BKModuleIDField, []string{common.BKOwnerIDField, common.BKModuleNameField, common.BKAppIDField, common.BKSetIDField}, append(filled, common.BKModuleIDField))
		if err != nil {
			blog.Error("add defaultResModule error ", err.Error())
			return err
		}

		inputFaultModuleInfo := make(map[string]interface{})
		inputFaultModuleInfo[common.BKSetIDField] = setID
		inputFaultModuleInfo[common.BKInstParentStr] = setID
		inputFaultModuleInfo[common.BKAppIDField] = bizID
		inputFaultModuleInfo[common.BKModuleNameField] = common.DefaultFaultModuleName
		inputFaultModuleInfo[common.BKDefaultField] = common.DefaultFaultModuleFlag
		inputFaultModuleInfo[common.BKOwnerIDField] = conf.OwnerID
		filled = fillEmptyFields(inputFaultModuleInfo, ModuleRow())
		_, _, err = upgrader.Upsert(ctx, db, "cc_ModuleBase", inputFaultModuleInfo, common.BKModuleIDField, []string{common.BKOwnerIDField, common.BKModuleNameField, common.BKAppIDField, common.BKSetIDField}, append(filled, common.BKModuleIDField))
		if err != nil {
			blog.Error("add defaultFaultModule error ", err.Error())
			return err
		}
	}
	return nil
}
