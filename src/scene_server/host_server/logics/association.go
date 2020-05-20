package logics

import (
	"configdatabase/src/common"
	"configdatabase/src/common/blog"
	"configdatabase/src/common/mapstr"
	"configdatabase/src/common/metadata"
	hutil "configdatabase/src/scene_server/host_server/util"
	"context"
)

func (lgc *Logics) SearchInstAssociation(ctx context.Context, input *metadata.QueryCondition) ([]metadata.InstAsst, error) {
	rsp, err := lgc.CoreAPI.CoreService().Association().ReadInstAssociation(ctx, lgc.header, input)
	if nil != err {
		blog.Errorf("[search-asst] failed to request object controller, err: %s, rid: %s", err.Error(), lgc.rid)
		return nil, lgc.ccErr.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !rsp.Result {
		blog.Errorf("[search-asst] failed to search host association info, query: %#v, err: %s, rid: %s", input, rsp.ErrMsg, lgc.rid)
		return nil, lgc.ccErr.New(rsp.Code, rsp.ErrMsg)
	}

	return rsp.Data.Info, nil
}

func (lgc *Logics) GetAllHostAssociation(ctx context.Context, iHostIDArr []int64) ([]metadata.InstAsst, error) {
	cond := mapstr.MapStr{}
	multiOrOper := make([]map[string]interface{}, 0)
	objCond := hutil.NewOperation().WithObjID(common.BKInnerObjIDHost).
		WithInstIDx(map[string]interface{}{common.BKDBIN: iHostIDArr}).Data()
	asstObjCond := hutil.NewOperation().WithAssoObjID(common.BKInnerObjIDHost).
		WithAssoInstID(map[string]interface{}{common.BKDBIN: iHostIDArr}).Data()
	multiOrOper = append(multiOrOper, objCond, asstObjCond)
	cond[common.BKDBOR] = multiOrOper

	query := &metadata.QueryCondition{
		Condition: cond,
	}
	return lgc.SearchInstAssociation(ctx, query)
}