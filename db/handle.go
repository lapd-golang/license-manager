package db

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/sevren/test/core"
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

	l := core.GenerateLicenses(licenseRefs)
	render.JSON(w, r, LicenseResp{l})
}
