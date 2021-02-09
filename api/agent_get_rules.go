package api

import (
	"github.com/evilsocket/islazy/log"
	"github.com/evilsocket/shieldwall/database"
	"gorm.io/datatypes"
	"net/http"
	"sync"
	"time"
)

type cachedRules struct {
	CachedAt time.Time
	Rules    datatypes.JSON
}

var cacheByAgentToken = sync.Map{}

func (api *API) GetRules(w http.ResponseWriter, r *http.Request) {
	agentIP := clientIP(r)
	agentToken := r.Header.Get("X-ShieldWall-Agent-Token")
	agentUA := r.Header.Get("User-Agent")

	if agentToken == "" {
		log.Warning("[%s %s] received rules request with no token", agentIP, agentUA)
		JSON(w, http.StatusBadRequest, nil)
		return
	}

	// check cache first
	entry, found := cacheByAgentToken.Load(agentToken)
	if found {
		// expired?
		cached := entry.(*cachedRules)
		if int64(time.Since(cached.CachedAt).Seconds()) >= api.config.CacheTTL {
			log.Debug("agent cache expired")
			cacheByAgentToken.Delete(agentToken)
		} else {
			w.Header().Set("shieldwall-cache", "hit")
			JSON(w, http.StatusOK, cached.Rules)
			return
		}
	}

	w.Header().Set("shieldwall-cache", "miss")

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

	// save to cache
	cacheByAgentToken.Store(agentToken, &cachedRules{
		CachedAt: time.Now(),
		Rules:    agent.Rules,
	})

	// log.Debug("[%s %s] successfully authenticated", agentIP, agentUA)

	agent.UpdatedAt = time.Now()
	agent.Address = agentIP
	agent.UserAgent = agentUA

	if err = agent.Save(); err != nil {
		log.Error("error updating agent: %v", err)
	}

	JSON(w, http.StatusOK, agent.Rules)
}
