package account

import (
	"configdatabase/src/auth/meta"
	"context"
)

type authFilter struct {

}

// Register resource filter
func (af *authFilter) needRegisterResource(ctx context.Context, rs meta.ResourceAttribute) bool {
	//if rs.Type == meta.Business {
	//	return true
	//}

	return false
}