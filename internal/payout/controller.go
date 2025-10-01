package payout

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"

	"order/internal/transfer"
)

// Controller exposes payout HTTP handlers.
type Controller struct {
	service  *Service
	validate *validator.Validate
}

// NewController constructs a payout controller with registered validations.
func NewController(service *Service) *Controller {
	v := validator.New()
	RegisterValidations(v)
	return &Controller{service: service, validate: v}
}

// RegisterRoutes mounts payout routes under the router.
func (c *Controller) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/payouts", c.create).Methods(http.MethodPost)
}

func (c *Controller) create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := c.validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	payout, err := c.service.Initiate(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, transfer.ErrRateNotFound):
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	respondJSON(w, http.StatusCreated, payout)
}

func respondJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
