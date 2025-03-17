package store

import (
	"context"
	"database/sql"
	"errors"
)

type Comment struct {
	ID        int64  `json:"id"`
	Content   string `json:"content"`
	Likes     uint32 `json:"likes"`
	Dislikes  uint32 `json:"dislikes"`
	MangaID   int64  `json:"manga_id"`
	ChapterID int64  `json:"chapter_id"`
	UserID    int64  `json:"user_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	User      User   `json:"-"`
}

type CommentStore struct {
	db *sql.DB
}

func (s *CommentStore) Create(ctx context.Context, comment *Comment) error {
	query := `
  INSERT INTO comments (content, manga_id, chapter_id, user_id)
  VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
  `

	ctx, cancel := context.WithTimeout(context.Background(), QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		comment.Content,
		comment.MangaID,
		comment.ChapterID,
		comment.UserID,
	).Scan(
		&comment.ID,
		&comment.CreatedAt,
		&comment.UpdatedAt,
	)

	return err
}

func (s *CommentStore) GetByMangaID(ctx context.Context, mangaID int64) (*[]Comment, error) {
	query := `
	SELECT 
		c.content, 
		c.likes, 
		c.dislikes, 
		c.manga_id, 
		c.user_id, 
		c.created_at, 
		c.updated_at,
		users.id,
		users.username
	FROM comments c
	JOIN users ON users.id = c.user_id
	WHERE c.manga_id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), QueryTimeoutDuration)
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

	var comments []Comment
	for rows.Next() {
		var comment Comment
		comment.User = User{}
		if err := rows.Scan(
			&comment.Content,
			&comment.Likes,
			&comment.Dislikes,
			&comment.MangaID,
			&comment.UserID,
			&comment.CreatedAt,
			&comment.UpdatedAt,
			&comment.User.ID,
			&comment.User.Username,
		); err != nil {
			return nil, err
		}

		comments = append(comments, comment)
	}

	return &comments, nil
}
