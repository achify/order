package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

// Controller exposes HTTP handlers for deliveries.
type Controller struct {
	Service  *Service
	Validate *validator.Validate
}

func NewController(s *Service) *Controller {
	return &Controller{Service: s, Validate: validator.New()}
}

func (c *Controller) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/deliveries", c.list).Methods("GET")
	r.HandleFunc("/deliveries/{id}", c.get).Methods("GET")
	r.HandleFunc("/deliveries", c.create).Methods("POST")
	r.HandleFunc("/deliveries/{id}", c.update).Methods("PATCH")
	r.HandleFunc("/deliveries/{id}", c.delete).Methods("DELETE")
	r.HandleFunc("/locations/{provider}", c.locations).Methods("GET")
}

func (c *Controller) list(w http.ResponseWriter, r *http.Request) {
	list, err := c.Service.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondJSON(w, http.StatusOK, list)
}

func (c *Controller) get(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	d, err := c.Service.Get(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if d == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	respondJSON(w, http.StatusOK, d)
}

func (c *Controller) create(w http.ResponseWriter, r *http.Request) {
	var dto CreateDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := c.Validate.Struct(dto); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	d, err := c.Service.Create(r.Context(), dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondJSON(w, http.StatusCreated, d)
}

func (c *Controller) update(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var dto UpdateDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	d, err := c.Service.Update(r.Context(), id, dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if d == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	respondJSON(w, http.StatusOK, d)
}

func (c *Controller) delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := c.Service.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (c *Controller) locations(w http.ResponseWriter, r *http.Request) {
	provider := mux.Vars(r)["provider"]
	locs, err := c.Service.Repo.Locations(r.Context(), provider)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondJSON(w, http.StatusOK, locs)
}

func respondJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
