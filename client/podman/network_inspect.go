package podman

// InspectNetworkSettings holds information about the network settings of the
// container.
// Many fields are maintained only for compatibility with `docker inspect` and
// are unused within Libpod.
type InspectNetworkSettings struct {
	InspectBasicNetworkConfig

	Bridge                 string                       `json:"Bridge"`
	SandboxID              string                       `json:"SandboxID"`
	HairpinMode            bool                         `json:"HairpinMode"`
	LinkLocalIPv6Address   string                       `json:"LinkLocalIPv6Address"`
	LinkLocalIPv6PrefixLen int                          `json:"LinkLocalIPv6PrefixLen"`
	Ports                  map[string][]InspectHostPort `json:"Ports"`
	SandboxKey             string                       `json:"SandboxKey"`
	// Networks contains information on non-default networks this
	// container has joined.
	// It is a map of network name to network information.
	Networks map[string]*InspectAdditionalNetwork `json:"Networks,omitempty"`
}

// InspectBasicNetworkConfig holds basic configuration information (e.g. IP
// addresses, MAC address, subnet masks, etc) that are common for all networks
// (both additional and main).
type InspectBasicNetworkConfig struct {
	// EndpointID is unused, maintained exclusively for compatibility.
	EndpointID string `json:"EndpointID"`
	// Gateway is the IP address of the gateway this network will use.
	Gateway string `json:"Gateway"`
	// IPAddress is the IP address for this network.
	IPAddress string `json:"IPAddress"`
	// IPPrefixLen is the length of the subnet mask of this network.
	IPPrefixLen int `json:"IPPrefixLen"`
	// SecondaryIPAddresses is a list of extra IP Addresses that the
	// container has been assigned in this network.
	SecondaryIPAddresses []Address `json:"SecondaryIPAddresses,omitempty"`
	// IPv6Gateway is the IPv6 gateway this network will use.
	IPv6Gateway string `json:"IPv6Gateway"`
	// GlobalIPv6Address is the global-scope IPv6 Address for this network.
	GlobalIPv6Address string `json:"GlobalIPv6Address"`
	// GlobalIPv6PrefixLen is the length of the subnet mask of this network.
	GlobalIPv6PrefixLen int `json:"GlobalIPv6PrefixLen"`
	// SecondaryIPv6Addresses is a list of extra IPv6 Addresses that the
	// container has been assigned in this network.
	SecondaryIPv6Addresses []Address `json:"SecondaryIPv6Addresses,omitempty"`
	// MacAddress is the MAC address for the interface in this network.
	MacAddress string `json:"MacAddress"`
	// AdditionalMacAddresses is a set of additional MAC Addresses beyond
	// the first. CNI may configure more than one interface for a single
	// network, which can cause this.
	AdditionalMacAddresses []string `json:"AdditionalMACAddresses,omitempty"`
}
