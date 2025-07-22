package rbac

import (
	"order/internal/user"
)

// Permission represents an action that can be performed on a resource.
type Permission string

const (
	PermissionViewOrders     Permission = "orders:view"
	PermissionEditOrders     Permission = "orders:edit"
	PermissionViewItems      Permission = "items:view"
	PermissionEditItems      Permission = "items:edit"
	PermissionViewBaskets    Permission = "baskets:view"
	PermissionEditBaskets    Permission = "baskets:edit"
	PermissionViewPayments   Permission = "payments:view"
	PermissionEditPayments   Permission = "payments:edit"
	PermissionViewDeliveries Permission = "deliveries:view"
	PermissionEditDeliveries Permission = "deliveries:edit"
	PermissionCreateUsers    Permission = "users:create"
)

// Endpoint describes an HTTP endpoint that requires a permission.
type Endpoint struct {
	Method     string `json:"method"`
	Path       string `json:"path"`
	Permission string `json:"permission"`
}

var rolePermissions = map[user.Role][]Permission{
	user.RoleAdmin: {
		PermissionViewOrders, PermissionEditOrders,
		PermissionViewItems, PermissionEditItems,
		PermissionViewBaskets, PermissionEditBaskets,
		PermissionViewPayments, PermissionEditPayments,
		PermissionViewDeliveries, PermissionEditDeliveries,
		PermissionCreateUsers,
	},
	user.RoleSupportManager: {
		PermissionViewOrders, PermissionEditOrders,
		PermissionViewItems, PermissionEditItems,
		PermissionViewBaskets, PermissionEditBaskets,
		PermissionViewPayments, PermissionEditPayments,
		PermissionViewDeliveries, PermissionEditDeliveries,
		PermissionCreateUsers,
	},
	user.RoleSeller: {
		PermissionViewOrders, PermissionEditOrders,
		PermissionViewItems,
		PermissionViewBaskets, PermissionEditBaskets,
	},
	user.RoleCustomer: {
		PermissionViewOrders,
		PermissionViewBaskets,
	},
	user.RoleLogisticsManager: {
		PermissionViewDeliveries, PermissionEditDeliveries,
	},
	user.RolePaymentManager: {
		PermissionViewPayments, PermissionEditPayments,
	},
	user.RoleAuditor: {
		PermissionViewOrders,
	},
}

var permissionEndpoints = map[Permission][]Endpoint{
	PermissionViewOrders: {
		{Method: "GET", Path: "/orders", Permission: string(PermissionViewOrders)},
		{Method: "GET", Path: "/orders/{id}", Permission: string(PermissionViewOrders)},
	},
	PermissionEditOrders: {
		{Method: "POST", Path: "/orders", Permission: string(PermissionEditOrders)},
		{Method: "PATCH", Path: "/orders/{id}", Permission: string(PermissionEditOrders)},
		{Method: "DELETE", Path: "/orders/{id}", Permission: string(PermissionEditOrders)},
	},
	PermissionViewItems: {
		{Method: "GET", Path: "/items", Permission: string(PermissionViewItems)},
		{Method: "GET", Path: "/items/{id}", Permission: string(PermissionViewItems)},
	},
	PermissionEditItems: {
		{Method: "POST", Path: "/items", Permission: string(PermissionEditItems)},
		{Method: "PATCH", Path: "/items/{id}", Permission: string(PermissionEditItems)},
		{Method: "DELETE", Path: "/items/{id}", Permission: string(PermissionEditItems)},
	},
	PermissionViewBaskets: {
		{Method: "GET", Path: "/baskets", Permission: string(PermissionViewBaskets)},
		{Method: "GET", Path: "/baskets/{id}", Permission: string(PermissionViewBaskets)},
		{Method: "GET", Path: "/baskets/{id}/items", Permission: string(PermissionViewBaskets)},
	},
	PermissionEditBaskets: {
		{Method: "POST", Path: "/baskets", Permission: string(PermissionEditBaskets)},
		{Method: "PATCH", Path: "/baskets/{id}", Permission: string(PermissionEditBaskets)},
		{Method: "DELETE", Path: "/baskets/{id}", Permission: string(PermissionEditBaskets)},
		{Method: "POST", Path: "/baskets/{id}/items", Permission: string(PermissionEditBaskets)},
		{Method: "PATCH", Path: "/baskets/{id}/items/{item_id}", Permission: string(PermissionEditBaskets)},
		{Method: "DELETE", Path: "/baskets/{id}/items/{item_id}", Permission: string(PermissionEditBaskets)},
	},
	PermissionViewPayments: {
		{Method: "GET", Path: "/payments/{id}", Permission: string(PermissionViewPayments)},
	},
	PermissionEditPayments: {
		{Method: "POST", Path: "/payments", Permission: string(PermissionEditPayments)},
	},
	PermissionViewDeliveries: {
		{Method: "GET", Path: "/deliveries", Permission: string(PermissionViewDeliveries)},
		{Method: "GET", Path: "/deliveries/{id}", Permission: string(PermissionViewDeliveries)},
		{Method: "GET", Path: "/locations/{provider}", Permission: string(PermissionViewDeliveries)},
	},
	PermissionEditDeliveries: {
		{Method: "POST", Path: "/deliveries", Permission: string(PermissionEditDeliveries)},
		{Method: "PATCH", Path: "/deliveries/{id}", Permission: string(PermissionEditDeliveries)},
		{Method: "DELETE", Path: "/deliveries/{id}", Permission: string(PermissionEditDeliveries)},
	},
	PermissionCreateUsers: {
		{Method: "POST", Path: "/users", Permission: string(PermissionCreateUsers)},
	},
}

// AllowedEndpointsForRoles returns allowed endpoints for the specified roles.
func AllowedEndpointsForRoles(roles []string) []Endpoint {
	permSet := map[Permission]struct{}{}
	for _, r := range roles {
		role := user.Role(r)
		if perms, ok := rolePermissions[role]; ok {
			for _, p := range perms {
				permSet[p] = struct{}{}
			}
		}
	}
	endpMap := map[string]Endpoint{}
	for p := range permSet {
		for _, e := range permissionEndpoints[p] {
			key := e.Method + " " + e.Path
			endpMap[key] = e
		}
	}
	result := make([]Endpoint, 0, len(endpMap))
	for _, e := range endpMap {
		result = append(result, e)
	}
	return result
}
