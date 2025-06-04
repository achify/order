package order

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateOrderValidation(t *testing.T) {
	repo := newMock()
	svc := NewService(repo)
	ctrl := NewController(svc)

	cases := []struct {
		name string
		body string
		want int
	}{
		{"empty body", `{}`, http.StatusUnprocessableEntity},
		{"missing receiver", `{"account_id":"a","seller_id":"s","delivery_id":"d","basket_id":"b"}`, http.StatusUnprocessableEntity},
		{"missing account", `{"receiver_id":"r","seller_id":"s","delivery_id":"d","basket_id":"b"}`, http.StatusUnprocessableEntity},
		{"missing seller", `{"receiver_id":"r","account_id":"a","delivery_id":"d","basket_id":"b"}`, http.StatusUnprocessableEntity},
		{"missing delivery", `{"receiver_id":"r","account_id":"a","seller_id":"s","basket_id":"b"}`, http.StatusUnprocessableEntity},
		{"missing basket", `{"receiver_id":"r","account_id":"a","seller_id":"s","delivery_id":"d"}`, http.StatusUnprocessableEntity},
		{"invalid json", `{"receiver_id":`, http.StatusBadRequest},
		{"valid", `{"receiver_id":"r","account_id":"a","seller_id":"s","delivery_id":"d","basket_id":"b"}`, http.StatusCreated},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/orders", strings.NewReader(tc.body))
			w := httptest.NewRecorder()
			ctrl.createOrder(w, req)
			if w.Code != tc.want {
				t.Fatalf("want %d got %d body: %s", tc.want, w.Code, w.Body.String())
			}
		})
	}
}
