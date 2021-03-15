package api

import (
	"fmt"
	"github.com/evilsocket/islazy/log"
	"github.com/evilsocket/islazy/str"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	cfIPv4URL = "https://www.cloudflare.com/ips-v4"
	cfIPv6URL = "https://www.cloudflare.com/ips-v6"
)

var (
	cfURLs    = []string{cfIPv4URL, cfIPv6URL}
	cfTTL     = time.Minute * time.Duration(60)
	cfTime    = time.Time{}
	cfLock    = sync.Mutex{}
	cfSubnets = []string(nil)
)

func cfGetSubnets() ([]string, error) {
	cfLock.Lock()
	defer cfLock.Unlock()

	if cfTime.IsZero() || time.Since(cfTime) >= cfTTL {
		cfSubnets = nil
		for _, url := range cfURLs {
			log.Info("updating cloudflare subnets from %s ...", url)
			if resp, err := http.Get(url); err != nil {
				return nil, err
			} else if resp.StatusCode != http.StatusOK {
				return nil, fmt.Errorf("%s [%d] %s", url, resp.StatusCode, resp.Status)
			} else {
				defer resp.Body.Close()
				if raw, err := ioutil.ReadAll(resp.Body); err != nil {
					return nil, fmt.Errorf("%s %v", url, err)
				} else if parts := strings.Split(string(raw), "\n"); len(parts) == 0 {
					return nil, fmt.Errorf("%s unexpected response: %s", url, string(raw))
				} else {
					cfTime = time.Now()
					for _, part := range parts {
						if part = str.Trim(part); part != "" {
							log.Info("cf: %s", part)
							cfSubnets = append(cfSubnets, part)
						}
					}
				}
			}
		}
	}

	return cfSubnets, nil
}

func (api *API) GetCloudflareSubnets(w http.ResponseWriter, r *http.Request) {
	if user := api.authorized(w, r); user != nil {
		if nets, err := cfGetSubnets(); err != nil {
			log.Error("error getting cloudflare subnets: %v", err)
			ERROR(w, http.StatusInternalServerError, err)
		} else {
			JSON(w, http.StatusOK, nets)
		}
	} else {
		JSON(w, http.StatusForbidden, nil)
	}
}
