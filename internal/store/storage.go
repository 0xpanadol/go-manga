package store

import (
	"context"
	"database/sql"
)

type Storage struct {
	Mangas interface {
		Create(context.Context, *Manga) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Mangas: &MangaStore{db},
	}
}
