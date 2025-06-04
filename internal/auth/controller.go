package auth

import (
	"encoding/json"
	"net/http"
)

// Controller exposes auth handlers to obtain and refresh tokens.
type Controller struct {
	Service *Service
}

func NewController(s *Service) *Controller { return &Controller{Service: s} }

// Login authenticates static credentials and returns JWT tokens.
func (c *Controller) Login(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if creds.Username != "admin" || creds.Password != "password" {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	token, refresh, err := c.Service.GenerateToken(creds.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{
		"token":         token,
		"refresh_token": refresh,
	})
}

// Refresh exchanges a refresh token for new tokens.
func (c *Controller) Refresh(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	token, refresh, err := c.Service.Refresh(req.RefreshToken)
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{
		"token":         token,
		"refresh_token": refresh,
	})
}

func respondJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
