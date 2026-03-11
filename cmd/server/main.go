package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/goosebananovy/paste/internal/handler"
	"github.com/goosebananovy/paste/internal/storage"
)

func main() {
	ms := storage.NewMemoryStorage()

	ph := handler.NewPasteHandler(ms)
	
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
