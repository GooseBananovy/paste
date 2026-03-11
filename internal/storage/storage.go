package storage

import (
	"context"
	"errors"

	"github.com/goosebananovy/paste/internal/model"
)

var ErrNotFound = errors.New("not found")

type Storage interface {
	Create(ctx context.Context, content string) (string, error)
	Get(ctx context.Context, ID string) (*model.Paste, error)
	Delete(ctx context.Context, ID string) error
}
