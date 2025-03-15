package main

import (
	"log"
	"net/http"
	"time"

	"github.com/0xpanadol/go-manga/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config config
	store  store.Storage
}

type config struct {
	addr string
	db   dbConfig
	env  string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  time.Duration
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)

		r.Route("/manga", func(r chi.Router) {
			r.Post("/", app.createMangaHandler)

			// Mangas
			r.Route("/{mangaID}", func(r chi.Router) {
				r.Use(app.mangaContextMiddleware)
				r.Get("/", app.getMangaHandler)

				// Chapters GET, POST, PATCH, DELETE
				r.Route("/chapters", func(r chi.Router) {
					r.Get("/", app.getChapterHandler)
					r.Post("/", app.createChapterHandler)

					r.Route("/{chapterID}", func(r chi.Router) {
						r.Delete("/", app.deleteChapterHandler)
					})
				})
			})
		})
	})
	return r
}

func (app *application) run(mux http.Handler) error {
	server := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Minute,
	}

	log.Printf("server is running at %s", app.config.addr)
	return server.ListenAndServe()
}
