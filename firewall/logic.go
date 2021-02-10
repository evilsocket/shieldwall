package firewall

import (
	"fmt"
	"github.com/evilsocket/islazy/log"
	"github.com/evilsocket/islazy/str"
	"os/exec"
	"sync"
)

// TODO: implement factory pattern with generic interface in order to support more firewalls

var lock = sync.Mutex{}

func cmd(bin string, args ...string) (string, error) {
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
		log.Debug("reset: %s", out)
	}
	return nil
}

func Reset() error {
	lock.Lock()
	defer lock.Unlock()
	return reset()
}

func Apply(rules []Rule) (err error) {
	lock.Lock()
	defer lock.Unlock()

	if err = reset(); err != nil {
		return fmt.Errorf("error while resetting firewall: %v", err)
	}

	/*
	Make this an option?

	out, err := cmd(binary, "-A", "INPUT", "-m", "conntrack", "--ctstate", "RELATED,ESTABLISHED", "-j" , "ACCEPT")
	if err != nil {
		return fmt.Errorf("error running conntrack step: %v", err)
	} else {
		log.Debug("conntrack: %s", out)
	}
	*/

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
			action := "ACCEPT"
			if rule.Type == RuleBlock {
				action = "DROP"
			}
			for _, port := range rule.Ports {
				// for each port
				out, err := cmd(binary,
					"-A", "INPUT",
					"-s", rule.Address,
					"-p", proto,
					"--dport", port,
					"-j", action)
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

	// drop the rest
	if out, err := cmd(binary, "-A", "INPUT", "-j", "DROP"); err != nil {
		return fmt.Errorf("error running drop rule: %v", err)
	} else {
		log.Debug("drop: %s", out)
	}

	return
}
