package main

import (
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/evilsocket/islazy/log"
	"github.com/evilsocket/shieldwall/firewall"
	"github.com/evilsocket/shieldwall/version"
	"os"
	"os/signal"
	"sort"
	"syscall"
	"time"
)

var (
	err     = (error)(nil)
	conf    = (* Config)(nil)
	state   = (*State)(nil)
	signals = make(chan os.Signal, 1)
)

func signalHandler() {
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	s := <-signals
	log.Raw("\n")
	log.Warning("RECEIVED SIGNAL: %s", s)
	os.Exit(1)
}

func addAllowRules(s *State) {
	for _, address := range conf.Allow {
		state.Rules = append(state.Rules, firewall.Rule{
			Type:     firewall.RuleAllow,
			Address:  address,
			Protocol: firewall.ProtoAll,
			Ports:    []string{firewall.AllPorts},
		})
	}
}

func hashObject(v interface{}) (string, error) {
	if raw, err := json.Marshal(v); err != nil {
		return "", err
	} else {
		return fmt.Sprintf("%x", sha256.Sum256(raw)), nil
	}
}

func rulesHash(rules []firewall.Rule) string {
	// make sure the order is always the same
	sort.Slice(rules, func(i, j int) bool {
		return rules[i].CreatedAt.Before(rules[j].CreatedAt)
	})
	hash, err := hashObject(rules)
	if err != nil {
		log.Warning("can't hash rules: %v", err)
	}
	return hash
}

func main() {
	flag.Parse()

	if showVersion {
		fmt.Printf("shiedwall agent v%s\n", version.Version)
		return
	}

	setupLogging()
	go signalHandler()

	log.Info("shieldwall agent v%s", version.Version)

	// load configuration
	if conf, err = LoadAgentConfig(confFile); err != nil {
		log.Fatal("error reading %s: %v", confFile, err)
	}

	// initialize firewall
	if err = firewall.SetPath(conf.IPTablesPath); err != nil {
		log.Fatal("%v", err)
	}

	// load saved state and run rules
	if state, err = LoadState(conf.DataPath); err != nil {
		log.Warning("%v", err)
	}

	// new state, add the entries allowed by configuration
	if len(state.Rules) == 0 && len(conf.Allow) > 0 {
		addAllowRules(state)
	}

	// apply previous rules from the saved state
	if err = firewall.Apply(state.Rules, conf.Drops); err != nil {
		log.Fatal("%v", err)
	}

	if conf.Update {
		go updater()
	}

	api := NewAPI(conf.API)
	// main loop
	for {
		prevHash := rulesHash(state.Rules)
		if rules, err := api.FetchRules(); err != nil {
			log.Error("error polling api: %v", err)
		} else {
			state.Rules = rules
			if len(conf.Allow) > 0 {
				addAllowRules(state)
			}
			newHash := rulesHash(state.Rules)
			if prevHash != newHash {
				log.Info("applying %d rules", len(state.Rules))
				if err = firewall.Apply(state.Rules, conf.Drops); err != nil {
					log.Fatal("%v", err)
				}
	 		} else {
	 			log.Debug("no changes")
			}
		}

		if err = state.Save(conf.DataPath); err != nil {
			log.Error("error saving state to %s: %v", conf.DataPath, err)
		} else {
			log.Debug("state saved to %s", conf.DataPath)
		}

		time.Sleep(time.Second * time.Duration(conf.API.Period))
	}
}
