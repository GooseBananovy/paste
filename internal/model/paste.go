package model

import "time"

type Paste struct {
	ID        string
	Content   string
	CreatedAt time.Time
	ExpiresAt *time.Time
}

