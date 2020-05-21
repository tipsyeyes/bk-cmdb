// the main service code

package service

import (
	"context"
	"io"
	"net/http"

	"configdatabase/src/auth/extensions"
	"configdatabase/src/common"
	"configdatabase/src/common/backbone"
	"configdatabase/src/common/errors"
	"configdatabase/src/common/language"
	"configdatabase/src/common/metadata"
	"configdatabase/src/common/metric"
	"configdatabase/src/common/rdapi"
	"configdatabase/src/common/types"
	"configdatabase/src/common/util"
	"configdatabase/src/scene_server/regal_server/app/options"
	"configdatabase/src/storage/dal"

	"github.com/emicklei/go-restful"
	"gopkg.in/redis.v5"
)

type Service struct {
	*options.Config
	*backbone.Engine

	CacheDB *redis.Client
	DB      dal.RDB

	AuthManager *extensions.AuthManager
}

type srvComm struct {
	header        http.Header
	rid           string
	ccErr         errors.DefaultCCErrorIf
	ccLang        language.DefaultCCLanguageIf
	ctx           context.Context
	ctxCancelFunc context.CancelFunc
	user          string
	ownerID       string
	// lgc           *logics.Logics
}

func (s *Service) newSrvComm(header http.Header) *srvComm {
	rid := util.GetHTTPCCRequestID(header)
	lang := util.GetLanguage(header)
	user := util.GetUser(header)
	ctx, cancel := s.Engine.CCCtx.WithCancel()
	ctx = context.WithValue(ctx, common.ContextRequestIDField, rid)
	ctx = context.WithValue(ctx, common.ContextRequestUserField, user)

	return &srvComm{
		header:        header,
		rid:           util.GetHTTPCCRequestID(header),
		ccErr:         s.CCErr.CreateDefaultCCErrorIf(lang),
		ccLang:        s.Language.CreateDefaultCCLanguageIf(lang),
		ctx:           ctx,
		ctxCancelFunc: cancel,
		user:          util.GetUser(header),
		ownerID:       util.GetOwnerID(header),
		//lgc:           logics.NewLogics(s.Engine, header, s.AuthManager),
	}
}

func (s *Service) WebService() *restful.Container {
	container := restful.NewContainer()

	getErrFunc := func() errors.CCErrorIf {
		return s.CCErr
	}
	api := new(restful.WebService)
	api.Path("/regal/v3/").Filter(rdapi.AllGlobalFilter(getErrFunc)).Produces(restful.MIME_JSON)
	api.Route(api.GET("/hello").To(s.hello))
	container.Add(api)

	healthzAPI := new(restful.WebService).Produces(restful.MIME_JSON)
	healthzAPI.Route(healthzAPI.GET("/healthz").To(s.healthz))
	container.Add(healthzAPI)

	return container;
}

func (s *Service) healthz(req *restful.Request, resp *restful.Response) {
	meta := metric.HealthMeta{IsHealthy: true}

	// zk health status
	zkItem := metric.HealthItem{IsHealthy: true, Name: types.CCFunctionalityServicediscover}
	if err := s.Engine.Ping(); err != nil {
		zkItem.IsHealthy = false
		zkItem.Message = err.Error()
	}
	meta.Items = append(meta.Items, zkItem)

	// mongodb
	healthItem := metric.NewHealthItem(types.CCFunctionalityMongo, s.DB.Ping())
	meta.Items = append(meta.Items, healthItem)

	// redis
	redisItem := metric.NewHealthItem(types.CCFunctionalityRedis, s.CacheDB.Ping().Err())
	meta.Items = append(meta.Items, redisItem)

	for _, item := range meta.Items {
		if item.IsHealthy == false {
			meta.IsHealthy = false
			meta.Message = "event server is unhealthy"
			break
		}
	}

	info := metric.HealthInfo{
		Module:     types.CC_MODULE_REGAL,
		HealthMeta: meta,
		AtTime:     metadata.Now(),
	}

	answer := metric.HealthResponse{
		Code:    common.CCSuccess,
		Data:    info,
		OK:      meta.IsHealthy,
		Result:  meta.IsHealthy,
		Message: meta.Message,
	}
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteEntity(answer)
}

func (s *Service) hello(req *restful.Request, resp *restful.Response) {
	io.WriteString(resp, "world")
}