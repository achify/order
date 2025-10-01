package transfer

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

// Controller wires HTTP transport with the transfer service.
type Controller struct {
	service  *Service
	validate *validator.Validate
}

// NewController builds a transfer controller with opinionated validations.
func NewController(service *Service) *Controller {
	v := validator.New()
	RegisterValidations(v)
	return &Controller{service: service, validate: v}
}

// RegisterRoutes registers transfer endpoints under the provided router.
func (c *Controller) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/transfers/calculate", c.calculate).Methods(http.MethodPost)
}

func (c *Controller) calculate(w http.ResponseWriter, r *http.Request) {
	var req CalculateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := c.validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	quote, err := c.service.Calculate(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, ErrRateNotFound):
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	respondJSON(w, http.StatusOK, quote)
}

func respondJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
