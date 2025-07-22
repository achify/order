package rbac

import (
	"encoding/json"
	"net/http"

	"order/internal/auth"
)

// Controller exposes RBAC related handlers.
type Controller struct{}

func NewController() *Controller { return &Controller{} }

// AllowedEndpoints returns a list of endpoints accessible by the current user.
func (c *Controller) AllowedEndpoints(w http.ResponseWriter, r *http.Request) {
	roles := auth.RolesFromContext(r.Context())
	eps := AllowedEndpointsForRoles(roles)
	respondJSON(w, http.StatusOK, eps)
}

func respondJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
