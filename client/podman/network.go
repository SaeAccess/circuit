package podman

type Network interface {
	// Connect a container to the network
	Connect(opts *NetworkConnectOptions) error

	// Disconnect a container from the network
	Disconnect(opts *NetworkDisconnectOptions) error

	// Exists checks if the network exists
	Exists() error

	// Inspect the network
	Inspect() (*InspectNetworkSettings, error)

	// List() []Network

	PeekBytes() []byte

	// Reload firewall rules for one or more containers
	Reload() error

	// Remove the network (rm)
	Scrub()

	// Update updates an existing network
	Update(opts *NetworkUpdateOptions) error
}

// NetworkConnectOptions describes options for connecting
// a container to a network
type NetworkConnectOptions struct {
	Container string `json:"container,omitempty"`
	// StaticIPs for this container. Optional.
	StaticIP  string `json:"static_ip,omitempty"`
	StaticIP6 string `json:"static_ip6,omitempty"`
	// Aliases contains a list of names which the dns server should resolve
	// to this container. Should only be set when DNSEnabled is true on the Network.
	// If aliases are set but there is no dns support for this network the
	// network interface implementation should ignore this and NOT error.
	// Optional.
	Aliases []string `json:"aliases,omitempty"`
	// StaticMac for this container. Optional.
	// swagger:strfmt string
	StaticMAC string `json:"static_mac,omitempty"`
}

func (c *NetworkConnectOptions) CmdLine(network string) []string {
	var args = []string{"network", "connect"}

	args = appendSA(args, "--alias", c.Aliases)
	args = appendS(args, "--ip", c.StaticIP)
	args = appendS(args, "--ip6", c.StaticIP6)
	args = appendS(args, "--mac-address", c.StaticMAC)

	args = append(args, network, c.Container)
	return args
}

// NetworkCreateOptions describes options to create a network
type NetworkCreateOptions struct {
	DisableDNS        bool     `json:"disable_dns,omitempty"`
	Driver            string   `json:"driver,omitempty"`
	Gateways          []string `json:"gateways,omitempty"`
	Internal          bool     `json:"internal,omitempty"`
	Labels            []string `json:"labels,omitempty"`
	MacVLAN           string   `json:"mac_vlan,omitempty"`
	NetworkDNSServers []string `json:"network_dns_servers,omitempty"`
	Ranges            []string `json:"ranges,omitempty"`
	Subnets           []string `json:"subnets,omitempty"`
	Routes            []string `json:"routes,omitempty"`
	IPv6              bool     `json:"ipv6,omitempty"`
	Name              string   `json:"name,omitempty"`
	// Mapping of driver options and values.
	Options []string `json:"options,omitempty"`
	// IgnoreIfExists if true, do not fail if the network already exists
	IgnoreIfExists bool `json:"ignore_if_exists,omitempty"`
	// InterfaceName sets the NetworkInterface in the network config
	InterfaceName string `json:"interface_name,omitempty"`

	Scrub bool `json:"scrub,omitempty"`
}

func (c *NetworkCreateOptions) CmdLine() []string {
	var args = []string{"network", "create"}

	args = appendB(args, "--disable-dns", c.DisableDNS)
	args = appendSA(args, "--dns", c.NetworkDNSServers)
	args = appendS(args, "--driver", c.Driver)
	args = appendSA(args, "--gateway", c.Gateways)
	args = appendB(args, "--ignore", c.IgnoreIfExists)
	args = appendS(args, "--interface-name", c.InterfaceName)
	args = appendB(args, "--internal", c.Internal)
	args = appendSA(args, "--ip-range", c.Ranges)
	args = append(args, "--ipam-driver", c.MacVLAN) // TODO check this
	args = appendB(args, "--ipv6", c.IPv6)
	args = appendSA(args, "--label", c.Labels)
	args = appendSA(args, "--opt", c.Options)
	args = appendSA(args, "--route", c.Routes)
	args = appendSA(args, "--subnet", c.Subnets)

	args = append(args, c.Name)

	return args
}

// NetworkDisconnectOptions describes options for disconnecting
// containers from networks
type NetworkDisconnectOptions struct {
	Container string `json:"container,omitempty"`
	Force     bool   `json:"force,omitempty"`
}

func (c *NetworkDisconnectOptions) CmdLine(network string) []string {
	var args = []string{"network", "disconnect"}

	args = appendB(args, "--force", c.Force)
	args = append(args, network, c.Container)

	return args
}

type NetworkRemoveOptions struct {
	Force bool `json:"force,omitempty"`
	Time  int  `json:"time,omitempty"`
}

func (c *NetworkRemoveOptions) CmdLine(network string) []string {
	var args = []string{"network", "remove"}

	args = appendB(args, "--force", c.Force)
	args = appendI(args, "--time", c.Time)
	args = append(args, network)

	return args
}

// NetworkUpdateOptions describes options to update a network
type NetworkUpdateOptions struct {
	AddDNSServers    []string `json:"adddnsservers"`
	RemoveDNSServers []string `json:"removednsservers"`
}

func (c *NetworkUpdateOptions) CmdLine(network string) []string {
	var args = []string{"network", "update"}

	args = appendSA(args, "--dns-add", c.AddDNSServers)
	args = appendSA(args, "--dns-drop", c.RemoveDNSServers)
	args = append(args, network)

	return args
}
