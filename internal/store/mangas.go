package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type Manga struct {
	ID          int64       `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Slug        string      `json:"slug"`
	Cover       string      `json:"cover"`
	Tags        []string    `json:"tags"`
	Type        int64       `json:"type"`
	Status      int64       `json:"status"`
	Painter     interface{} `json:"painter"`
	UserID      int64       `json:"user_id"`
	CreatedAt   string      `json:"created_at"`
	UpdatedAt   string      `json:"updated_at"`
	Comments    []Comment   `json:"comments"`
	Version     int         `json:"version"`
}

type MangaStore struct {
	db *sql.DB
}

func (s *MangaStore) Create(ctx context.Context, manga *Manga) error {
	query := `
	INSERT INTO mangas (title, description, slug, cover, tags, user_id)
	VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

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

func (s *MangaStore) GetByID(ctx context.Context, mangaID int64) (*Manga, error) {
	query := `
	SELECT
		id,
		title,
		description,
		slug,
		cover,
		tags, 
		type, 
		status, 
		painter, 
		user_id, 
		created_at, 
		updated_at, 
		version
	FROM mangas WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var manga Manga
	err := s.db.QueryRowContext(ctx, query, mangaID).Scan(
		&manga.ID,
		&manga.Title,
		&manga.Description,
		&manga.Slug,
		&manga.Cover,
		pq.Array(&manga.Tags),
		&manga.Type,
		&manga.Status,
		&manga.Painter,
		&manga.UserID,
		&manga.CreatedAt,
		&manga.UpdatedAt,
		&manga.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &manga, nil
}
