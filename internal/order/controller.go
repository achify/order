package order

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

// Controller exposes HTTP handlers

type Controller struct {
	Service  *Service
	Validate *validator.Validate
}

func NewController(s *Service) *Controller {
	return &Controller{Service: s, Validate: validator.New()}
}

func (c *Controller) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/orders", c.listOrders).Methods("GET")
	r.HandleFunc("/orders/{id}", c.getOrder).Methods("GET")
	r.HandleFunc("/orders", c.createOrder).Methods("POST")
	r.HandleFunc("/orders/{id}", c.patchOrder).Methods("PATCH")
	r.HandleFunc("/orders/{id}", c.deleteOrder).Methods("DELETE")
}

func (c *Controller) listOrders(w http.ResponseWriter, r *http.Request) {
	deliveryID := r.URL.Query().Get("delivery_id")
	list, err := c.Service.List(r.Context(), deliveryID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var filtered []Order
	for i := range list {
		if canView(r.Context(), &list[i]) {
			filtered = append(filtered, list[i])
		}
	}
	respondJSON(w, http.StatusOK, filtered)
}

func (c *Controller) getOrder(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	o, err := c.Service.Get(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if o == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if !canView(r.Context(), o) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	respondJSON(w, http.StatusOK, o)
}

func (c *Controller) createOrder(w http.ResponseWriter, r *http.Request) {
	var dto OrderCreateDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := c.Validate.Struct(dto); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	o, err := c.Service.Create(r.Context(), dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondJSON(w, http.StatusCreated, o)
}

func (c *Controller) patchOrder(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var dto OrderUpdateDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	o, err := c.Service.Get(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if o == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if !canEdit(r.Context(), o) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	o, err = c.Service.Update(r.Context(), id, dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if o == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	respondJSON(w, http.StatusOK, o)
}

func (c *Controller) deleteOrder(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	o, err := c.Service.Get(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if o == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if !canEdit(r.Context(), o) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	if err := c.Service.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func respondJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
