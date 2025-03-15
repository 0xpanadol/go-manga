package store

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type Manga struct {
	ID          int64    `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Slug        string   `json:"slug"`
	Cover       string   `json:"cover"`
	Tags        []string `json:"tags"`
	Type        int64    `json:"type"`
	Status      int64    `json:"status"`
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
	query := `
	INSERT INTO mangas (title, description, slug, cover, tags, user_id)
	VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at
	`

	err := s.db.QueryRowContext(
		ctx,
		query,
		manga.Title,
		manga.Description,
		manga.Slug,
		manga.Cover,
		pq.Array(manga.Tags),
		manga.UserID,
	).Scan(
		&manga.ID,
		&manga.CreatedAt,
		&manga.UpdatedAt,
	)

	return err
}
