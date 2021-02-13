package main

import (
	"github.com/evilsocket/shieldwall/firewall"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type APIConfig struct {
	// api server hostname
	Server string `yaml:"server"`
	// auth token
	Token string `yaml:"token"`
	// api polling period
	Period int `yaml:"period"`
	// api timeout
	Timeout int `yaml:"timeout"`
}

type Config struct {
	// TODO: refactor to support different firewalls
	IPTablesPath string `yaml:"iptables"`
	// where lists are stored
	DataPath string `yaml:"data"`
	// list of addresses to pre allow
	Allow []string `yaml:"allow"`
	// check for newer versions and self update the agent
	Update bool `yaml:"update"`
	// dropped packet logging
	Drops firewall.DropConfig `yaml:"drops"`
	// api config
	API APIConfig `yaml:"api"`
}

func LoadAgentConfig(fileName string) (*Config, error) {
	raw, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	conf := Config{}

	err = yaml.Unmarshal(raw, &conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}
