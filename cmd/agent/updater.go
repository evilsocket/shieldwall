package main

import (
	"github.com/evilsocket/islazy/log"
	"github.com/evilsocket/shieldwall/version"
	"net/http"
	"regexp"
	"time"
)

var repo = "evilsocket/shieldwall"

var versionParser = regexp.MustCompile("^https://github\\.com/" + repo + "/releases/tag/v([\\d\\.a-z]+)$")

func updater() {
	log.Info("update checker started with a %s period", updateCheckPeriod)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	url := "https://github.com/" + repo + "/releases/latest"
	for {
		log.Debug("checking for updates %s", url)

		req, _ := http.NewRequest("GET", url, nil)
		resp, err := client.Do(req)
		if err != nil {
			log.Error("error while checking latest version: %v", err)
			continue
		}
		defer resp.Body.Close()

		location := resp.Header.Get("Location")

		log.Debug("location header = '%s'", location)

		m := versionParser.FindStringSubmatch(location)
		if len(m) == 2 {
			latest := m[1]
			log.Debug("Latest version is '%s'", latest)
			if version.Version != latest {
				log.Important("update to %s available at %s", latest, location)
			} else {
				log.Debug("no updates available")
			}
		} else {
			log.Debug("unexpected location header: '%s'", location)
		}

		time.Sleep(updateCheckPeriod)
	}
}