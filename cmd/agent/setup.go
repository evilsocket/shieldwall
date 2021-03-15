package main

import (
	"flag"
	"github.com/evilsocket/islazy/log"
	"github.com/evilsocket/shieldwall/firewall"
	"time"
)

var (
	confFile          = ""
	debug             = false
	showVersion       = false
	logfile           = ""
	updateCheckPeriod = time.Duration(10) * time.Minute
)

func init() {
	flag.StringVar(&confFile, "config", "/etc/shieldwall/config.yaml", "YAML configuration file.")
	flag.BoolVar(&debug, "debug", debug, "Enable debug logs.")
	flag.StringVar(&logfile, "log", logfile, "Log messages to this file instead of standard error.")
	flag.DurationVar(&updateCheckPeriod, "update-check-period", updateCheckPeriod, "Self update polling period.")

	flag.BoolVar(&showVersion, "version", showVersion, "Show version and exit.")
	flag.BoolVar(&firewall.DryRun, "dry-run", firewall.DryRun, "Do not execute firewall commands.")
}

func setupLogging() {
	if logfile != "" {
		log.Output = logfile
	}

	if debug == true {
		log.Level = log.DEBUG
	} else {
		log.Level = log.INFO
	}

	log.DateFormat = "06-Jan-02"
	log.TimeFormat = "15:04:05"
	log.DateTimeFormat = "2006-01-02 15:04:05"
	log.Format = "{datetime} {level:color}{level:name}{reset} {message}"

	if err := log.Open(); err != nil {
		panic(err)
	}
}
