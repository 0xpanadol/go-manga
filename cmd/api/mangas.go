package main

import (
	"net/http"

	"github.com/0xpanadol/go-manga/internal/store"
)

type CreateMangaPayload struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Slug        string   `json:"slug"`
	Cover       string   `json:"cover"`
	Tags        []string `json:"tags"`
}

func (app *application) createMangaHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateMangaPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	// Update Later
	userID := 1

	manga := &store.Manga{
		Title:       payload.Title,
		Description: payload.Description,
		Slug:        payload.Slug,
		Cover:       payload.Cover,
		Tags:        payload.Tags,
		UserID:      int64(userID),
	}

	ctx := r.Context()
	if err := app.store.Mangas.Create(ctx, manga); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusOK, manga); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
