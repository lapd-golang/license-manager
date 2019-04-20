package db

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/sevren/test/models"
	log "github.com/sirupsen/logrus"
)

type PostReq struct {
	Password string `json:"password"`
}

// Should lookup the user from the databse here
func (d *Dao) AuthUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var p PostReq
		err := decoder.Decode(&p)
		if err != nil {
			panic(err)
		}
		log.Println(p.Password)
		user := chi.URLParam(r, "user")

		u := models.User_licenses{}
		if err = d.DB.Where(&models.User_licenses{Username: user, Password: p.Password}).First(&u).Error; err != nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		fmt.Printf("FROM DB: %v", u)
		ctx := context.WithValue(r.Context(), "user", u.Username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
