package storage

import (
	"net/http"
)

type ItemStore interface {
	Get()
	Save()
	User(http.ResponseWriter, *http.Request)
	Licenses(http.ResponseWriter, *http.Request)
	AuthUser(http.Handler) http.Handler
}
