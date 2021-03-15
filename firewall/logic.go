package firewall

import (
	"fmt"
	"github.com/evilsocket/islazy/log"
	"github.com/evilsocket/islazy/str"
	"github.com/evilsocket/islazy/tui"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

// TODO: implement factory pattern with generic interface in order to support more firewalls

const (
	swInputChain   = "SHIELDWALL"
	swLogDropChain = "LOGNDROP"
)

var DryRun = false

type DropConfig struct {
	Log    bool   `yaml:"log"`
	Limit  string `yaml:"limit"`
	Prefix string `yaml:"prefix"`
	Level  int    `yaml:"level"`
}

var lock = sync.Mutex{}

func cmd(bin string, args ...string) (string, error) {
	log.Debug("# %s %s", bin, args)
	if DryRun {
		log.Info("%s %s %s", tui.Dim("<dry run>"), bin, strings.Join(args, " "))
		return "", nil
	} else {
		raw, err := exec.Command(bin, args...).CombinedOutput()
		if err != nil {
			log.Warning("%s", str.Trim(string(raw)))
			return "", err
		} else {
			return str.Trim(string(raw)), nil
		}
	}
}

func reset(binary string) error {
	commands := []string{
		fmt.Sprintf("-F %s", swInputChain),
		fmt.Sprintf("-D INPUT -j %s", swInputChain),
		fmt.Sprintf("-X %s", swInputChain),
		fmt.Sprintf("-F %s", swLogDropChain),
		fmt.Sprintf("-D INPUT -j %s", swLogDropChain),
		fmt.Sprintf("-X %s", swLogDropChain),
	}

	for _, c := range commands {
		out, err := cmd(binary, strings.Split(c, " ")...)
		if err != nil {
			// this is not fatal, some chains/rules might not exist yet
			log.Warning("firewall reset: %v", err)
			continue
		} else {
			log.Debug("reset(%s): %s", c, out)
		}
	}

	return nil
}

func Reset() error {
	lock.Lock()
	defer lock.Unlock()
	return reset(binary4)
}

func allowRelated(binary string) error {
	// allow related connections for responses to local clients (such as DNS)
	out, err := cmd(binary, "-A", "INPUT", "-m", "conntrack", "--ctstate", "RELATED,ESTABLISHED", "-j", "ACCEPT")
	if err != nil {
		return fmt.Errorf("error running conntrack step: %v", err)
	} else {
		log.Debug("conntrack: %s", out)
	}
	return nil
}

func createInputChain(binary string) error {
	// create custom chain
	if out, err := cmd(binary, "-N", swInputChain); err != nil {
		return fmt.Errorf("error creating chain-shieldwall: %v", err)
	} else {
		log.Debug("error creating chain %s: %s", swInputChain, out)
	}

	// Accept everything on loopback
	if out, err := cmd(binary, "-A", swInputChain, "-i", "lo", "-j", "ACCEPT"); err != nil {
		return fmt.Errorf("error applying loopback rule: %v", err)
	} else {
		log.Debug("loopback rule applied: %s", out)
	}

	return nil
}

func createLogAndDropChain(binary string, drops DropConfig) error {
	log.Debug("enabling logging of dropped packets: %#v", drops)

	if out, err := cmd(binary, "-N", swLogDropChain); err != nil {
		return fmt.Errorf("error creating %s: %v", swLogDropChain, err)
	} else {
		log.Debug("%s: %s", swLogDropChain, out)
	}

	out, err := cmd(binary,
		"-A", swLogDropChain,
		"-m", "limit", "--limit", drops.Limit,
		"-j", "LOG",
		"--log-prefix", fmt.Sprintf("%s: ", drops.Prefix),
		"--log-level", strconv.FormatInt(int64(drops.Level), 10))
	if err != nil {
		return fmt.Errorf("error enabling logging: %v", err)
	} else {
		log.Debug("logging: %s", out)
	}

	if out, err = cmd(binary, "-A", swLogDropChain, "-j", "DROP"); err != nil {
		return fmt.Errorf("error dropping LOGNDROP: %v", err)
	} else {
		log.Debug("dropping: %s", out)
	}

	return nil
}

func directInputTo(binary string, target string) error {
	// Apply custom chain on INPUT
	if out, err := cmd(binary, "-A", "INPUT", "-j", swInputChain); err != nil {
		return fmt.Errorf("error running %s rule: %v", swInputChain, err)
	} else {
		log.Debug("%s applied: %s", swInputChain, out)
	}

	// drop the rest
	if out, err := cmd(binary, "-A", "INPUT", "-j", target); err != nil {
		return fmt.Errorf("error running drop rule: %v", err)
	} else {
		log.Debug("drop: %s", out)
	}

	return nil
}

func Apply(rules []Rule, drops DropConfig) (err error) {
	binaries := [] string {
		binary4,
		binary6,
	}

	lock.Lock()
	defer lock.Unlock()

	for _, bin := range binaries {
		// skip if not available
		if bin == "" {
			continue
		}

		if err = reset(bin); err != nil {
			return fmt.Errorf("%s: %v", bin, err)
		}

		// allow related connections for responses to local clients (such as DNS)
		if err = allowRelated(bin); err != nil {
			return fmt.Errorf("%s: %v", bin, err)
		}

		if err = createInputChain(bin); err != nil {
			return fmt.Errorf("%s: %v", bin, err)
		}
	}

	// apply each rule
	for _, rule := range rules {
		// select either the ipv4 or ipv6 binary
		bin := binary4
		if rule.IPType() == IPv6 {
			if binary6 != "" {
				bin = binary6
			} else {
				log.Warning("found IPv6 rule but non ip6tables binary, ignored: %v", rule)
				continue
			}
		}

		// for each protocol
		for _, proto := range rule.Protocols() {
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
				args := []string{"-A", swInputChain}
				args = append(args, source...)
				args = append(args, "-p", proto, "--dport", port, "-j", action)
				out, err := cmd(bin, args...)
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

	for _, bin := range binaries {
		// just drop by default
		target := "DROP"
		// enable logging?
		if drops.Log {
			if err = createLogAndDropChain(bin, drops); err != nil {
				return fmt.Errorf("%s: %v", bin, err)
			}
			// first log then drop
			target = swLogDropChain
		}
		// this is what drops everything else
		if err = directInputTo(bin, target); err != nil {
			return fmt.Errorf("%s: %v", bin, err)
		}
	}

	return nil
}
