package api

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/evilsocket/islazy/log"
	"github.com/evilsocket/shieldwall/database"
	"io"
	"net/http"
	"time"
)

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (api *API) tokenFor(user *database.User) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = user.ID
	claims["expires_at"] = time.Now().Add(time.Duration(api.config.TokenTTL) * time.Second).Format(time.RFC3339)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(api.config.Secret))
	if err != nil {
		return "", err
	}
	return signed, nil
}

func (api *API) UserLogin(w http.ResponseWriter, r *http.Request) {
	var req UserLoginRequest

	defer r.Body.Close()

	client := clientIP(r)
	reader := io.LimitReader(r.Body, api.config.ReqMaxSize)
	decoder := json.NewDecoder(reader)

	err := decoder.Decode(&req)
	if err != nil {
		log.Warning("[%s] error decoding user login request: %v", client, err)
		JSON(w, http.StatusBadRequest, nil)
		return
	}

	user, err := database.LoginUser(client, req.Email, req.Password)
	if err != nil {
		ERROR(w, http.StatusUnauthorized, err)
		return
	} else if user == nil {
		JSON(w, http.StatusUnauthorized, "invalid credentials")
		return
	} else if token, err := api.tokenFor(user); err != nil {
		log.Error("error creating token for user %d: %v", user.ID, err)
		ERROR(w, http.StatusInternalServerError, err)
		return
	} else {
		log.Debug("[%s] user %s logged in", client, user.Email)
		JSON(w, http.StatusOK, struct {
			Token string `json:"token"`
		}{
			Token: token,
		})
	}
}
