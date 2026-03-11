package storage

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/goosebananovy/paste/internal/model"
)

type MemoryStorage struct {
	pool map[string]*model.Paste
	mu   sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		pool: make(map[string]*model.Paste),
	}
}

func (ms *MemoryStorage) Create(ctx context.Context, content string) (ID string, err error) {
	rawID := make([]byte, 4)

	if _, err = rand.Read(rawID); err != nil {
		return "", fmt.Errorf("falied to generate ID: %w", err)
	}

	ID = hex.EncodeToString(rawID)

	ms.mu.Lock()
	ms.pool[ID] = &model.Paste{
		ID:        ID,
		Content:   content,
		CreatedAt: time.Now(),
		ExpiresAt: nil,
	}
	ms.mu.Unlock()

	return ID, nil
}

func (ms *MemoryStorage) Get(ctx context.Context, ID string) (res *model.Paste, err error) {
	ok := true

	ms.mu.RLock()
	defer ms.mu.RUnlock()

	res, ok = ms.pool[ID]

	if !ok {
		return nil, ErrNotFound
	}

	return res, nil
}

func (ms *MemoryStorage) Delete(ctx context.Context, ID string) (err error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if _, ok := ms.pool[ID]; !ok {
		return ErrNotFound
	}

	delete(ms.pool, ID)

	return nil
}
