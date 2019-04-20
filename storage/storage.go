package storage

import (
	"net/http"
)

type ItemStore interface {
	StoreUsedLicenses(string)
	GetLicenses(string) ([]string, error)

	HandleLicenses(http.ResponseWriter, *http.Request)
	AuthUser(http.Handler) http.Handler
}
