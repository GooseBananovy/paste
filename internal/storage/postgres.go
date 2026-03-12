package storage

import (
	"context"
	"errors"
	"fmt"
	"log"
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

	ps := &PostgresStorage{pool: pool}

	go ps.cleanup(context.Background())

	return ps, nil
}

func (ps *PostgresStorage) cleanup(ctx context.Context) {
	for {
		time.Sleep(1 * time.Minute)
		if _, err := ps.pool.Exec(ctx, "DELETE FROM pastes WHERE expires_at IS NOT NULL AND expires_at < NOW()"); err != nil {
			log.Printf("failed to clean pastes: %v", err)
		}
	}
}

func (ps *PostgresStorage) Create(ctx context.Context, content string, ttl *time.Duration) (ID string, err error) {

	var expiresAt *time.Time
	if ttl != nil {
		t := time.Now().Add(*ttl)
		expiresAt = &t
	}

	if err = ps.pool.QueryRow(ctx, "INSERT INTO pastes (content, created_at, expires_at) VALUES ($1, $2, $3) RETURNING id", content, time.Now(), expiresAt).Scan(&ID); err != nil {
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
