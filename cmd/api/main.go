package main

import (
	"flag"
	"github.com/evilsocket/islazy/log"
	"github.com/evilsocket/shieldwall/api"
	"github.com/evilsocket/shieldwall/database"
	"github.com/evilsocket/shieldwall/mailer"
	"github.com/evilsocket/shieldwall/version"
)

var (
	err  = (error)(nil)
	conf = (* Config)(nil)
	mail = (*mailer.Mailer)(nil)
)

func main() {
	flag.Parse()

	setupLogging()

	log.Info("shieldwall api v%s", version.Version)

	if conf, err = LoadConfig(confFile); err != nil {
		log.Fatal("error reading %s: %v", confFile, err)
	}

	if mail, err = mailer.New(conf.Email.SMTP); err != nil {
		log.Fatal("error creating mailer: %v", err)
	}

	if err = database.Setup(conf.Database); err != nil {
		log.Fatal("error connecting to database: %v", err)
	}

	server := api.Setup(conf.Api, conf.Email, mail)

	server.Run()
}
