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

		// attempts to decode the JSON body
		// This maps to a POST body {"password":"something"}
		decoder := json.NewDecoder(r.Body)
		var p PostReq
		err := decoder.Decode(&p)
		if err != nil {
			panic(err)
		}

		// extracts the username from the url
		user := chi.URLParam(r, "user")

		// looks up in the database the user and attempts to match the password
		// if this fails we get 403, forbidden and no response
		u := models.User_licenses{}
		if err = d.DB.Where(&models.User_licenses{Username: user, Password: p.Password}).First(&u).Error; err != nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		// Pass the "authenticated" user in the context down to the next handler
		// In our case the HandleLicenses function from the db package would be called next.
		ctx := context.WithValue(r.Context(), "user", u.Username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
