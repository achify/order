package order

import (
	"context"

	"order/internal/auth"
	"order/internal/user"
)

func hasRole(roles []string, r user.Role) bool {
	for _, v := range roles {
		if v == string(r) {
			return true
		}
	}
	return false
}

func canView(ctx context.Context, o *Order) bool {
	roles := auth.RolesFromContext(ctx)
	uid := auth.UserIDFromContext(ctx)
	if hasRole(roles, user.RoleAdmin) || hasRole(roles, user.RoleSupportManager) || hasRole(roles, user.RoleAuditor) {
		return true
	}
	if hasRole(roles, user.RoleSeller) && o.SellerID == uid {
		return true
	}
	if hasRole(roles, user.RoleCustomer) && o.AccountID == uid {
		return true
	}
	return false
}

func canEdit(ctx context.Context, o *Order) bool {
	roles := auth.RolesFromContext(ctx)
	uid := auth.UserIDFromContext(ctx)
	if hasRole(roles, user.RoleAdmin) || hasRole(roles, user.RoleSupportManager) {
		return true
	}
	if hasRole(roles, user.RoleSeller) && o.SellerID == uid {
		return true
	}
	return false
}
