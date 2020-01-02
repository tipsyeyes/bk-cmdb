/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package service

import (
	"encoding/json"
	"net/http"

	"configdatabase/src/common"
	"configdatabase/src/common/blog"
	"configdatabase/src/common/condition"
	meta "configdatabase/src/common/metadata"
	"configdatabase/src/common/util"
	"github.com/emicklei/go-restful"
)

func (s *Service) FindModuleHost(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	defErr := srvData.ccErr

	body := new(meta.HostModuleFind)
	if err := json.NewDecoder(req.Request.Body).Decode(body); err != nil {
		blog.Errorf("find host failed with decode body err: %#v, rid:%s", err, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	host, err := srvData.lgc.FindHostByModuleIDs(srvData.ctx, body, false)
	if err != nil {
		blog.Errorf("find host failed, err: %#v, input:%#v, rid:%s", err, body, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrHostGetFail)})
		return
	}

	_ = resp.WriteEntity(meta.SearchHostResult{
		BaseResp: meta.SuccessBaseResp,
		Data:     *host,
	})
}

// ListHosts list host under business specified by path parameter
func (s *Service) ListBizHosts(req *restful.Request, resp *restful.Response) {
	header := req.Request.Header
	ctx := util.NewContextFromHTTPHeader(header)
	rid := util.ExtractRequestIDFromContext(ctx)
	srvData := s.newSrvComm(header)
	defErr := srvData.ccErr

	parameter := &meta.ListHostsParameter{}
	if err := json.NewDecoder(req.Request.Body).Decode(parameter); err != nil {
		blog.Errorf("ListHostByTopoNode failed, decode body failed, err: %#v, rid:%s", err, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}
	bizID, err := util.GetInt64ByInterface(req.PathParameter("appid"))
	if err != nil {
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrCommParamsInvalid, "bk_app_id")})
		return
	}
	if bizID == 0 {
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrCommParamsInvalid, "bk_app_id")})
		return
	}

	if parameter.Page.Limit == 0 {
		parameter.Page.Limit = common.BKMaxPageSize
	}
	if parameter.Page.Limit > common.BKMaxPageLimit {
		blog.Errorf("ListBizHosts failed, page limit %d exceed max pageSize %d, rid:%s", parameter.Page.Limit, common.BKMaxPageLimit, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommPageLimitIsExceeded)})
		return
	}
	option := meta.ListHosts{
		BizID:              bizID,
		SetIDs:             parameter.SetIDs,
		ModuleIDs:          parameter.ModuleIDs,
		HostPropertyFilter: parameter.HostPropertyFilter,
		Page:               parameter.Page,
	}
	host, err := s.CoreAPI.CoreService().Host().ListHosts(ctx, header, option)
	if err != nil {
		blog.Errorf("find host failed, err: %s, input:%#v, rid:%s", err.Error(), parameter, rid)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrHostGetFail)})
		return
	}

	result := meta.NewSuccessResponse(host)
	_ = resp.WriteEntity(result)
}

// ListHostsWithNoBiz list host for no biz case merely
func (s *Service) ListHostsWithNoBiz(req *restful.Request, resp *restful.Response) {
	header := req.Request.Header
	ctx := util.NewContextFromHTTPHeader(header)
	rid := util.ExtractRequestIDFromContext(ctx)
	srvData := s.newSrvComm(header)
	defErr := srvData.ccErr

	parameter := &meta.ListHostsWithNoBizParameter{}
	if err := json.NewDecoder(req.Request.Body).Decode(parameter); err != nil {
		blog.Errorf("ListHostsWithNoBiz failed, decode body failed, err: %#v, rid:%s", err, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	if parameter.Page.Limit == 0 {
		parameter.Page.Limit = common.BKMaxPageSize
	}
	if parameter.Page.Limit > common.BKMaxPageLimit {
		blog.Errorf("ListHostsWithNoBiz failed, page limit %d exceed max pageSize %d, rid:%s", parameter.Page.Limit, common.BKMaxPageLimit, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommPageLimitIsExceeded)})
		return
	}
	option := meta.ListHosts{
		HostPropertyFilter: parameter.HostPropertyFilter,
		Page:               parameter.Page,
	}
	host, err := s.CoreAPI.CoreService().Host().ListHosts(ctx, header, option)
	if err != nil {
		blog.Errorf("find host failed, err: %s, input:%#v, rid:%s", err.Error(), parameter, rid)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrHostGetFail)})
		return
	}

	result := meta.NewSuccessResponse(host)
	_ = resp.WriteEntity(result)
}

// ListBizHostsTopo list hosts under business specified by path parameter with their topology information
func (s *Service) ListBizHostsTopo(req *restful.Request, resp *restful.Response) {
	header := req.Request.Header
	ctx := util.NewContextFromHTTPHeader(header)
	rid := util.ExtractRequestIDFromContext(ctx)
	srvData := s.newSrvComm(header)
	defErr := srvData.ccErr

	parameter := &meta.ListHostsWithNoBizParameter{}
	if err := json.NewDecoder(req.Request.Body).Decode(parameter); err != nil {
		blog.Errorf("ListHostByTopoNode failed, decode body failed, err: %#v, rid:%s", err, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}
	bizID, err := util.GetInt64ByInterface(req.PathParameter("bk_biz_id"))
	if err != nil {
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrCommParamsInvalid, "bk_app_id")})
		return
	}
	if bizID == 0 {
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrCommParamsInvalid, "bk_app_id")})
		return
	}

	if parameter.Page.Limit == 0 {
		parameter.Page.Limit = common.BKMaxPageSize
	}
	if parameter.Page.Limit > common.BKMaxPageLimit {
		blog.Errorf("ListBizHosts failed, page limit %d exceed max pageSize %d, rid:%s", parameter.Page.Limit, common.BKMaxPageLimit, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommPageLimitIsExceeded)})
		return
	}

	// search all hosts
	option := meta.ListHosts{
		BizID:              bizID,
		HostPropertyFilter: parameter.HostPropertyFilter,
		Page:               parameter.Page,
	}
	hosts, err := s.CoreAPI.CoreService().Host().ListHosts(ctx, header, option)
	if err != nil {
		blog.Errorf("find host failed, err: %s, input:%#v, rid:%s", err.Error(), parameter, rid)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrHostGetFail)})
		return
	}

	if len(hosts.Info) == 0 {
		result := meta.NewSuccessResponse(hosts)
		_ = resp.WriteEntity(result)
		return
	}

	// search all hosts' host module relations
	hostIDs := make([]int64, 0)
	for _, host := range hosts.Info {
		hostID, err := util.GetInt64ByInterface(host[common.BKHostIDField])
		if err != nil {
			blog.ErrorJSON("host: %s bk_host_id field invalid, rid:%s", host, rid)
			_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
			return
		}
		hostIDs = append(hostIDs, hostID)
	}
	relationCond := meta.HostModuleRelationRequest{
		ApplicationID: bizID,
		HostIDArr:     hostIDs,
	}
	relations, err := srvData.lgc.GetConfigByCond(srvData.ctx, relationCond)
	if nil != err {
		blog.ErrorJSON("read host module relation error: %s, input: %s, rid: %s", err, hosts, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
		return
	}

	// search all module and set info
	setIDs := make([]int64, 0)
	moduleIDs := make([]int64, 0)
	relation := make(map[int64]map[int64][]int64)
	for _, r := range relations {
		setIDs = append(setIDs, r.SetID)
		moduleIDs = append(moduleIDs, r.ModuleID)
		if setModule, ok := relation[r.HostID]; ok {
			setModule[r.SetID] = append(setModule[r.SetID], r.ModuleID)
			relation[r.HostID] = setModule
		} else {
			setModule := make(map[int64][]int64)
			setModule[r.SetID] = append(setModule[r.SetID], r.ModuleID)
			relation[r.HostID] = setModule
		}
	}
	setIDs = util.IntArrayUnique(setIDs)
	moduleIDs = util.IntArrayUnique(moduleIDs)

	cond := condition.CreateCondition()
	cond.Field(common.BKSetIDField).In(setIDs)
	query := &meta.QueryCondition{
		Fields:    []string{common.BKSetIDField, common.BKSetNameField},
		Condition: cond.ToMapStr(),
	}
	sets, err := s.CoreAPI.CoreService().Instance().ReadInstance(ctx, header, common.BKInnerObjIDSet, query)
	if err != nil {
		blog.ErrorJSON("get set by condition: %s failed, err: %+v, rid: %s", cond.ToMapStr(), err, rid)
		_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
		return
	}
	setMap := make(map[int64]string)
	for _, set := range sets.Data.Info {
		setID, err := set.Int64(common.BKSetIDField)
		if err != nil {
			blog.ErrorJSON("set %s id invalid, error: %s, rid: %s", set, err, rid)
			_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
			return
		}
		setName, err := set.String(common.BKSetNameField)
		if err != nil {
			blog.ErrorJSON("set %s name invalid, error: %s, rid: %s", set, err, rid)
			_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
			return
		}
		setMap[setID] = setName
	}

	cond = condition.CreateCondition()
	cond.Field(common.BKModuleIDField).In(moduleIDs)
	query = &meta.QueryCondition{
		Fields:    []string{common.BKModuleIDField, common.BKModuleNameField},
		Condition: cond.ToMapStr(),
	}
	modules, err := s.CoreAPI.CoreService().Instance().ReadInstance(ctx, header, common.BKInnerObjIDModule, query)
	if err != nil {
		blog.ErrorJSON("get module by condition: %s failed, err: %+v, rid: %s", cond.ToMapStr(), err, rid)
		_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
		return
	}
	moduleMap := make(map[int64]string)
	for _, module := range modules.Data.Info {
		moduleID, err := module.Int64(common.BKModuleIDField)
		if err != nil {
			blog.ErrorJSON("module %s id invalid, error: %s, rid: %s", module, err, rid)
			_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
			return
		}
		moduleName, err := module.String(common.BKModuleNameField)
		if err != nil {
			blog.ErrorJSON("module %s name invalid, error: %s, rid: %s", module, err, rid)
			_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
			return
		}
		moduleMap[moduleID] = moduleName
	}

	// format the output
	hostTopos := meta.HostTopoResult{
		Count: hosts.Count,
	}
	for _, host := range hosts.Info {
		hostTopo := meta.HostTopo{
			Host: host,
		}
		topos := make([]meta.Topo, 0)
		hostID, _ := util.GetInt64ByInterface(host[common.BKHostIDField])
		if setModule, ok := relation[hostID]; ok {
			for setID, moduleIDs := range setModule {
				topo := meta.Topo{
					SetID:   setID,
					SetName: setMap[setID],
				}
				modules := make([]meta.Module, 0)
				for _, moduleID := range moduleIDs {
					module := meta.Module{
						ModuleID:   moduleID,
						ModuleName: moduleMap[moduleID],
					}
					modules = append(modules, module)
				}
				topo.Module = modules
				topos = append(topos, topo)
			}
		}
		hostTopo.Topo = topos
		hostTopos.Info = append(hostTopos.Info, hostTopo)
	}

	result := meta.NewSuccessResponse(hostTopos)
	_ = resp.WriteEntity(result)
}
