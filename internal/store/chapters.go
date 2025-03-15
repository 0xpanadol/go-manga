package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type Chapter struct {
	ID        int64    `json:"id"`
	Title     string   `json:"title"`
	Content   []string `json:"content"`
	Number    float64  `json:"chapter_number"`
	MangaID   int64    `json:"manga_id"`
	UserID    int64    `json:"user_id"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
	Version   int      `json:"version"`
}

type ChapterStore struct {
	db *sql.DB
}

func (s *ChapterStore) Create(ctx context.Context, chapter *Chapter) error {
	query := `
  INSERT INTO chapters (title, content, chapter_number, manga_id, user_id)
  VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at
  `

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		chapter.Title,
		pq.Array(chapter.Content),
		chapter.Number,
		chapter.MangaID,
		chapter.UserID,
	).Scan(
		&chapter.ID,
		&chapter.CreatedAt,
		&chapter.UpdatedAt,
	)

	return err
}

func (s *ChapterStore) GetByMangaID(ctx context.Context, mangaID int64) (*[]Chapter, error) {
	query := `
  SELECT 
    c.id,
    c.title,
    c.content,
    c.chapter_number,
    c.manga_id,
    c.user_id,
    c.created_at,
    c.updated_at,
    c.version,
    users.id,
    users.username
  FROM chapters c
  JOIN users ON c.user_id = users.id
  WHERE c.manga_id = $1
  ORDER BY c.chapter_number DESC;
  `

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, mangaID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	chapters := []Chapter{}
	for rows.Next() {
		var chapter Chapter
		if err := rows.Scan(
			&chapter.ID,
			&chapter.Content,
			&chapter.Number,
			&chapter.MangaID,
			&chapter.UserID,
			&chapter.CreatedAt,
			&chapter.UpdatedAt,
			&chapter.Version,
		); err != nil {
			return nil, err
		}

		chapters = append(chapters, chapter)
	}

	return &chapters, nil
}

func (s *ChapterStore) Delete(ctx context.Context, chapterID int64, mangaID int64) error {
	query := `DELETE FROM chapters WHERE id = $1 AND manga_id = $2`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	result, err := s.db.ExecContext(ctx, query, chapterID, mangaID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
