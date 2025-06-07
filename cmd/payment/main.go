package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	ord "order/internal/order"
	"order/internal/payment"
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

	payRepo := payment.NewPostgresRepository(db)
	orderRepo := ord.NewPostgresRepository(db)
	svc := payment.NewService(payRepo, orderRepo)
	ctrl := payment.NewController(svc)

	f, err := os.OpenFile("payment.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("log file: %v", err)
	}
	log.SetOutput(io.MultiWriter(os.Stdout, f))

	r := mux.NewRouter()
	ctrl.RegisterRoutes(r)

	log.Println("payment service listening on :8091")
	log.Fatal(http.ListenAndServe(":8091", r))
}
