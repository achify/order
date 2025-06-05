package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"order/internal/delivery"
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

	repo := delivery.NewPostgresRepository(db)
	svc := delivery.NewService(repo)
	ctrl := delivery.NewController(svc)

	f, err := os.OpenFile("delivery.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("log file: %v", err)
	}
	log.SetOutput(io.MultiWriter(os.Stdout, f))

	r := mux.NewRouter()
	ctrl.RegisterRoutes(r)

	// daily sync of provider locations
	go func() {
		for {
			if err := svc.SyncLocations(context.Background()); err != nil {
				log.Printf("sync locations: %v", err)
			}
			time.Sleep(24 * time.Hour)
		}
	}()

	log.Println("delivery service listening on :8090")
	log.Fatal(http.ListenAndServe(":8090", r))
}
