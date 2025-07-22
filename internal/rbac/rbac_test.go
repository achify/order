package rbac

import (
	"testing"
)

func TestAllowedEndpointsForRoles(t *testing.T) {
	eps := AllowedEndpointsForRoles([]string{"admin"})
	if len(eps) == 0 {
		t.Fatalf("no endpoints returned")
	}
	// expect one known endpoint
	found := false
	for _, e := range eps {
		if e.Method == "GET" && e.Path == "/orders" && e.Permission == string(PermissionViewOrders) {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected orders GET endpoint for admin")
	}
}
