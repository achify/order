package router

import (
	"net/http"

	"github.com/gorilla/mux"

	"order/internal/auth"
	"order/internal/metrics"
	ord "order/internal/order"
)

// New sets up application routes with middleware
func New(orderCtrl *ord.Controller, secret []byte, authCtrl *auth.Controller) http.Handler {
	r := mux.NewRouter()

	r.Use(metrics.Middleware)

	// public endpoints
	r.HandleFunc("/auth/login", authCtrl.Login).Methods("POST")
	r.HandleFunc("/auth/refresh", authCtrl.Refresh).Methods("POST")
	r.HandleFunc("/metrics", metrics.Handler)

	// protected routes
	api := r.PathPrefix("").Subrouter()
	api.Use(auth.Middleware(secret))
	orderCtrl.RegisterRoutes(api)

	return r
}
