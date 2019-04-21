package storage

import (
	"net/http"
)

// Create an interface definition so we can attach methods to custom types
// We implement this interface on the Database Access object
type ItemStore interface {
	StoreUsedLicenses(string) error
	GetLicenses(string) ([]string, error)

	HandleLicenses(http.ResponseWriter, *http.Request)
	AuthUser(http.Handler) http.Handler
}
