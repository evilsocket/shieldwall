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
	Use2FA      bool   `json:"use_2fa"`
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

		if _, err := database.UpdateUser(user, client, req.NewPassword, req.Use2FA); err != nil {
			log.Debug("[%s] %v", client, err)
			ERROR(w, http.StatusBadRequest, err)
			return
		}

		token, err := api.tokenFor(user, !user.Use2FA)
		if err != nil {
			log.Error("error updating token for user %d: %v", user.ID, err)
		}

		JSON(w, http.StatusOK, UserResponse{
			Token:   token,
			User:    user,
			Address: client,
		})
	} else {
		JSON(w, http.StatusForbidden, nil)
	}
}
