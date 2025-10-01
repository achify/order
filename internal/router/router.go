package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/swaggo/http-swagger/v2"
	_ "order/internal/docs"

	"order/internal/auth"
	"order/internal/basket"
	"order/internal/item"
	"order/internal/metrics"
	ord "order/internal/order"
	"order/internal/payout"
	"order/internal/transfer"
	usr "order/internal/user"
)

// New sets up application routes with middleware
func New(orderCtrl *ord.Controller, itemCtrl *item.Controller, basketCtrl *basket.Controller, secret []byte, authCtrl *auth.Controller, userCtrl *usr.Controller, transferCtrl *transfer.Controller, payoutCtrl *payout.Controller) http.Handler {
	r := mux.NewRouter()

	r.Use(metrics.Middleware)

	// public endpoints
	r.HandleFunc("/auth/login", authCtrl.Login).Methods("POST")
	r.HandleFunc("/auth/refresh", authCtrl.Refresh).Methods("POST")
	r.HandleFunc("/metrics", metrics.Handler)
	r.PathPrefix("/docs/").Handler(httpSwagger.WrapHandler)

	// protected routes
	api := r.PathPrefix("").Subrouter()
	api.Use(auth.Middleware(secret))
	orderCtrl.RegisterRoutes(api)
	itemCtrl.RegisterRoutes(api)
	basketCtrl.RegisterRoutes(api)
	transferCtrl.RegisterRoutes(api)
	payoutCtrl.RegisterRoutes(api)
	api.HandleFunc("/users", userCtrl.CreateUser).Methods("POST")

	return r
}
