package db

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

func (store *Dao) HandleLicenses(w http.ResponseWriter, r *http.Request) {

	user, ok := r.Context().Value(USER_KEY).(string)
	if !ok {
		errPayload := ErrorPayload{ErrorResponse{Code: http.StatusBadRequest, Message: "User not authenticated"}}
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, errPayload)
		return
	}

	licenseRefs, err := store.GetLicenses(user)
	if err != nil {
		errPayload := ErrorPayload{ErrorResponse{Code: http.StatusOK, Message: err.Error()}}
		render.Status(r, http.StatusOK)
		render.JSON(w, r, errPayload)
		return
	}

	challenge3features, ok := r.Context().Value(CHALLENGE_3_KEY).(bool)
	if !ok {
		challenge3features = false
	}

	l := []string{}

	if !challenge3features {
		log.Info("Now generating licenses with Challenge 1 features")
		l = core.GenerateLicenses(licenseRefs)
	} else {
		log.Info("Now generating licenses with Challenge 3 features")
		l = core.GenerateBetterLicenses(user, licenseRefs)
	}

	render.JSON(w, r, LicenseResp{l})
}
