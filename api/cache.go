package api

import (
	"encoding/json"
	"github.com/evilsocket/islazy/log"
	"github.com/evilsocket/shieldwall/database"
	"github.com/evilsocket/shieldwall/firewall"
	"gorm.io/datatypes"
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
