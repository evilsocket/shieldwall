package firewall

import (
	"os"
	"path/filepath"
	"github.com/evilsocket/islazy/log"
)

var binary = "/sbin/iptables"

func SetPath(bin string) (err error) {
	if bin, err = filepath.Abs(bin); err != nil {
		return
	} else if _, err = os.Stat(bin); err != nil {
		return
	}
	binary = bin
	log.Info("using firewall %s", binary)
	return
}