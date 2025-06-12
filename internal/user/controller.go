package user

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

// Controller exposes handlers for managing users.
type Controller struct {
	Service  *Service
	Validate *validator.Validate
}

func NewController(s *Service) *Controller { return &Controller{Service: s, Validate: validator.New()} }

// CreateUser registers a new user.
func (c *Controller) CreateUser(w http.ResponseWriter, r *http.Request) {
	var dto UserCreateDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := c.Validate.Struct(dto); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	u, err := c.Service.Create(r.Context(), dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondJSON(w, http.StatusCreated, u)
}

func respondJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
