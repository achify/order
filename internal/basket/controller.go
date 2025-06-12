package basket

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

// Controller exposes HTTP handlers for baskets.
type Controller struct {
	Service  *Service
	Validate *validator.Validate
}

func NewController(s *Service) *Controller {
	return &Controller{Service: s, Validate: validator.New()}
}

func (c *Controller) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/baskets", c.listBaskets).Methods("GET")
	r.HandleFunc("/baskets/{id}", c.getBasket).Methods("GET")
	r.HandleFunc("/baskets", c.createBasket).Methods("POST")
	r.HandleFunc("/baskets/{id}", c.patchBasket).Methods("PATCH")
	r.HandleFunc("/baskets/{id}", c.deleteBasket).Methods("DELETE")
	r.HandleFunc("/baskets/{id}/items", c.addItem).Methods("POST")
	r.HandleFunc("/baskets/{id}/items/{item_id}", c.patchItem).Methods("PATCH")
	r.HandleFunc("/baskets/{id}/items/{item_id}", c.deleteItem).Methods("DELETE")
	r.HandleFunc("/baskets/{id}/items", c.listItems).Methods("GET")
}

func (c *Controller) listBaskets(w http.ResponseWriter, r *http.Request) {
	list, err := c.Service.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondJSON(w, http.StatusOK, list)
}

func (c *Controller) getBasket(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	b, err := c.Service.Get(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if b == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	respondJSON(w, http.StatusOK, b)
}

func (c *Controller) createBasket(w http.ResponseWriter, r *http.Request) {
	var dto BasketCreateDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := c.Validate.Struct(dto); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	b, err := c.Service.Create(r.Context(), dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondJSON(w, http.StatusCreated, b)
}

func (c *Controller) patchBasket(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var dto BasketUpdateDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	b, err := c.Service.Update(r.Context(), id, dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if b == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	respondJSON(w, http.StatusOK, b)
}

func (c *Controller) deleteBasket(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := c.Service.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (c *Controller) addItem(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var dto ItemDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := c.Validate.Struct(dto); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	if err := c.Service.AddItem(r.Context(), id, dto); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (c *Controller) patchItem(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	itemID := mux.Vars(r)["item_id"]
	var dto ItemDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := c.Service.UpdateItem(r.Context(), id, itemID, dto); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (c *Controller) deleteItem(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	itemID := mux.Vars(r)["item_id"]
	if err := c.Service.DeleteItem(r.Context(), id, itemID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (c *Controller) listItems(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	items, err := c.Service.ListItems(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func respondJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
