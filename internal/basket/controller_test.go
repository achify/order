package basket

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateBasketValidation(t *testing.T) {
	repo := newMock()
	svc := NewService(repo)
	ctrl := NewController(svc)

	cases := []struct {
		name string
		body string
		want int
	}{
		{"empty", `{}`, http.StatusUnprocessableEntity},
		{"missing account", `{"total_price":1}`, http.StatusUnprocessableEntity},
		{"valid", `{"account_id":"a"}`, http.StatusCreated},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/baskets", strings.NewReader(tc.body))
			w := httptest.NewRecorder()
			ctrl.createBasket(w, req)
			if w.Code != tc.want {
				t.Fatalf("want %d got %d", tc.want, w.Code)
			}
		})
	}
}
