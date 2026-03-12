package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/goosebananovy/paste/internal/handler"
	"github.com/goosebananovy/paste/internal/storage"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	connString := os.Getenv("DATABASE_URL")
	if len(connString) == 0 {
		log.Fatal("DATABASE_URL is not set")
	}

	ctx := context.Background()
	ps, err := storage.NewPostgresStorage(ctx, connString)
	if err != nil {
		log.Fatalf("failed to create connection: %v", err)
	}

	ph := handler.NewPasteHandler(ps)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		type healthResponse struct {
			Status string `json:"status"`
		}

		response := healthResponse{
			Status: "ok",
		}
		json.NewEncoder(w).Encode(response)
	})

	mux.HandleFunc("POST /paste", ph.Create)
	mux.HandleFunc("GET /paste/{id}", ph.Get)
	mux.HandleFunc("DELETE /paste/{id}", ph.Delete)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
