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
	Server  string
	Token   string
	Timeout int
}

func (a API) FetchRules() ([]firewall.Rule, error) {
	client := &http.Client{}
	if a.Timeout > 0 {
		client.Timeout = time.Duration(a.Timeout) * time.Second
	}

	if strings.Index(a.Server, "://") == -1 {
		a.Server = "https://" + a.Server
	}
	url := fmt.Sprintf("%s/api/v1/rules", a.Server)

	log.Debug("polling %s", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// agent authentication
	req.Header.Set("X-ShieldWall-Agent-Token", a.Token)
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
