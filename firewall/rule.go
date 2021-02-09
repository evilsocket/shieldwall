package firewall

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
	RuleBlock = RuleType("block") // not used for now
	RuleAllow = RuleType("allow")
)

type Rule struct {
	Type     RuleType `json:"type"` // always RuleBlock for now
	Address  string   `json:"address"`
	Protocol Protocol `json:"protocol"`
	Ports    []string `json:"ports"` // strings to also allow ranges
}