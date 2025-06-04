package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	ord "order/internal/order"
	"order/internal/router"
)

func main() {
	_ = godotenv.Load()
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5552/order?sslmode=disable"
	}
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}

	repo := ord.NewPostgresRepository(db)
	svc := ord.NewService(repo)
	ctrl := ord.NewController(svc)

	secret := []byte(os.Getenv("JWT_SECRET"))
	if len(secret) == 0 {
		secret = []byte("secret")
	}

	// log to file
	f, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("log file: %v", err)
	}
	log.SetOutput(io.MultiWriter(os.Stdout, f))

	r := router.New(ctrl, secret)
	log.Println("server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
