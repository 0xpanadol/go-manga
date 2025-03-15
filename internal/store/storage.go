package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("resource not found")
	QueryTimeoutDuration = time.Second * 5
)

type Storage struct {
	Mangas interface {
		Create(context.Context, *Manga) error
		GetByID(context.Context, int64) (*Manga, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Mangas: &MangaStore{db},
	}
}
