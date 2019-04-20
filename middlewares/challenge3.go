package middlewares

import (
	"context"
	"net/http"
)

type PostReq struct {
	Password string `json:"password"`
}

// Middleware to set wether or not challenge 3 features are allowed

func Challenge3(feature bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "challenge3-features", feature)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
