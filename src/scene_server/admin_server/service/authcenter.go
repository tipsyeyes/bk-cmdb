/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2019 THL A29 Limited, a Tencent company. All rights reserved.
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
	"net/http"

	"configdatabase/src/auth/meta"
	"configdatabase/src/common"
	"configdatabase/src/common/auth"
	"configdatabase/src/common/blog"
	"configdatabase/src/common/metadata"
	"configdatabase/src/common/util"

	"github.com/emicklei/go-restful"
)

func (s *Service) InitAuthCenter(req *restful.Request, resp *restful.Response) {
	rHeader := req.Request.Header
	rid := util.GetHTTPCCRequestID(rHeader)
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(rHeader))
	if !auth.IsAuthed() {
		blog.Errorf("received auth center initialization request, but auth center not enabled, rid: %s", rid)
		result := &metadata.RespError{
			Msg: defErr.Error(common.CCErrCommAuthCenterIsNotEnabled),
		}
		resp.WriteError(http.StatusBadRequest, result)
		return
	}

	bizs := make([]metadata.BizInst, 0)
	bizFilter := map[string]interface{}{
		common.BKDefaultField: map[string]interface{}{
			common.BKDBNE: common.DefaultAppFlag,
		},
	}
	if err := s.db.Table(common.BKTableNameBaseApp).Find(bizFilter).All(s.ctx, &bizs); err != nil {
		blog.Errorf("init auth center failed, list businesses failed, err: %v, rid: %s", err, rid)
		result := &metadata.RespError{
			Msg: defErr.Errorf(common.CCErrCommInitAuthCenterFailed, err.Error()),
		}
		resp.WriteError(http.StatusInternalServerError, result)
		return
	}

	noRscPoolBiz := make([]metadata.BizInst, 0)
	resourcePoolNames := []string{"资源池", "resource pool"}
	for _, biz := range bizs {
		if util.InArray(biz.BizName, resourcePoolNames) {
			continue
		}
		noRscPoolBiz = append(noRscPoolBiz, biz)
	}

	cls := make([]metadata.Classification, 0)
	clsFilter := map[string]interface{}{
		common.BKIsPre: map[string]interface{}{
			common.BKDBNE: true,
		},
	}
	if err := s.db.Table(common.BKTableNameObjClassifiction).Find(clsFilter).All(s.ctx, &cls); err != nil {
		blog.Errorf("init auth center failed, list classifications failed, err: %+v, rid: %s", err, rid)
		result := &metadata.RespError{
			Msg: defErr.Errorf(common.CCErrCommInitAuthCenterFailed, err.Error()),
		}
		resp.WriteError(http.StatusInternalServerError, result)
		return
	}

	models := make([]metadata.Object, 0)
	modelFilter := map[string]interface{}{
		common.BKObjIDField: map[string]interface{}{
			common.BKDBNIN: []string{common.BKInnerObjIDProc, common.BKInnerObjIDPlat},
		},
	}
	if err := s.db.Table(common.BKTableNameObjDes).Find(modelFilter).All(s.ctx, &models); err != nil {
		blog.Errorf("init auth center failed, list models failed, err: %v, rid: %s", err, rid)
		result := &metadata.RespError{
			Msg: defErr.Errorf(common.CCErrCommInitAuthCenterFailed, err.Error()),
		}
		resp.WriteError(http.StatusInternalServerError, result)
		return
	}

	associationKinds := make([]metadata.AssociationKind, 0)
	associationFilter := map[string]interface{}{
		common.BKIsPre: true,
	}
	if err := s.db.Table(common.BKTableNameAsstDes).Find(associationFilter).All(s.ctx, &associationKinds); err != nil {
		blog.Errorf("init auth center with association kind failed, get details association kind failed, err: %+v, rid: %s", err, rid)
		result := &metadata.RespError{
			Msg: defErr.Errorf(common.CCErrCommInitAuthCenterFailed, err.Error()),
		}
		resp.WriteError(http.StatusInternalServerError, result)
		return
	}
	assoKinds := make([]metadata.AssociationKind, 0)
	for ak := range associationKinds {
		// filter bk_mainline kind, do not register to auth center.
		if associationKinds[ak].AssociationKindID != common.AssociationKindMainline {
			assoKinds = append(assoKinds, associationKinds[ak])
		}
	}

	initCfg := meta.InitConfig{
		Bizs:             noRscPoolBiz,
		Models:           models,
		Classifications:  cls,
		AssociationKinds: assoKinds,
	}
	if err := s.authCenter.Init(s.ctx, rHeader, initCfg); nil != err {
		blog.Errorf("init auth center failed, err: %+v, rid: %s", err, rid)
		result := &metadata.RespError{
			Msg: defErr.Errorf(common.CCErrCommInitAuthCenterFailed, err.Error()),
		}
		resp.WriteError(http.StatusInternalServerError, result)
		return
	}
	resp.WriteEntity(metadata.NewSuccessResp("init auth center success"))
}

// init new auth center
func (s *Service) InitAuthAccount(req *restful.Request, resp *restful.Response) {
	rHeader := req.Request.Header
	rid := util.GetHTTPCCRequestID(rHeader)
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(rHeader))
	if !auth.IsAuthed() {
		blog.Errorf("received auth center initialization request, but auth center not enabled, rid: %s", rid)
		result := &metadata.RespError{
			Msg: defErr.Error(common.CCErrCommAuthCenterIsNotEnabled),
		}
		resp.WriteError(http.StatusBadRequest, result)
		return
	}

	bizs := make([]metadata.BizInst, 0)
	bizFilter := map[string]interface{}{
		common.BKDefaultField: map[string]interface{}{
			common.BKDBNE: common.DefaultAppFlag,
		},
		// 不注册资源池、默认业务
		common.BKAppIDField: map[string]interface{}{
			common.BKDBNIN: []int{1, 2},
		},
	}
	if err := s.db.Table(common.BKTableNameBaseApp).Find(bizFilter).All(s.ctx, &bizs); err != nil {
		blog.Errorf("init auth center failed, list businesses failed, err: %v, rid: %s", err, rid)
		result := &metadata.RespError{
			Msg: defErr.Errorf(common.CCErrCommInitAuthCenterFailed, err.Error()),
		}
		resp.WriteError(http.StatusInternalServerError, result)
		return
	}

	noRscPoolBiz := make([]metadata.BizInst, 0)
	resourcePoolNames := []string{"资源池", "resource pool"}
	for _, biz := range bizs {
		if util.InArray(biz.BizName, resourcePoolNames) {
			continue
		}
		noRscPoolBiz = append(noRscPoolBiz, biz)
	}

	// search rbiz object instance
	rBizs := make([]metadata.RBizInst, 0)
	rBizFilter := map[string]interface{}{
		common.BKObjIDField: common.BKInnerObjIDRealBiz,
		// 过滤默认业务下 rbiz
		common.BKAppIDField: map[string]interface{}{
			common.BKDBNIN: []int{2},
		},
	}
	if err := s.db.Table(common.BKTableNameBaseInst).Find(rBizFilter).All(s.ctx, &rBizs); err != nil {
		blog.Errorf("init auth center failed, list rbiz failed, err: %v, rid: %s", err, rid)
		result := &metadata.RespError{
			Msg: defErr.Errorf(common.CCErrCommInitAuthCenterFailed, err.Error()),
		}
		resp.WriteError(http.StatusInternalServerError, result)
		return
	}

	initCfg := meta.InitConfig{
		Bizs:             noRscPoolBiz,
		RBizs:  		  rBizs,
	}
	if err := s.authCenter.Init(s.ctx, rHeader, initCfg); nil != err {
		blog.Errorf("init auth center failed, err: %+v, rid: %s", err, rid)
		result := &metadata.RespError{
			Msg: defErr.Errorf(common.CCErrCommInitAuthCenterFailed, err.Error()),
		}
		resp.WriteError(http.StatusInternalServerError, result)
		return
	}
	resp.WriteEntity(metadata.NewSuccessResp("init auth center success"))
}