package api

import (
	"encoding/json"
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

	// parse resource information from POST body if present
	var resources *database.AgentResources
	if r.Method == http.MethodPost && r.Body != nil {
		resources = &database.AgentResources{}
		if err := json.NewDecoder(r.Body).Decode(resources); err != nil {
			log.Debug("could not decode resources from agent: %v", err)
			resources = nil
		}
	}

	cacheWhat := "miss"

	// check cache first
	entry, found := cacheByAgentToken.Load(agentToken)
	if found {
		// expired?
		cached := entry.(*cachedRules)
		if int(time.Since(cached.CachedAt).Seconds()) >= api.config.CacheTTL {
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
				// still update resources even when cache hit
				if resources != nil {
					go api.updateAgentResources(agentToken, agentIP, agentUA, resources)
				}
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

	agent.SeenAt = time.Now()
	agent.Address = agentIP
	agent.UserAgent = agentUA

	// update resource information if provided
	if resources != nil {
		if interfacesJSON, err := json.Marshal(resources.Interfaces); err == nil {
			agent.Interfaces = interfacesJSON
		}
		agent.ActiveInterface = resources.ActiveInterface
		if resourcesJSON, err := json.Marshal(resources.Resources); err == nil {
			agent.Resources = resourcesJSON
		}
	}

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

// updateAgentResources updates resource info asynchronously during cache hits
func (api *API) updateAgentResources(token, ip, ua string, resources *database.AgentResources) {
	agent, err := database.FindAgentByToken(token)
	if err != nil || agent == nil {
		return
	}

	agent.SeenAt = time.Now()
	agent.Address = ip
	agent.UserAgent = ua

	if resources != nil {
		if interfacesJSON, err := json.Marshal(resources.Interfaces); err == nil {
			agent.Interfaces = interfacesJSON
		}
		agent.ActiveInterface = resources.ActiveInterface
		if resourcesJSON, err := json.Marshal(resources.Resources); err == nil {
			agent.Resources = resourcesJSON
		}
	}

	if err = agent.Save(); err != nil {
		log.Error("error updating agent resources: %v", err)
	}
}
