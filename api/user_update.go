package api

import (
	"encoding/json"
	"github.com/evilsocket/islazy/log"
	"github.com/evilsocket/shieldwall/database"
	"io"
	"net/http"
)

type UserUpdateRequest struct {
	NewPassword string `json:"password"`
}

func (api *API) UserUpdate(w http.ResponseWriter, r *http.Request) {
	if user := api.authorized(w, r); user != nil {
		defer r.Body.Close()

		var req UserUpdateRequest

		client := clientIP(r)
		reader := io.LimitReader(r.Body, api.config.ReqMaxSize)
		decoder := json.NewDecoder(reader)

		if err := decoder.Decode(&req); err != nil {
			log.Warning("[%s] error decoding user update request: %v", client, err)
			JSON(w, http.StatusBadRequest, nil)
			return
		}

		if _, err := database.UpdateUser(user, client, req.NewPassword); err != nil {
			log.Debug("[%s] %v", client, err)
			ERROR(w, http.StatusBadRequest, err)
			return
		}

		JSON(w, http.StatusOK, "OK")
	}
}
