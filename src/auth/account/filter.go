package account

import (
	"configdatabase/src/auth/meta"
	"context"
)

type authFilter struct {

}

// Register a resource instance
func (af *authFilter) registerResourceFilter(ctx context.Context, rs meta.ResourceAttribute) bool {
	if rs.Type == meta.Business {
		return false
	}

	return true
}