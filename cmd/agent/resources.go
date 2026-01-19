package main

import (
	"net"
	"os"
	"runtime"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
)

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

func getNetworkInterfaces() ([]NetworkInterface, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var interfaces []NetworkInterface
	for _, iface := range ifaces {
		// skip loopback and interfaces that are down
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		var addresses []string
		for _, addr := range addrs {
			addresses = append(addresses, addr.String())
		}

		var flags []string
		if iface.Flags&net.FlagUp != 0 {
			flags = append(flags, "up")
		}
		if iface.Flags&net.FlagBroadcast != 0 {
			flags = append(flags, "broadcast")
		}
		if iface.Flags&net.FlagMulticast != 0 {
			flags = append(flags, "multicast")
		}
		if iface.Flags&net.FlagPointToPoint != 0 {
			flags = append(flags, "pointtopoint")
		}

		ni := NetworkInterface{
			Name:       iface.Name,
			MTU:        iface.MTU,
			MacAddress: iface.HardwareAddr.String(),
			Addresses:  addresses,
			Flags:      flags,
		}
		interfaces = append(interfaces, ni)
	}

	return interfaces, nil
}

// getActiveInterface returns the interface most likely used for external communication
func getActiveInterface(interfaces []NetworkInterface) string {
	for _, iface := range interfaces {
		// look for an interface that is up and has addresses
		hasUp := false
		for _, flag := range iface.Flags {
			if flag == "up" {
				hasUp = true
				break
			}
		}
		if hasUp && len(iface.Addresses) > 0 && iface.MacAddress != "" {
			return iface.Name
		}
	}
	if len(interfaces) > 0 {
		return interfaces[0].Name
	}
	return ""
}

func getResourceUsage() (ResourceUsage, error) {
	usage := ResourceUsage{
		NumGoroutines: runtime.NumGoroutine(),
	}

	// get current process
	pid := os.Getpid()
	proc, err := process.NewProcess(int32(pid))
	if err != nil {
		return usage, err
	}

	// get cpu percent for this process
	cpuPercent, err := proc.CPUPercent()
	if err == nil {
		usage.CPUPercent = cpuPercent
	}

	// get memory info for this process
	memInfo, err := proc.MemoryInfo()
	if err == nil && memInfo != nil {
		usage.MemoryUsed = memInfo.RSS
	}

	// get overall system memory for percentage calculation
	vmem, err := mem.VirtualMemory()
	if err == nil && vmem != nil && vmem.Total > 0 && memInfo != nil {
		usage.MemoryPercent = float64(memInfo.RSS) / float64(vmem.Total) * 100
	}

	return usage, nil
}

// collectResources gathers all resource information
func collectResources() (*AgentResources, error) {
	interfaces, err := getNetworkInterfaces()
	if err != nil {
		return nil, err
	}

	resources, err := getResourceUsage()
	if err != nil {
		// non-fatal, continue with partial data
		resources = ResourceUsage{
			NumGoroutines: runtime.NumGoroutine(),
		}
	}

	// warm up cpu measurement (first call often returns 0)
	if resources.CPUPercent == 0 {
		if pid := os.Getpid(); pid > 0 {
			if proc, err := process.NewProcess(int32(pid)); err == nil {
				if cpuPercent, err := cpu.Percent(0, false); err == nil && len(cpuPercent) > 0 {
					// use system-wide cpu percent as fallback
					resources.CPUPercent = cpuPercent[0]
				} else if cpuPercent, _ := proc.CPUPercent(); cpuPercent > 0 {
					resources.CPUPercent = cpuPercent
				}
			}
		}
	}

	return &AgentResources{
		Interfaces:      interfaces,
		ActiveInterface: getActiveInterface(interfaces),
		Resources:       resources,
	}, nil
}
