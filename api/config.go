package api

import "github.com/evilsocket/shieldwall/mailer"

type EmailConfig struct {
	From string        `yaml:"from"`
	SMTP mailer.Config `yaml:"smtp"`
}

type Config struct {
	URL           string   `yaml:"url"`
	SSL           bool     `yaml:"ssl"`
	CertsCache    string   `yaml:"certs_cache"`
	Domains       []string `yaml:"domains"`
	Address       string   `yaml:"address"`
	ReqMaxSize    int64    `yaml:"req_max_size"`
	TokenTTL      int      `yaml:"token_ttl"`
	Secret        string   `yaml:"secret"`
	MaxAgents     int      `yaml:"max_agents_per_user"`
	CacheTTL      int      `yaml:"cache_ttl"`
	AllowNewUsers bool     `yaml:"allow_new_users"`
}
