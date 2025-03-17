package main

import (
	"net/http"

	"github.com/0xpanadol/go-manga/internal/store"
)

type CreateCommentPayload struct {
	Content string `json:"content" validate:"required,min=6,max=1000"`
}

func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateCommentPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	// TODO: Update later
	userID := 1
	mangaID := 7

	comment := &store.Comment{
		Content:   payload.Content,
		UserID:    int64(userID),
		MangaID:   int64(mangaID),
		ChapterID: 0,
	}

	if err := app.store.Comments.Create(r.Context(), comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
