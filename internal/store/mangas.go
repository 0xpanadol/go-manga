package store

import (
	"context"
	"database/sql"
)

type Manga struct {
	ID          int64    `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Slug        string   `json:"slug"`
	Cover       string   `json:"cover"`
	Tags        []string `json:"tags"`
	Type        string   `json:"type"`
	Status      string   `json:"status"`
	Painter     string   `json:"painter"`
	UserID      int64    `json:"user_id"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
	Version     int      `json:"version"`
}

type MangaStore struct {
	db *sql.DB
}

func (s *MangaStore) Create(ctx context.Context, manga *Manga) error {
	return nil
}
