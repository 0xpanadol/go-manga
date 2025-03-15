package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/0xpanadol/go-manga/internal/store"
	"github.com/go-chi/chi/v5"
)

type mangakey string

const mangaCtx mangakey = "manga"

type CreateMangaPayload struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Slug        string   `json:"slug"`
	Cover       string   `json:"cover"`
	Tags        []string `json:"tags"`
}

func (app *application) getMangaHandler(w http.ResponseWriter, r *http.Request) {
	manga := getMangaFromCtx(r.Context())

	if err := app.jsonResponse(w, http.StatusOK, manga); err != nil {
		app.internalServerError(w, r, err)
		return
	}
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

	if err := app.jsonResponse(w, http.StatusOK, manga); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) mangaContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mangaID := chi.URLParam(r, "mangaID")
		id, err := strconv.ParseInt(mangaID, 10, 64)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		ctx := r.Context()
		manga, err := app.store.Mangas.GetByID(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFound(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, mangaCtx, manga)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getMangaFromCtx(ctx context.Context) *store.Manga {
	manga, _ := ctx.Value(mangaCtx).(*store.Manga)
	return manga
}
