package api

import (
	"encoding/json"
	"fmt"
	"github.com/evilsocket/islazy/log"
	"github.com/evilsocket/shieldwall/database"
	"github.com/go-chi/chi"
	"io"
	"net/http"
	"strconv"
)

func (api *API) authorized(w http.ResponseWriter, r *http.Request) *database.User {
	client := clientIP(r)
	tokenHeader := reqToken(r)
	if tokenHeader == "" {
		log.Debug("unauthenticated request from %s", client)
		ERROR(w, http.StatusUnauthorized, ErrUnauthorized)
		return nil
	}

	claims, err := api.validateToken(tokenHeader)
	if err != nil {
		log.Warning("token error for %s: %v", client, err)
		ERROR(w, http.StatusUnauthorized, ErrUnauthorized)
		return nil
	}

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

func (api *API) UserCreateAgent(w http.ResponseWriter, r *http.Request) {
	if user := api.authorized(w, r); user != nil {
		if api.config.MaxAgents > 0 && len(user.Agents) >= api.config.MaxAgents {
			ERROR(w, http.StatusForbidden, fmt.Errorf("max %d agents per user reached", api.config.MaxAgents))
			return
		}

		var req AgentCreationRequest

		defer r.Body.Close()

		client := clientIP(r)
		reader := io.LimitReader(r.Body, api.config.ReqMaxSize)
		decoder := json.NewDecoder(reader)

		err := decoder.Decode(&req)
		if err != nil {
			log.Warning("[%s] error decoding agent creation request: %v", client, err)
			JSON(w, http.StatusBadRequest, nil)
			return
		}

		agent, err := database.RegisterAgent(user, req.Name, req.Rules)
		if err != nil {
			ERROR(w, http.StatusBadRequest, err)
			return
		}

		log.Info("registered new agent %s for %s", agent.Name, user.Email)

		JSON(w, http.StatusOK, agent)
	} else {
		JSON(w, http.StatusForbidden, nil)
	}
}

func (api *API) UserGetAgents(w http.ResponseWriter, r *http.Request) {
	if user := api.authorized(w, r); user != nil {
		JSON(w, http.StatusOK, user.Agents)
	} else {
		JSON(w, http.StatusForbidden, nil)
	}
}

func (api *API) UserGetAgent(w http.ResponseWriter, r *http.Request) {
	if user := api.authorized(w, r); user != nil {
		id := chi.URLParam(r, "id")
		if id == "" {
			JSON(w, http.StatusBadRequest, nil)
			return
		}

		idNum, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			ERROR(w, http.StatusBadRequest, err)
			return
		}

		for _, agent := range user.Agents {
			if agent.ID == uint(idNum) {
				JSON(w, http.StatusOK, agent)
				return
			}
		}

		ERROR(w, http.StatusNotFound, fmt.Errorf("not found"))
	} else {
		JSON(w, http.StatusForbidden, nil)
	}
}

func (api *API) UserUpdateAgent(w http.ResponseWriter, r *http.Request) {
	if user := api.authorized(w, r); user != nil {
		id := chi.URLParam(r, "id")
		if id == "" {
			JSON(w, http.StatusBadRequest, nil)
			return
		}

		idNum, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			ERROR(w, http.StatusBadRequest, err)
			return
		}

		var req AgentUpdateRequest

		defer r.Body.Close()

		client := clientIP(r)
		reader := io.LimitReader(r.Body, api.config.ReqMaxSize)
		decoder := json.NewDecoder(reader)

		if err = decoder.Decode(&req); err != nil {
			log.Warning("[%s] error decoding agent creation request: %v", client, err)
			JSON(w, http.StatusBadRequest, nil)
			return
		}

		for _, agent := range user.Agents {
			if agent.ID == uint(idNum) {
				if err = database.UpdateAgent(&agent, req.Name, req.Rules); err != nil {
					ERROR(w, http.StatusBadRequest, err)
				} else {
					cacheByAgentToken.Delete(agent.Token)
					JSON(w, http.StatusOK, agent)
				}
				return
			}
		}

		ERROR(w, http.StatusNotFound, fmt.Errorf("not found"))
	} else {
		JSON(w, http.StatusForbidden, nil)
	}
}

func (api *API) UserDeleteAgent(w http.ResponseWriter, r *http.Request) {
	if user := api.authorized(w, r); user != nil {
		id := chi.URLParam(r, "id")
		if id == "" {
			JSON(w, http.StatusBadRequest, nil)
			return
		}

		idNum, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			ERROR(w, http.StatusBadRequest, err)
			return
		}

		for _, agent := range user.Agents {
			if agent.ID == uint(idNum) {
				if err = database.Delete(&agent); err != nil {
					ERROR(w, http.StatusInternalServerError, err)
				} else {
					cacheByAgentToken.Delete(agent.Token)
					JSON(w, http.StatusOK, "agent deleted")
				}
				return
			}
		}

		ERROR(w, http.StatusNotFound, fmt.Errorf("not found"))
	} else {
		JSON(w, http.StatusForbidden, nil)
	}
}
