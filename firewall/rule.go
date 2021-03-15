package firewall

import (
	"fmt"
	"github.com/evilsocket/islazy/log"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	AllPorts = "1:65535"
)

type IPType int

const (
	IPv4 = IPType(4)
	IPv6 = IPType(6)
)

type Protocol string

const (
	ProtoTCP = Protocol("tcp")
	ProtoUDP = Protocol("udp")
	ProtoAll = Protocol("all")
)

type AddressType string

const (
	AddressNone   = AddressType("")
	AddressSimple = AddressType("simple")
	AddressMask   = AddressType("mask")
	AddressRange  = AddressType("range")
)

type RuleType string

const (
	RuleBlock = RuleType("block")
	RuleAllow = RuleType("allow")
)

type Rule struct {
	CreatedAt   time.Time   `json:"created_at"`
	TTL         int         `json:"ttl"`  // used from the api to delete expired rules
	Type        RuleType    `json:"type"` // always RuleBlock for now
	Address     string      `json:"address"`
	AddressType AddressType `json:"address_type"`
	Protocol    Protocol    `json:"protocol"`
	Ports       []string    `json:"ports"` // strings to also allow ranges
	Comment     string      `json:"comment"`
}

func (r Rule) Expires() bool {
	return r.TTL > 0
}

func (r Rule) Expired() bool {
	return r.Expires() && time.Since(r.CreatedAt).Seconds() >= float64(r.TTL)
}

func (r Rule) IPType() IPType {
	if strings.Contains(r.Address, ":") {
		return IPv6
	}
	return IPv4
}

func (r Rule) Protocols() []string {
	switch r.Protocol {
	case ProtoTCP:
		return []string{string(ProtoTCP)}
	case ProtoUDP:
		return []string{string(ProtoUDP)}
	}
	return []string{
		string(ProtoTCP),
		string(ProtoUDP),
	}
}

func (r *Rule) Validate() error {
	if r.TTL < 0 {
		return fmt.Errorf("really? %d", r.TTL)
	}

	for _, port := range r.Ports {
		if strings.Index(port, ":") != -1 {
			// parse as range
			if parts := strings.Split(port, ":"); len(parts) != 2 {
				return fmt.Errorf("%s is not a valid port range", port)
			} else if from, err := strconv.ParseInt(parts[0], 10, 32); err != nil {
				return fmt.Errorf("%s is not a valid port", parts[0])
			} else if to, err := strconv.ParseInt(parts[1], 10, 32); err != nil {
				return fmt.Errorf("%s is not a valid port", parts[1])
			} else if to <= from {
				return fmt.Errorf("bad port range, %d is not >= %d", to, from)
			} else if from < 1 || from > 65535 {
				return fmt.Errorf("%d is outside the valid ports range", from)
			} else if to < 1 || to > 65535 {
				return fmt.Errorf("%d is outside the valid ports range", to)
			}
		} else {
			// parse as number
			if p, err := strconv.ParseInt(port, 10, 32); err != nil {
				return fmt.Errorf("%s is not a valid port", port)
			} else if p < 1 || p > 65535 {
				return fmt.Errorf("%d is outside the valid ports range", p)
			}
		}
	}

	if addr := net.ParseIP(r.Address); addr != nil {
		r.AddressType = AddressSimple
		log.Debug("net.ParseIP('%s') = %#v", r.Address, addr)
	} else if ip, netw, err := net.ParseCIDR(r.Address); err == nil && ip != nil && netw != nil {
		r.AddressType = AddressMask
	} else if strings.Index(r.Address, "-") != -1 {
		if parts := strings.Split(r.Address, "-"); len(parts) == 2 {
			if net.ParseIP(parts[0]) != nil && net.ParseIP(parts[1]) != nil {
				r.AddressType = AddressRange
			}
		}
	}

	log.Debug("r.AddressType is %s", r.AddressType)

	if len(r.AddressType) == 0 {
		return fmt.Errorf("%s is not a valid IP address, mask or range", r.Address)

	}

	return nil
}
