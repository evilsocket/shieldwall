package api

import (
	"encoding/json"
	"fmt"
	"github.com/evilsocket/islazy/log"
	"github.com/evilsocket/shieldwall/database"
	"io"
	"net/http"
)

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	Token   string         `json:"token"`
	User    *database.User `json:"data"`
	Address string         `json:"address"`
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
		ERROR(w, http.StatusUnauthorized, fmt.Errorf("invalid credentials")) // TODO: Change Errorf to errors
		return
	} else if token, err := api.tokenFor(user, !user.Use2FA); err != nil {
		log.Error("error creating token for user %d: %v", user.ID, err)
		ERROR(w, http.StatusInternalServerError, err)
		return
	} else {
		log.Debug("[%s] user %s logged in", client, user.Email)
		if user.Use2FA {
			log.Info("[%s] sending verification email to %s", client, user.Email)

			emailSubject := "shieldwall.me login verification"
			emailBody := fmt.Sprintf("Your verification code is %s<br/><br/>", user.Verification) +
				fmt.Sprintf("The address %s tried to login to your shieldwall account, "+
					"if that wasn't you, you should change your credentials immediately.", client)

			if err = api.sendmail.Send(api.mail.From, user.Email, emailSubject, emailBody); err != nil {
				log.Error("error sending verification email to %s: %v", user.Email, err)
				ERROR(w, http.StatusInternalServerError, fmt.Errorf("error sending verification email"))
				return
			} else {
				log.Debug("verification email sent to %s (%s)", user.Email, user.Verification)
			}
		}

		JSON(w, http.StatusOK, UserResponse{
			Token:   token,
			User:    user,
			Address: client,
		})
	}
}

func (api *API) authorizedForStep2(w http.ResponseWriter, r *http.Request) *database.User {
	client := clientIP(r)
	tokenHeader := reqToken(r)
	if tokenHeader == "" {
		log.Debug("unauthenticated request from %s", client)
		ERROR(w, http.StatusUnauthorized, ErrUnauthorized)
		return nil
	}

	claims, err := api.validateToken(tokenHeader)
	if err != nil && err != ErrToken2FA {
		log.Warning("token error for %s: %v", client, err)
		ERROR(w, http.StatusUnauthorized, ErrUnauthorized)
		return nil
	}

	log.Debug("user claims for step2: %#v", claims)
	user, err := database.FindUserByID(int(claims["user_id"].(float64)))
	if err != nil {
		ERROR(w, http.StatusUnauthorized, err)
		return nil
	} else if user == nil {
		log.Warning("client %s tried to authenticated with unknown claims '%v'", client, claims)
		ERROR(w, http.StatusUnauthorized, ErrUnauthorized)
		return nil
	}

	user.Address = client
	return user
}

type Step2Request struct {
	Code string `json:"code"`
}

func (api *API) UserSecondStep(w http.ResponseWriter, r *http.Request) {
	if user := api.authorizedForStep2(w, r); user != nil {
		var req Step2Request

		defer r.Body.Close()

		client := clientIP(r)
		reader := io.LimitReader(r.Body, api.config.ReqMaxSize)
		decoder := json.NewDecoder(reader)

		err := decoder.Decode(&req)
		if err != nil {
			log.Warning("[%s] error decoding user step2 request: %v", client, err)
			JSON(w, http.StatusBadRequest, nil)
			return
		}

		if req.Code != user.Verification {
			log.Warning("[%s] wrong code %s", client, req.Code)
			JSON(w, http.StatusForbidden, nil)
			return
		}

		if token, err := api.tokenFor(user, true); err != nil {
			log.Error("error creating token for user %d: %v", user.ID, err)
			ERROR(w, http.StatusInternalServerError, err)
			return
		} else {
			JSON(w, http.StatusOK, UserResponse{
				Token:   token,
				User:    user,
				Address: client,
			})
		}
	} else {
		JSON(w, http.StatusForbidden, nil)
	}
}
