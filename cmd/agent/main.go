package main

import (
	"flag"
	"github.com/evilsocket/islazy/log"
	"github.com/evilsocket/shieldwall/firewall"
	"github.com/evilsocket/shieldwall/version"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	err               = (error)(nil)
	conf              = (* Config)(nil)
	state             = (*State)(nil)
	signals           = make(chan os.Signal, 1)
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

func main() {
	flag.Parse()

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
		log.Fatal("%v", err)
	}

	// new state, add the entries allowed by configuration
	if len(state.Rules) == 0 && len(conf.Allow) > 0 {
		addAllowRules(state)
	}

	// apply previous rules from the saved state
	if err = firewall.Apply(state.Rules); err != nil {
		log.Fatal("%v", err)
	}

	if conf.Update {
		go updater()
	}

	// TODO: split in APIConfig
	api := API{
		Server:  conf.Server,
		Token:   conf.Token,
		Timeout: conf.Timeout,
	}

	// main loop
	for {
		if rules, err := api.FetchRules(); err != nil {
			log.Error("error polling api: %v", err)
		} else {
			state.Rules = rules

			if num := len(state.Rules); num > 0 {
				log.Debug("got rules baby %#v", state.Rules)
				log.Info("applying %d rules", num)
			}

			if len(conf.Allow) > 0 {
				addAllowRules(state)
			}

			if err = firewall.Apply(state.Rules); err != nil {
				log.Fatal("%v", err)
			}
		}

		if err = state.Save(conf.DataPath); err != nil {
			log.Error("error saving state to %s: %v", conf.DataPath, err)
		} else {
			log.Debug("state saved to %s", conf.DataPath)
		}

		time.Sleep(time.Second * time.Duration(conf.Period))
	}
}
