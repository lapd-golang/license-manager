package db

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/sevren/test/helpers"
	"github.com/sevren/test/models"
)

type LicenseResp struct {
	Licenses []string `json:"licenses"`
}

type EmptyResponse struct{}

const CTXLICENSE_KEY = "licenses"

func (store *Dao) User(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	render.JSON(w, r, EmptyResponse{})
}

func (store *Dao) Licenses(w http.ResponseWriter, r *http.Request) {

	licenseRefs, ok := r.Context().Value(CTXLICENSE_KEY).(models.Licenses)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	l := helpers.GenerateLicenses(licenseRefs)
	render.JSON(w, r, LicenseResp{l})
}
