package extensions

import (
	"configdatabase/src/auth/meta"
	"context"
	"strconv"
)

// add by elias 07/22
// 查询 rbiz的授权列表
func (am *AuthManager) ListAuthorizedRBizIDs(ctx context.Context, username string) ([]int64, error) {
	authorizedResources, err := am.Authorize.ListAuthorizedResources(ctx, username, 0, meta.RBusiness, meta.FindMany)
	if err != nil {
		return nil, err
	}

	authorizedIDs := make([]int64, 0)
	for _, app := range authorizedResources {
		id, err := strconv.ParseInt(app.ResourceID, 10, 64)
		if err != nil {
			return authorizedIDs, err
		}
		authorizedIDs = append(authorizedIDs, id)
	}
	return authorizedIDs, nil
}
