package api

import (
	"github.com/evilsocket/shieldwall/database"
	"github.com/go-chi/chi"
	"net/http"
)

func (api *API) UserVerify(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "verification")
	if code == "" {
		JSON(w, http.StatusBadRequest, nil)
		return
	}

	if err := database.VerifyUser(code); err != nil {
		JSON(w, http.StatusBadRequest, err)
		return
	}

	JSON(w, http.StatusOK, "user successfully verified")
}