package api

import "github.com/evilsocket/shieldwall/mailer"

type EmailConfig struct {
	From string        `yaml:"from"`
	SMTP mailer.Config `yaml:"smtp"`
}

type Config struct {
	URL        string `yaml:"url"`
	Address    string `yaml:"address"`
	ReqMaxSize int64  `yaml:"req_max_size"`
	TokenTTL   int64  `yaml:"token_ttl"`
	Secret     string `yaml:"secret"`
	MaxAgents  int64  `yaml:"max_agents_per_user"`
}
