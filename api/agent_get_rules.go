package api

import (
	"encoding/json"
	"github.com/evilsocket/islazy/log"
	"github.com/evilsocket/shieldwall/database"
	"github.com/evilsocket/shieldwall/firewall"
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

func (api *API) expireRules(jsonbRules datatypes.JSON, doLog bool) (datatypes.JSON, int, error) {
	var rules []firewall.Rule

	expired := 0
	notExpired := make([]firewall.Rule, 0)

	if err := json.Unmarshal(jsonbRules, &rules); err != nil {
		return nil, 0, err
	}

	for _, rule := range rules {
		if rule.Expired() {
			if doLog {
				log.Info("rule expired %#v", rule)
			}
			expired++
		} else {
			notExpired = append(notExpired, rule)
		}
	}

	return database.ToJSONB(notExpired), expired, nil
}

func (api *API) GetRules(w http.ResponseWriter, r *http.Request) {
	agentIP := clientIP(r)
	agentToken := r.Header.Get("X-ShieldWall-Agent-Token")
	agentUA := r.Header.Get("User-Agent")

	if agentToken == "" {
		log.Warning("[%s %s] received rules request with no token", agentIP, agentUA)
		JSON(w, http.StatusBadRequest, nil)
		return
	}

	cacheWhat := "miss"

	// check cache first
	entry, found := cacheByAgentToken.Load(agentToken)
	if found {
		// expired?
		cached := entry.(*cachedRules)
		if int64(time.Since(cached.CachedAt).Seconds()) >= api.config.CacheTTL {
			log.Debug("agent cache expired")
			cacheByAgentToken.Delete(agentToken)
		} else {
			// check expired rules
			_, expired, err := api.expireRules(cached.Rules, false)
			if err != nil {
				log.Error("error checking rules expiration: %v", err)
				JSON(w, http.StatusInternalServerError, nil)
				return
			}
			// bypass and invalidate cache if there are expired rules
			// in order to cache a fresh copy of the model
			if expired == 0 {
				w.Header().Set("shieldwall-cache", "hit")
				JSON(w, http.StatusOK, cached.Rules)
				return
			} else {
				cacheByAgentToken.Delete(agentToken)
				cacheWhat = "purge" // let the client know what happened ^_^
			}
		}
	}

	w.Header().Set("shieldwall-cache", cacheWhat)

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

	// check expired rules
	agent.Rules, _, err = api.expireRules(agent.Rules, true)
	if err != nil {
		log.Error("error checking rules expiration: %v", err)
		JSON(w, http.StatusInternalServerError, nil)
		return
	}

	// log.Debug("[%s %s] successfully authenticated", agentIP, agentUA)

	agent.UpdatedAt = time.Now()
	agent.Address = agentIP
	agent.UserAgent = agentUA

	if err = agent.Save(); err != nil {
		log.Error("error updating agent: %v", err)
	}

	// save to cache
	cacheByAgentToken.Store(agentToken, &cachedRules{
		CachedAt: time.Now(),
		Rules:    agent.Rules,
	})

	JSON(w, http.StatusOK, agent.Rules)
}
