package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/goosebananovy/paste/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStorage struct {
	pool *pgxpool.Pool
}

func NewPostgresStorage(ctx context.Context, connString string) (*PostgresStorage, error) {
	pool, err := pgxpool.New(ctx, connString)

	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	return &PostgresStorage{pool: pool}, nil
}

func (ps *PostgresStorage) Create(ctx context.Context, content string) (ID string, err error) {
	if err = ps.pool.QueryRow(ctx, "INSERT INTO pastes (content, created_at) VALUES ($1, $2) RETURNING id", content, time.Now()).Scan(&ID); err != nil {
		return "", fmt.Errorf("failed to create paste: %w", err)
	}

	return ID, nil
}

func (ps *PostgresStorage) Get(ctx context.Context, ID string) (*model.Paste, error) {
	var paste model.Paste

	if err := ps.pool.QueryRow(ctx, "SELECT id, content, created_at, expires_at FROM pastes WHERE id = $1", ID).Scan(&paste.ID, &paste.Content, &paste.CreatedAt, &paste.ExpiresAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("failed to get paste: %w", err)
	}

	return &paste, nil
}

func (ps *PostgresStorage) Delete(ctx context.Context, ID string) error {
	tag, err := ps.pool.Exec(ctx, "DELETE FROM pastes WHERE id = $1", ID)

	if err != nil {
		return fmt.Errorf("failed to delete paste: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}
