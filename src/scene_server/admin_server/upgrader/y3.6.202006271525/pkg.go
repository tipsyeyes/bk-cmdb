package y3_6_202006271525

import (
	"configdatabase/src/common/blog"
	"configdatabase/src/scene_server/admin_server/upgrader"
	"configdatabase/src/storage/dal"
	"context"
)

func init() {
	upgrader.RegistUpgrader("y3.6.202006271525", upgrade)
}

func upgrade(ctx context.Context, db dal.RDB, conf *upgrader.Config) (err error) {
	err = addCustomModel(ctx, db, conf)
	if err != nil {
		blog.Errorf("[upgrade y3_6_202006271525] addCustomModel error  %s", err.Error())
		return
	}
	return
}