package main

import (
	"encoding/json"
	"fmt"
	"github.com/evilsocket/islazy/fs"
	"github.com/evilsocket/islazy/log"
	"github.com/evilsocket/shieldwall/firewall"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

const stateFileName = "state.json"

type State struct {
	UpdatedAt time.Time       `json:"updated_at"`
	Rules     []firewall.Rule `json:"rules"`
}

func LoadState(path string) (*State, error) {
	if !fs.Exists(path) {
		log.Info("creating %s", path)
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return nil, err
		}
	}

	fileName := filepath.Join(path, stateFileName)
	if !fs.Exists(fileName) {
		log.Debug("%s doesn't exist, returning initial state", fileName)
		return &State{
			UpdatedAt: time.Now(),
		}, nil
	}

	log.Debug("loading state from %s ...", fileName)

	var state State

	raw, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("error reading %s: %v", fileName, err)
	}

	if err = json.Unmarshal(raw, &state); err != nil {
		return nil, fmt.Errorf("error decoding %s: %v", fileName, err)
	}

	return &state, nil
}

func (s *State) Save(path string) error {
	s.UpdatedAt = time.Now()
	raw, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(path, stateFileName), raw, os.ModePerm)
}