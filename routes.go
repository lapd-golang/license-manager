package main

import (
	"net/http"

	"github.com/sevren/test/middlewares"
	"github.com/sevren/test/storage"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"

	"github.com/go-chi/render"
)

type EmptyResponse struct{}

func Routes(store storage.ItemStore, challenge3features bool) (*chi.Mux, error) {

	corsConf := corsConfig()
	r := chi.NewRouter()
	r.Use(corsConf.Handler)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(middleware.Logger)

	// If rabbitmq is connected then we can use challenge 3 stuff
	// this sets a middleware which adds the challenge 3 context to the request.
	if challenge3features {
		r.Use(middlewares.Challenge3(challenge3features))
	}

	r.Route("/{user}", func(r chi.Router) {
		r.Post("/", handleUser)
		r.Route("/licenses", func(r chi.Router) {
			r.Use(store.AuthUser)
			r.Post("/", store.HandleLicenses)
		})
	})

	return r, nil
}

// handles the use case for a simple /{user} rest call
func handleUser(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusBadRequest)
	render.JSON(w, r, EmptyResponse{})
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
