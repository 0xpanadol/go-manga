package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/0xpanadol/go-manga/internal/store"
	"github.com/go-chi/chi/v5"
)

func (app *application) getChapterHandler(w http.ResponseWriter, r *http.Request) {
	manga := getMangaFromCtx(r.Context())

	ctx := r.Context()
	chapters, err := app.store.Chapters.GetByMangaID(ctx, manga.ID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFound(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, chapters); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

type CreateChapterPayload struct {
	Title   string   `json:"title" validate:"required,max=255"`
	Content []string `json:"content" validate:"required"`
	Number  float64  `json:"chapter_number" validate:"required"`
}

func (app *application) createChapterHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateChapterPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	mangaID := getMangaFromCtx(r.Context()).ID

	// Update Later
	userID := 1

	chapter := &store.Chapter{
		Title:   payload.Title,
		Content: payload.Content,
		Number:  payload.Number,
		MangaID: mangaID,
		UserID:  int64(userID),
	}

	err := app.store.Chapters.Create(r.Context(), chapter)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	data := map[string]any{
		"status":  "success",
		"message": "chapter has been published successfuly",
		"chapter": chapter,
	}

	if err := app.jsonResponse(w, http.StatusOK, data); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) deleteChapterHandler(w http.ResponseWriter, r *http.Request) {
	mangaID := chi.URLParam(r, "mangaID")
	mangaIDInt, err := strconv.ParseInt(mangaID, 10, 64)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	chapterID := chi.URLParam(r, "chapterID")
	chapterIDInt, err := strconv.ParseInt(chapterID, 10, 64)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.store.Chapters.Delete(r.Context(), chapterIDInt, mangaIDInt); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFound(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, map[string]string{"status": "success", "message": "chapter has been deleted successfuly."}); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
