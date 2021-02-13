package main

import (
	"encoding/json"
	"fmt"
	"github.com/evilsocket/islazy/log"
	"github.com/evilsocket/shieldwall/firewall"
	"github.com/evilsocket/shieldwall/version"
	"net/http"
	"runtime"
	"strings"
	"time"
)

// API client
type API struct {
	config APIConfig
}

func NewAPI(config APIConfig) *API {
	return &API{
		config: config,
	}
}

func (a API) FetchRules() ([]firewall.Rule, error) {
	client := &http.Client{}
	if a.config.Timeout > 0 {
		client.Timeout = time.Duration(a.config.Timeout) * time.Second
	}

	if strings.Index(a.config.Server, "://") == -1 {
		a.config.Server = "https://" + a.config.Server
	}
	url := fmt.Sprintf("%s/api/v1/rules", a.config.Server)

	log.Debug("polling %s", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// agent authentication
	req.Header.Set("X-ShieldWall-Agent-Token", a.config.Token)
	req.Header.Set("User-Agent", fmt.Sprintf(
		"ShieldWall Agent v%s (%s %s)",
		version.Version,
		runtime.GOOS,
		runtime.GOARCH))

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%d (%s)", res.StatusCode, res.Status)
	}

	var rules []firewall.Rule
	if err = json.NewDecoder(res.Body).Decode(&rules); err != nil {
		return nil, err
	}

	return rules, nil
}
