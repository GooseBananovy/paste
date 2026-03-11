package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/goosebananovy/paste/internal/storage"
)

type PasteHandler struct {
	stg storage.Storage
}

type createResponse struct {
	ID string `json:"id"`
}

type getResponse struct {
	Content string `json:"content"`
}

func NewPasteHandler(stg storage.Storage) *PasteHandler {
	return &PasteHandler{
		stg: stg,
	}
}

func (h *PasteHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	content := string(body)

	if len(content) == 0 {
		http.Error(w, "empty content in request", http.StatusBadRequest)
		return
	}

	ID, err := h.stg.Create(ctx, content)
	if err != nil {
		log.Printf("failed to create: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(createResponse{ID: ID})
}

func (h *PasteHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ID := r.PathValue("id")

	if len(ID) == 0 {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	paste, err := h.stg.Get(ctx, ID)

	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			http.Error(w, "paste not found", http.StatusNotFound)
		} else {
			log.Printf("failed to delete: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(getResponse{Content: paste.Content})
}

func (h *PasteHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ID := r.PathValue("id")

	if len(ID) == 0 {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	if err := h.stg.Delete(ctx, ID); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			http.Error(w, "paste not found", http.StatusNotFound)
		} else {
			log.Printf("failed to get: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
