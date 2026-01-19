package database

// NetworkInterface represents a network interface on the agent host
type NetworkInterface struct {
	Name       string   `json:"name"`
	MTU        int      `json:"mtu"`
	MacAddress string   `json:"mac_address"`
	Addresses  []string `json:"addresses"`
	Flags      []string `json:"flags"`
}

// ResourceUsage represents the agent's current resource consumption
type ResourceUsage struct {
	CPUPercent    float64 `json:"cpu_percent"`
	MemoryPercent float64 `json:"memory_percent"`
	MemoryUsed    uint64  `json:"memory_used"`
	NumGoroutines int     `json:"num_goroutines"`
}

// AgentResources contains all resource information reported by an agent
type AgentResources struct {
	Interfaces      []NetworkInterface `json:"interfaces"`
	ActiveInterface string             `json:"active_interface"`
	Resources       ResourceUsage      `json:"resources"`
}
