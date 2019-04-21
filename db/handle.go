package db

// This file is responsible for bridging the logic between the REST controller and the database.
// A request will come inn from the routes.go file, The middlewares will inspect the request
// The user will be "authenticated" and if successful the request will reach the handler here
// This handler will then call the appropriate database functions and respond to the request.

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/sevren/test/core"
	log "github.com/sirupsen/logrus"
)

type LicenseResp struct {
	Licenses []string `json:"licenses"`
}

type ErrorPayload struct {
	Error ErrorResponse `json:"error"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

const USER_KEY = "user"
const CHALLENGE_3_KEY = "challenge3-features"

//HandleLicenses - Calls the database functions to retieve the license information from the database
// and then generates a proper license depending on wether or not the microservice is connected to rabbitmq
func (store *Dao) HandleLicenses(w http.ResponseWriter, r *http.Request) {

	// Get the authenticated user from the middleware
	user, ok := r.Context().Value(USER_KEY).(string)
	if !ok {
		errPayload := ErrorPayload{ErrorResponse{Code: http.StatusBadRequest, Message: "User not authenticated"}}
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, errPayload)
		return
	}

	// Get the licenses from the SQL lite database given by the authenticated user
	licenseRefs, err := store.GetLicenses(user)
	if err != nil {
		errPayload := ErrorPayload{ErrorResponse{Code: http.StatusOK, Message: err.Error()}}
		render.Status(r, http.StatusOK)
		render.JSON(w, r, errPayload)
		return
	}

	// Sets a flag if we should be using challenge 3 stuff
	// This will only be true if and only if the microservice is connected to rabbitmq on startup
	challenge3features, ok := r.Context().Value(CHALLENGE_3_KEY).(bool)
	if !ok {
		challenge3features = false
	}

	l := []string{}

	// If we are not connected to rabbitmq - we generate licenses using challenge 1 design - base64 encoding
	// if we are connected to rabbitmq - we generate licenses using challenge 3 - we use hashids and a salt of user:license-from-database
	if !challenge3features {
		log.Info("Now generating licenses with Challenge 1 features")
		l = core.GenerateLicenses(licenseRefs)
	} else {
		log.Info("Now generating licenses with Challenge 3 features")
		l = core.GenerateBetterLicenses(user, licenseRefs)
	}

	render.JSON(w, r, LicenseResp{l})
}
