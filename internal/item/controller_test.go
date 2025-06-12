package item

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateItemValidation(t *testing.T) {
	repo := newMock()
	svc := NewService(repo)
	ctrl := NewController(svc)

	cases := []struct {
		name string
		body string
		want int
	}{
		{"empty", `{}`, http.StatusUnprocessableEntity},
		{"missing name", `{"price":1,"category_id":"c"}`, http.StatusUnprocessableEntity},
		{"missing price", `{"name":"n","category_id":"c"}`, http.StatusUnprocessableEntity},
		{"missing category", `{"name":"n","price":1}`, http.StatusUnprocessableEntity},
		{"valid", `{"name":"n","price":1,"category_id":"c"}`, http.StatusCreated},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/items", strings.NewReader(tc.body))
			w := httptest.NewRecorder()
			ctrl.createItem(w, req)
			if w.Code != tc.want {
				t.Fatalf("want %d got %d", tc.want, w.Code)
			}
		})
	}
}
