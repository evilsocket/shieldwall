package main

import (
	"flag"
	"github.com/evilsocket/islazy/log"
)

var (
	confFile = ""
	debug    = false
	logfile  = ""
)

func init() {
	flag.StringVar(&confFile, "config", "/etc/shieldwall/config.yaml", "YAML configuration file.")
	flag.BoolVar(&debug, "debug", debug, "Enable debug logs.")
	flag.StringVar(&logfile, "log", logfile, "Log messages to this file instead of standard error.")

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
