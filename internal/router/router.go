package router

import (
    "net/http"

    "github.com/gorilla/mux"

    "order/internal/auth"
    ord "order/internal/order"
    "order/internal/metrics"
)

// New sets up application routes with middleware
func New(c *ord.Controller, secret []byte) http.Handler {
    r := mux.NewRouter()

    r.Use(metrics.Middleware)
    r.Use(auth.Middleware(secret))

    c.RegisterRoutes(r)
    r.HandleFunc("/metrics", metrics.Handler)
    return r
}
