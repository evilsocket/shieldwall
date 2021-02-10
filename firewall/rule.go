package firewall

import "time"

const (
	AllPorts = "1:65535"
)

type Protocol string

const (
	ProtoTCP = Protocol("tcp")
	ProtoUDP = Protocol("udp")
	ProtoAll = Protocol("all")
)

type RuleType string

const (
	RuleBlock = RuleType("block")
	RuleAllow = RuleType("allow")
)

type Rule struct {
	CreatedAt time.Time `json:"created_at"`
	TTL       int       `json:"ttl"`  // used from the api to delete expired rules
	Type      RuleType  `json:"type"` // always RuleBlock for now
	Address   string    `json:"address"`
	Protocol  Protocol  `json:"protocol"`
	Ports     []string  `json:"ports"` // strings to also allow ranges
}

func (r Rule) Expires() bool {
	return r.TTL > 0
}

func (r Rule) Expired() bool {
	return r.Expires() && time.Since(r.CreatedAt).Seconds() >= float64(r.TTL)
}
