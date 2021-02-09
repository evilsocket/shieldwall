package api

import (
	"github.com/evilsocket/islazy/log"
	"github.com/evilsocket/shieldwall/database"
	"net/http"
	"time"
)

func (api *API) GetRules(w http.ResponseWriter, r *http.Request) {
	agentIP := clientIP(r)
	agentToken := r.Header.Get("X-ShieldWall-Agent-Token")
	agentUA := r.Header.Get("User-Agent")

	if agentToken == "" {
		log.Warning("[%s %s] received rules request with no token", agentIP, agentUA)
		JSON(w, http.StatusBadRequest, nil)
		return
	}

	// TODO: add cache with ttl
	agent, err := database.FindAgentByToken(agentToken)
	if err != nil {
		log.Warning("[%s %s] error searching for token '%s': %v", agentIP, agentUA, agentToken, err)
		JSON(w, http.StatusBadRequest, nil)
		return
	} else if agent == nil {
		log.Warning("[%s %s] invalid token '%s'", agentIP, agentUA, agentToken)
		JSON(w, http.StatusUnauthorized, nil)
		return
	}

	log.Debug("[%s %s] successfully authenticated, returning %d rules", agentIP, agentUA, len(agent.Rules))

	agent.UpdatedAt = time.Now()
	agent.Address = agentIP
	agent.UserAgent = agentUA

	if err = agent.Save(); err != nil {
		log.Error("error updating agent: %v", err)
	}

	JSON(w, http.StatusOK, agent.Rules)
}
