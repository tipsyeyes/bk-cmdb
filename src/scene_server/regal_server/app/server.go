// the main app

package app

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"configdatabase/src/auth"
	"configdatabase/src/auth/authcenter"
	"configdatabase/src/auth/extensions"
	"configdatabase/src/common"
	"configdatabase/src/common/backbone"
	cc "configdatabase/src/common/backbone/configcenter"
	"configdatabase/src/common/blog"
	"configdatabase/src/common/types"
	"configdatabase/src/common/version"
	"configdatabase/src/scene_server/regal_server/app/options"
	svc "configdatabase/src/scene_server/regal_server/service"
	"configdatabase/src/storage/dal/mongo"
	"configdatabase/src/storage/dal/redis"
)

func Run(ctx context.Context, cancel context.CancelFunc, op *options.ServerOption) error {

	svrInfo, err := newServerInfo(op)
	if err != nil {
		blog.Errorf("fail to new server information. err: %s", err.Error())
		return fmt.Errorf("make server information failed, err:%v", err)
	}

	regal := new(RegalServer)
	service := new(svc.Service)

	input := &backbone.BackboneParameter{
		ConfigUpdate: regal.onRegalConfigUpdate,
		ConfigPath:   op.ServConf.ExConfig,
		Regdiscv:     op.ServConf.RegDiscover,
		SrvInfo:      svrInfo,
	}
	engine, err := backbone.NewBackbone(ctx, input)
	if err != nil {
		return fmt.Errorf("new backbone failed, err: %v", err)
	}
	configReady := false
	for sleepCnt := 0; sleepCnt < common.APPConfigWaitTime; sleepCnt++ {
		if nil == regal.Config {
			blog.Infof("waiting for config ready ...")
			time.Sleep(time.Second)
		} else {
			configReady = true
			break
		}
	}
	if false == configReady {
		blog.Errorf("waiting config timeout.")
		return fmt.Errorf("configuration item not found")
	}

	cacheDB, err := redis.NewFromConfig(regal.Config.Redis)
	if err != nil {
		blog.Errorf("new redis client failed, err: %s", err.Error())
		return fmt.Errorf("new redis client failed, err: %s", err.Error())
	}
	db, err := regal.Config.Mongo.GetMongoClient(engine)
	if err != nil {
		blog.Errorf("new mongo client failed, err: %s", err.Error())
		return fmt.Errorf("new mongo client failed, err: %s", err.Error())
	}

	// 创建权限中心资源授权处理接口
	// 目前需要操作的模块包括：
	// datacollection、synchronizeserver、hostserver、procserver、operationserver、adminserver、toposerver、eventserver
	authorize, err := auth.NewAuthorize(nil, regal.Config.Auth, engine.Metric().Registry())
	if err != nil {
		return fmt.Errorf("new authorize failed, err: %v", err)
	}
	authManager := extensions.NewAuthManager(engine.CoreAPI, authorize)

	service.AuthManager = authManager
	service.Engine = engine
	service.Config = regal.Config
	service.CacheDB = cacheDB
	service.DB = db
	regal.Core = engine
	regal.Service = service

	err = backbone.StartServer(ctx, cancel, engine, service.WebService(), true)
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		blog.Infof("process will exit!")
	}

	return nil
}

// main server
type RegalServer struct {
	Config  *options.Config
	Core    *backbone.Engine
	Service *svc.Service
}

var configLock sync.Mutex

func (r *RegalServer) onRegalConfigUpdate(previous, current cc.ProcessConfig) {
	var err error
	configLock.Lock()
	defer configLock.Unlock()

	if len(current.ConfigMap) <= 0 {
		blog.Errorf("configuration was empty.")
		return
	}
	if r.Config == nil {
		r.Config = new(options.Config)
	}
	// 解析redis信息
	r.Config.Redis.Address = current.ConfigMap["redis.host"]
	r.Config.Redis.Database = current.ConfigMap["redis.database"]
	r.Config.Redis.Password = current.ConfigMap["redis.pwd"]
	r.Config.Redis.Port = current.ConfigMap["redis.port"]
	r.Config.Redis.MasterName = current.ConfigMap["redis.user"]
	// r.Config.Mongo = mongo.ParseConfigFromKV("mongodb", current.ConfigMap)

	// 解析mongo信息
	r.Config.Mongo = mongo.ParseConfigFromKV("mongodb", current.ConfigMap)

	// 解析权限中心配置
	r.Config.Auth, err = authcenter.ParseConfigFromKV("auth", current.ConfigMap)
	if err != nil {
		blog.Warnf("parse auth center config failed: %v", err)
	}
}

func newServerInfo(op *options.ServerOption) (*types.ServerInfo, error) {
	ip, err := op.ServConf.GetAddress()
	if err != nil {
		return nil, err
	}

	port, err := op.ServConf.GetPort()
	if err != nil {
		return nil, err
	}

	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	svrInfo := &types.ServerInfo{
		IP:       ip,
		Port:     port,
		HostName: hostname,
		Scheme:   "http",
		Version:  version.GetVersion(),
		Pid:      os.Getpid(),
	}

	return svrInfo, nil
}
