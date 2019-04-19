package main

import (
	"github.com/sevren/test/storage"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"

	"github.com/go-chi/render"
)

func Routes(store storage.ItemStore) (*chi.Mux, error) {

	corsConf := corsConfig()
	r := chi.NewRouter()
	r.Use(corsConf.Handler)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(middleware.Logger)

	r.Route("/{user}", func(r chi.Router) {
		r.Post("/", store.User)
		r.Route("/licenses", func(r chi.Router) {
			r.Use(store.AuthUser)
			r.Post("/", store.Licenses)
		})
	})

	return r, nil
}

func corsConfig() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "X-CSRF-Token", "Cache-Control", "X-Requested-With"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
}
