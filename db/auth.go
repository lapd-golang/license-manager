package db

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/sevren/test/models"
)

type PostReq struct {
	Password string `json:"password"`
}

// Middleware to handle user Authentication from the database. This will run before every REST request

// Using the existing database connection we check to see if the posted user successfully authenticated
// if sucessful the request is passed down the chain to the function handling the call
// if unsucessful the request is shutdown immediatly
func (d *Dao) AuthUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var p PostReq
		err := decoder.Decode(&p)
		if err != nil {
			panic(err)
		}

		user := chi.URLParam(r, "user")

		u := models.User_licenses{}
		if err = d.DB.Where(&models.User_licenses{Username: user, Password: p.Password}).First(&u).Error; err != nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), "user", u.Username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
