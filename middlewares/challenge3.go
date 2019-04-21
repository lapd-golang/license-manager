package middlewares

import (
	"context"
	"net/http"
)

// Middleware to set wether or not challenge 3 features are allowed
// This is only ever true if the microservice connects to rabbitmq on startup
func Challenge3(feature bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "challenge3-features", feature)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
