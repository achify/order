package item

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

// Controller exposes HTTP handlers for items.
type Controller struct {
	Service  *Service
	Validate *validator.Validate
}

func NewController(s *Service) *Controller {
	return &Controller{Service: s, Validate: validator.New()}
}

func (c *Controller) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/items", c.listItems).Methods("GET")
	r.HandleFunc("/items/{id}", c.getItem).Methods("GET")
	r.HandleFunc("/items", c.createItem).Methods("POST")
	r.HandleFunc("/items/{id}", c.patchItem).Methods("PATCH")
	r.HandleFunc("/items/{id}", c.deleteItem).Methods("DELETE")
}

// listItems godoc
// @Summary List items
// @Tags items
// @Produce json
// @Success 200 {array} item.Item
// @Router /items [get]
// @Security BearerAuth
func (c *Controller) listItems(w http.ResponseWriter, r *http.Request) {
	list, err := c.Service.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondJSON(w, http.StatusOK, list)
}

// getItem godoc
// @Summary Get item
// @Tags items
// @Produce json
// @Param id path string true "Item ID"
// @Success 200 {object} item.Item
// @Router /items/{id} [get]
// @Security BearerAuth
func (c *Controller) getItem(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	it, err := c.Service.Get(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if it == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	respondJSON(w, http.StatusOK, it)
}

// createItem godoc
// @Summary Create item
// @Tags items
// @Accept json
// @Produce json
// @Param item body ItemCreateDTO true "New item"
// @Success 201 {object} item.Item
// @Router /items [post]
// @Security BearerAuth
func (c *Controller) createItem(w http.ResponseWriter, r *http.Request) {
	var dto ItemCreateDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := c.Validate.Struct(dto); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	it, err := c.Service.Create(r.Context(), dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondJSON(w, http.StatusCreated, it)
}

// patchItem godoc
// @Summary Update item
// @Tags items
// @Accept json
// @Produce json
// @Param id path string true "Item ID"
// @Param item body ItemUpdateDTO true "Fields to update"
// @Success 200 {object} item.Item
// @Router /items/{id} [patch]
// @Security BearerAuth
func (c *Controller) patchItem(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var dto ItemUpdateDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	it, err := c.Service.Update(r.Context(), id, dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if it == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	respondJSON(w, http.StatusOK, it)
}

// deleteItem godoc
// @Summary Delete item
// @Tags items
// @Param id path string true "Item ID"
// @Success 204 {string} string "no content"
// @Router /items/{id} [delete]
// @Security BearerAuth
func (c *Controller) deleteItem(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
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
