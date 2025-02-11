package podman

// Pod is the interface for managing a pod
type Pod interface {
	// Clone an existing pod
	Clone(opts *PodCloneOptions) (Pod, error)

	// Exists checks if pod exists in storage.
	Exists() error

	// Inspect the configuration of the pod
	Inspect() (*InspectPodData, error)

	//	IsDone() bool

	// List the containers in the pod
	//	List() []Container

	// Pause the pod
	Pause() error

	PeekBytes() []byte

	// Restart restarts all containers in the pod.
	Restart() error

	// Scrub removes the pod (rm)
	Scrub()

	// Signal sends the specified signal or SIGKILL to the container in the pod
	Signal(sig string) error

	// Start the pod
	Start(opts *PodStartOptions) error

	// Stop the pod
	Stop(opts *PodStopOptions) error

	// Unpause the pod
	Unpause() error
}

type PodCloneOptions struct {
	BlkIOWeight        string   `json:"blk_io_weight,omitempty"`
	BlkIOWeightDevice  []string `json:"blk_io_weight_device,omitempty"`
	CgroupParent       string   `json:"cgroup_parent,omitempty"`
	CPUs               float64  `json:"cpus,omitempty"`
	CPUShares          uint64   `json:"cpu_shares,omitempty"`
	CPUSetCPUs         string   `json:"cpuset_cpus,omitempty"`
	CPUSetMems         string   `json:"cpu_set_mems,omitempty"`
	Destroy            bool     `json:"destroy,omitempty"`
	Devices            []string `json:"devices,omitempty"`
	DeviceReadBPs      []string `json:"device_read_bps,omitempty"`
	DeviceWriteBPs     []string `json:"device_write_b_ps,omitempty"`
	GIDMap             []string `json:"gid_map,omitempty"`
	Hostname           string   `json:"hostname,omitempty"`
	InfraCommand       string   `json:"container_command,omitempty"`
	InfraConmonPidFile string   `json:"container_conmon_pidfile,omitempty"`
	InfraName          string   `json:"container_name,omitempty"`
	Labels             []string `json:"labels,omitempty"`
	LabelFile          string   `json:"label-file,omitempty"`
	Memory             string   `json:"memory,omitempty"`
	MemorySwap         string   `json:"memory_swap,omitempty"`
	Name               string   `json:"name,omitempty"`
	Pid                string   `json:"pid,omitempty"`
	Restart            string   `json:"restart,omitempty"`
	SecurityOpt        []string `json:"security_opt,omitempty"`
	ShmSize            string   `json:"shm_size,omitempty"`
	ShmSizeSystemd     string   `json:"shm_size_systemd,omitempty"`
	Start              string   `json:"start,omitempty"`
	SubGIDName         string   `json:"sub_gid_name,omitempty"`
	SubUIDName         string   `json:"sub_uid_name,omitempty"`
	Sysctl             []string `json:"sysctl,omitempty"`
	UIDMap             []string `json:"uid_map,omitempty"`
	UserNamespace      string   `json:"userns,omitempty"`
	Uts                string   `json:"uts,omitempty"`
	Volume             []string `json:"volume,omitempty"`
	VolumesFrom        []string `json:"volumes_from,omitempty"`
}

func (c *PodCloneOptions) CmdLine() []string {
	var args = []string{"pod", "clone"}

	args = append(args, "--name", c.Name)

	args = appendS(args, "--blkio-weight", c.BlkIOWeight)
	args = appendSA(args, "--blkio-weight-device", c.BlkIOWeightDevice)
	args = appendS(args, "--cgroup-parent", c.CgroupParent)
	args = appendI(args, "--cpu-shares", c.CPUShares)
	args = appendI(args, "--cpus", c.CPUs)
	args = appendS(args, "--cpuset-cpus", c.CPUSetCPUs)
	args = appendS(args, "--cpuset-mems", c.CPUSetMems)
	args = appendB(args, "--destroy", c.Destroy)
	args = appendSA(args, "--device", c.Devices)
	args = appendSA(args, "--device-read-bps", c.DeviceReadBPs)
	args = appendSA(args, "--device-write-bps", c.DeviceWriteBPs)
	args = appendSA(args, "--gidmap", c.GIDMap)
	args = appendS(args, "--hostname", c.Hostname)
	args = appendS(args, "--infra-command", c.InfraCommand)
	args = appendS(args, "--infra-conmon-pidfile", c.InfraConmonPidFile)
	args = appendS(args, "--infra-name", c.InfraName)
	args = appendSA(args, "--label", c.Labels)
	args = appendS(args, "--label-file", c.LabelFile)
	args = appendS(args, "--memory", c.Memory)
	args = appendS(args, "--memory-swap", c.MemorySwap)
	args = appendS(args, "--pid", c.Pid)
	args = appendS(args, "--restart", c.Restart)
	args = appendSA(args, "--security-opt", c.SecurityOpt)
	args = appendS(args, "--shm-size", c.ShmSize)
	args = appendS(args, "--shm-size-systemd", c.ShmSizeSystemd)
	args = appendS(args, "--subgidname", c.SubGIDName)
	args = appendS(args, "--subuidname", c.SubUIDName)
	args = appendSA(args, "--sysctl", c.Sysctl)
	args = appendSA(args, "--uidmap", c.UIDMap)
	args = appendS(args, "--userns", c.UserNamespace)
	args = appendS(args, "--uts", c.Uts)
	args = appendSA(args, "--volume", c.Volume)
	args = appendSA(args, "--volumes-from", c.VolumesFrom)

	return args

}

type PodCreateOptions struct {
	BlkIOWeight        string   `json:"blk_io_weight,omitempty"`
	BlkIOWeightDevice  []string `json:"blk_io_weight_device,omitempty"`
	CgroupParent       string   `json:"cgroup_parent,omitempty"`
	CPUs               float64  `json:"cpus,omitempty"`
	CPUShares          uint64   `json:"cpu_shares,omitempty"`
	CPUSetCPUs         string   `json:"cpuset_cpus,omitempty"`
	CPUSetMems         string   `json:"cpu_set_mems,omitempty"`
	Devices            []string `json:"devices,omitempty"`
	DeviceReadBPs      []string `json:"device_read_bps,omitempty"`
	DeviceWriteBPs     []string `json:"device_write_b_ps,omitempty"`
	ExitPolicy         string   `json:"exit_policy,omitempty"`
	GIDMap             []string `json:"gid_map,omitempty"`
	Hostname           string   `json:"hostname,omitempty"`
	Infra              bool     `json:"infra,omitempty"`
	InfraImage         string   `json:"infra_image,omitempty"`
	InfraName          string   `json:"container_name,omitempty"`
	InfraCommand       string   `json:"container_command,omitempty"`
	InfraConmonPidFile string   `json:"container_conmon_pidfile,omitempty"`
	Labels             []string `json:"labels,omitempty"`
	LabelFile          string   `json:"label-file,omitempty"`
	Memory             string   `json:"memory,omitempty"`
	MemorySwap         string   `json:"memory_swap,omitempty"`
	Name               string   `json:"name,omitempty"`
	Share              string   `json:"share,omitempty"`
	ShareParent        *bool    `json:"share_parent,omitempty"`
	Pid                string   `json:"pid,omitempty"`
	PidIDFile          string   `json:"pid_id_file,omitempty"`
	Replace            bool     `json:"replace,omitempty"`
	Restart            string   `json:"restart,omitempty"`
	ShmSize            string   `json:"shm_size,omitempty"`
	ShmSizeSystemd     string   `json:"shm_size_systemd,omitempty"`
	Volume             []string `json:"volume,omitempty"`
	VolumesFrom        []string `json:"volumes_from,omitempty"`
	SecurityOpt        []string `json:"security_opt,omitempty"`
	SubGIDName         string   `json:"sub_gid_name,omitempty"`
	SubUIDName         string   `json:"sub_uid_name,omitempty"`
	Sysctl             []string `json:"sysctl,omitempty"`
	UserNamespace      string   `json:"userns,omitempty"`
	Uts                string   `json:"uts,omitempty"`
	UIDMap             []string `json:"uid_map,omitempty"`

	AddHosts     []string `json:"add_host,omitempty"`
	Aliases      []string `json:"network_alias,omitempty"`
	StaticIP     string   `json:"static_ip,omitempty"`
	StaticIP6    string   `json:"static_ip6,omitempty"`
	StaticMAC    string   `json:"static_mac,omitempty"`
	DNSOptions   []string `json:"dns_option,omitempty"`
	DNSSearch    []string `json:"dns_search,omitempty"`
	DNSServers   []string `json:"dns_server,omitempty"`
	NoHosts      bool     `json:"no_manage_hosts,omitempty"`
	PublishPorts []string `json:"published_ports,omitempty"`
	Networks     []string `json:"network,omitempty"`

	Scrub bool `json:"scrub,omitempty"`
}

func (c *PodCreateOptions) CmdLine() []string {
	var args = []string{"pod", "create"}

	args = append(args, "--name", c.Name)

	args = appendSA(args, "--add-host", c.AddHosts)
	args = appendS(args, "--blkio-weight", c.BlkIOWeight)
	args = appendSA(args, "--blkio-weight-device", c.BlkIOWeightDevice)
	args = appendS(args, "--cgroup-parent", c.CgroupParent)
	args = appendI(args, "--cpu-shares", c.CPUShares)
	args = appendI(args, "--cpus", c.CPUs)
	args = appendS(args, "--cpuset-cpus", c.CPUSetCPUs)
	args = appendS(args, "--cpuset-mems", c.CPUSetMems)
	args = appendSA(args, "--device", c.Devices)
	args = appendSA(args, "--device-read-bps", c.DeviceReadBPs)
	args = appendSA(args, "--device-write-bps", c.DeviceWriteBPs)
	args = appendSA(args, "--dns-option", c.DNSOptions)
	args = appendSA(args, "--dns-search", c.DNSSearch)
	args = appendSA(args, "--dns", c.DNSServers)
	args = appendS(args, "--exit-policy", c.ExitPolicy)
	args = appendSA(args, "--gidmap", c.GIDMap)
	args = appendS(args, "--hostname", c.Hostname)
	args = appendB(args, "--infra", c.Infra)
	args = appendS(args, "--infra-command", c.InfraCommand)
	args = appendS(args, "--infra-conmon-pidfile", c.InfraConmonPidFile)
	args = appendS(args, "--infra-image", c.InfraImage)
	args = appendS(args, "--infra-name", c.InfraName)
	args = appendS(args, "--ip", c.StaticIP)
	args = appendS(args, "--ip6", c.StaticIP6)
	args = appendSA(args, "--label", c.Labels)
	args = appendS(args, "--label-file", c.LabelFile)
	args = appendS(args, "--mac-address", c.StaticMAC)
	args = appendS(args, "--memory", c.Memory)
	args = appendS(args, "--memory-swap", c.MemorySwap)
	args = appendSA(args, "--network-alias", c.Aliases)
	args = appendSA(args, "--network", c.Networks)
	args = appendB(args, "--no-hosts", c.NoHosts)
	args = appendS(args, "--pid", c.Pid)
	args = appendS(args, "--pod-id-file", c.PidIDFile)
	args = appendSA(args, "--publish", c.PublishPorts)
	args = appendB(args, "--replace", c.Replace)
	args = appendS(args, "--restart", c.Restart)
	args = appendSA(args, "--security-opt", c.SecurityOpt)
	args = appendS(args, "--share", c.Share)
	if c.ShareParent != nil {
		args = appendB(args, "--share-parent", *c.ShareParent)
	}
	args = appendS(args, "--shm-size", c.ShmSize)
	args = appendS(args, "--shm-size-systemd", c.ShmSizeSystemd)
	args = appendS(args, "--subgidname", c.SubGIDName)
	args = appendS(args, "--subuidname", c.SubUIDName)
	args = appendSA(args, "--sysctl", c.Sysctl)
	args = appendSA(args, "--uidmap", c.UIDMap)
	args = appendS(args, "--userns", c.UserNamespace)
	args = appendS(args, "--uts", c.Uts)
	//args = appendS(args, "--signature-policy", c.SignaturePolicy)
	args = appendSA(args, "--volume", c.Volume)
	args = appendSA(args, "--volumes-from", c.VolumesFrom)

	return args
}

type PodStartOptions struct {
	PodIdFile string `json:"pod_id_file,omitempty"`
}

func (c *PodStartOptions) CmdLine(name string) []string {
	var args = []string{"pod", "start"}
	args = appendS(args, "--pod-id-file", c.PodIdFile)
	args = append(args, name)
	return args
}

type PodStopOptions struct {
	Ignore    bool   `json:"ignore,omitempty"`
	PodIdFile string `json:"pod_id_file,omitempty"`
	Time      int    `json:"time,omitempty"`
}

func (c *PodStopOptions) CmdLine(name string) []string {
	var args = []string{"pod", "stop"}
	args = appendB(args, "--ignore", c.Ignore)
	args = appendS(args, "--pod-id-file", c.PodIdFile)
	args = appendI(args, "--time", c.Time)
	args = append(args, name)
	return args
}
