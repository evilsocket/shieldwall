package firewall

import (
	"fmt"
	"github.com/evilsocket/islazy/log"
	"github.com/evilsocket/islazy/str"
	"os/exec"
	"strconv"
	"sync"
)

// TODO: implement factory pattern with generic interface in order to support more firewalls

type DropConfig struct {
	Log    bool   `yaml:"log"`
	Limit  string `yaml:"limit"`
	Prefix string `yaml:"prefix"`
	Level  int    `yaml:"level"`
}

var lock = sync.Mutex{}

func cmd(bin string, args ...string) (string, error) {
	log.Debug("# %s %s", bin, args)
	raw, err := exec.Command(bin, args...).CombinedOutput()
	if err != nil {
		log.Error("%s", str.Trim(string(raw)))
		return "", err
	} else {
		return str.Trim(string(raw)), nil
	}
}

func reset() error {
	out, err := cmd(binary, "-F", "INPUT")
	if err != nil {
		return err
	} else {
		log.Debug("flush INPUT: %s", out)
	}

	out, err = cmd(binary, "-F", "LOGNDROP")
	if err != nil {
		return err
	} else {
		log.Debug("flush LOGNDROP: %s", out)
	}

	out, err = cmd(binary, "-X", "LOGNDROP")
	if err != nil {
		return err
	} else {
		log.Debug("delete LOGNDROP: %s", out)
	}

	return nil
}

func Reset() error {
	lock.Lock()
	defer lock.Unlock()
	return reset()
}

func Apply(rules []Rule, drops DropConfig) (err error) {
	lock.Lock()
	defer lock.Unlock()

	if err = reset(); err != nil {
		return fmt.Errorf("error while resetting firewall: %v", err)
	}

	out, err := cmd(binary, "-A", "INPUT", "-m", "conntrack", "--ctstate", "RELATED,ESTABLISHED", "-j", "ACCEPT")
	if err != nil {
		return fmt.Errorf("error running conntrack step: %v", err)
	} else {
		log.Debug("conntrack: %s", out)
	}

	// for each rule
	for _, rule := range rules {
		protos := []string{"tcp"}
		if rule.Protocol == ProtoUDP {
			protos[0] = "udp"
		} else {
			protos = []string{"tcp", "udp"}
		}

		// for each protocol
		for _, proto := range protos {
			source := []string{"-s", rule.Address}
			if rule.AddressType == AddressRange {
				// use iprange module
				source = []string{
					"-m",
					"iprange",
					"--src-range",
					rule.Address,
				}
			}

			action := "ACCEPT"
			if rule.Type == RuleBlock {
				action = "DROP"
			}

			for _, port := range rule.Ports {
				args := []string{"-A", "INPUT"}
				args = append(args, source...)
				args = append(args, "-p", proto, "--dport", port, "-j", action)
				out, err := cmd(binary, args...)
				if err != nil {
					return fmt.Errorf("error applying rule %s.%s.%s.%s: %v",
						rule.Type,
						rule.Address,
						port,
						proto,
						err)
				} else {
					log.Debug(" %s.%s.%s.%s: %s",
						rule.Type,
						rule.Address,
						port,
						proto,
						out)
				}
			}
		}
	}

	target := "DROP"

	// enable logging?
	if drops.Log {
		log.Debug("enabling logging of dropped packets: %#v", drops)

		// just in case
		// cmd(binary, "-X", "LOGNDROP");

		if out, err = cmd(binary, "-N", "LOGNDROP"); err != nil {
			return fmt.Errorf("error creating LOGNDROP: %v", err)
		} else {
			log.Debug("LOGNDROP: %s", out)
		}

		out, err := cmd(binary,
			"-A", "LOGNDROP",
			"-m", "limit", "--limit", drops.Limit,
			"-j", "LOG",
			"--log-prefix", fmt.Sprintf("%s: ", drops.Prefix),
			"--log-level", strconv.FormatInt(int64(drops.Level), 10))
		if err != nil {
			return fmt.Errorf("error enabling logging: %v", err)
		} else {
			log.Debug("logging: %s", out)
		}

		if out, err = cmd(binary, "-A", "LOGNDROP", "-j", "DROP"); err != nil {
			return fmt.Errorf("error dropping LOGNDROP: %v", err)
		} else {
			log.Debug("dropping: %s", out)
		}

		target = "LOGNDROP"
	}

	// drop the rest
	if out, err := cmd(binary, "-A", "INPUT", "-j", target); err != nil {
		return fmt.Errorf("error running drop rule: %v", err)
	} else {
		log.Debug("drop: %s", out)
	}

	return
}
