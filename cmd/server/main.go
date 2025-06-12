// Package main Order Service API
// @title Order API
// @version 1.0
// @description API for managing orders, items, and baskets
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"order/internal/auth"
	"order/internal/basket"
	"order/internal/item"
	ord "order/internal/order"
	"order/internal/router"
	usr "order/internal/user"
)

func main() {
	_ = godotenv.Load()
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN env not set")
	}
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}

	orderRepo := ord.NewPostgresRepository(db)
	orderSvc := ord.NewService(orderRepo)
	orderCtrl := ord.NewController(orderSvc)

	itemRepo := item.NewPostgresRepository(db)
	itemSvc := item.NewService(itemRepo)
	itemCtrl := item.NewController(itemSvc)

	basketRepo := basket.NewPostgresRepository(db)
	basketSvc := basket.NewService(basketRepo)
	basketCtrl := basket.NewController(basketSvc)

	userRepo := usr.NewPostgresRepository(db)
	userSvc := usr.NewService(userRepo)
	userCtrl := usr.NewController(userSvc)

	userRepo := usr.NewPostgresRepository(db)
	userSvc := usr.NewService(userRepo)
	userCtrl := usr.NewController(userSvc)

	secret := []byte(os.Getenv("JWT_SECRET"))
	if len(secret) == 0 {
		secret = []byte("secret")
	}
	authSvc := auth.NewService(secret)
	authCtrl := auth.NewController(authSvc, userSvc)

	// log to file
	f, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("log file: %v", err)
	}
	log.SetOutput(io.MultiWriter(os.Stdout, f))

	r := router.New(orderCtrl, itemCtrl, basketCtrl, secret, authCtrl, userCtrl)
	log.Println("server listening on :8089")
	log.Fatal(http.ListenAndServe(":8089", r))
}
