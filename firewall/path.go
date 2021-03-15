package firewall

import (
	"os"
	"path"
	"path/filepath"
	"github.com/evilsocket/islazy/log"
)

var (
	binary4 = "/sbin/iptables"
	binary6 = ""
)

func SetPath(bin string) (err error) {
	if bin, err = filepath.Abs(bin); err != nil {
		return
	} else if _, err = os.Stat(bin); err != nil {
		return
	}
	binary4 = bin
	log.Info("ipv4 firewall: %s", binary4)

	base := path.Dir(bin)
	ipv6 := filepath.Join(base, "ip6tables")

	if _, err = os.Stat(ipv6); err == nil {
		binary6 = ipv6
		log.Info("ipv6 firewall: %s", binary6)
	} else {
		log.Important("ipv6 firewall not found in %s", base)
	}

	return
}