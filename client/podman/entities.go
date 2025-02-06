package podman

import (
	"encoding/gob"
	"fmt"

	"golang.org/x/exp/constraints"
)

func init() {
	gob.Register(ContainerCreateOptions{})
}

func appendS(args []string, opt string, v string) []string {
	if v != "" {
		args = append(args, opt, v)
	}

	return args
}

func appendSA(args []string, opt string, v []string) []string {
	for _, s := range v {
		args = append(args, opt, s)
	}

	return args
}

func appendI[T constraints.Ordered](args []string, opt string, v T) []string {
	var zv T
	if v > zv {
		args = append(args, opt, fmt.Sprintf("%v", v))
	}

	return args
}

func appendB(args []string, opt string, v bool) []string {
	if v {
		args = append(args, opt)
	}

	return args
}

type ContainerRunOptions struct {
	ContainerCreateOptions
	Detach     bool   `json:"detach,omitempty"`
	DetachKeys string `json:"detach_keys,omitempty"`
	SigProxy   bool   `json:"sig_proxy,omitempty"`
	Passwd     bool   `json:"passwd,omitempty"`
}

func (c *ContainerRunOptions) CmdLine(name string) []string {
	var args = []string{"container", "run"}
	args = append(args, c.ContainerCreateOptions.CmdLine(name)[2:]...)

	// rest of run options
	args = appendB(args, "--detach", c.Detach)
	args = appendS(args, "--detach-keys", c.DetachKeys)
	args = appendB(args, "--sig-proxy", c.SigProxy)
	args = appendB(args, "--passwd", c.Passwd)

	return args
}

type ContainerCreateOptions struct {
	Annotation            []string `json:"annotation,omitempty"`
	Attach                []string `json:"attach,omitempty"`
	Authfile              string   `json:"authfile,omitempty"`
	BlkIOWeight           string   `json:"blk_io_weight,omitempty"`
	BlkIOWeightDevice     []string `json:"blk_io_weight_device,omitempty"`
	CapAdd                []string `json:"cap_add,omitempty"`
	CapDrop               []string `json:"cap_drop,omitempty"`
	CgroupNS              string   `json:"cgroup_ns,omitempty"`
	CgroupsMode           string   `json:"cgroups_mode,omitempty"`
	CgroupParent          string   `json:"cgroup_parent,omitempty"`
	CIDFile               string   `json:"cid_file,omitempty"`
	ConmonPIDFile         string   `json:"container_conmon_pidfile,omitempty"`
	CPUPeriod             uint64   `json:"cpu_period,omitempty"`
	CPUQuota              int64    `json:"cpu_quota,omitempty"`
	CPURTPeriod           uint64   `json:"cpurt_period,omitempty"`
	CPURTRuntime          int64    `json:"cpurt_runtime,omitempty"`
	CPUShares             uint64   `json:"cpu_shares,omitempty"`
	CPUS                  float64  `json:"cpus,omitempty"`
	CPUSetCPUs            string   `json:"cpuset_cpus,omitempty"`
	CPUSetMems            string   `json:"cpu_set_mems,omitempty"`
	Devices               []string `json:"devices,omitempty"`
	DeviceCgroupRule      []string `json:"device_cgroup_rule,omitempty"`
	DeviceReadBPs         []string `json:"device_read_bps,omitempty"`
	DeviceReadIOPs        []string `json:"device_read_io_ps,omitempty"`
	DeviceWriteBPs        []string `json:"device_write_b_ps,omitempty"`
	DeviceWriteIOPs       []string `json:"device_write_io_ps,omitempty"`
	Entrypoint            *string  `json:"container_command,omitempty"`
	Env                   []string `json:"env,omitempty"`
	EnvHost               bool     `json:"env_host,omitempty"`
	EnvFile               []string `json:"env_file,omitempty"`
	Expose                []string `json:"expose,omitempty"`
	GIDMap                []string `json:"gid_map,omitempty"`
	GPUs                  []string `json:"gp_us,omitempty"`
	GroupAdd              []string `json:"group_add,omitempty"`
	HealthCmd             string   `json:"health_cmd,omitempty"`
	HealthInterval        string   `json:"health_interval,omitempty"`
	HealthRetries         uint     `json:"health_retries,omitempty"`
	HealthLogDestination  string   `json:"health_log_destination,omitempty"`
	HealthMaxLogCount     uint     `json:"health_max_log_count,omitempty"`
	HealthMaxLogSize      uint     `json:"health_max_log_size,omitempty"`
	HealthStartPeriod     string   `json:"health_start_period,omitempty"`
	HealthStartupCmd      string   `json:"health_startup_cmd,omitempty"`
	HealthStartupInterval string   `json:"health_startup_interval,omitempty"`
	HealthStartupRetries  uint     `json:"health_startup_retries,omitempty"`
	HealthStartupSuccess  uint     `json:"health_startup_success,omitempty"`
	HealthStartupTimeout  string   `json:"health_startup_timeout,omitempty"`
	HealthTimeout         string   `json:"health_timeout,omitempty"`
	HealthOnFailure       string   `json:"health_on_failure,omitempty"`
	Hostname              string   `json:"hostname,omitempty"`
	HTTPProxy             bool     `json:"http_proxy,omitempty"`
	HostUsers             []string `json:"host_users,omitempty"`
	Image                 string   `json:"image,omitempty"`
	ImageVolume           string   `json:"image_volume,omitempty"`
	Init                  bool     `json:"init,omitempty"`
	InitContainerType     string   `json:"init_container_type,omitempty"`
	InitPath              string   `json:"init_path,omitempty"`
	IntelRdtClosID        string   `json:"intel_rdt_clos_id,omitempty"`
	Interactive           bool     `json:"interactive,omitempty"`
	IPC                   string   `json:"ipc,omitempty"`
	Label                 []string `json:"label,omitempty"`
	LabelFile             []string `json:"label_file,omitempty"`
	LogDriver             string   `json:"log_driver,omitempty"`
	LogOptions            []string `json:"log_options,omitempty"`
	Memory                string   `json:"memory,omitempty"`
	MemoryReservation     string   `json:"memory_reservation,omitempty"`
	MemorySwap            string   `json:"memory_swap,omitempty"`
	MemorySwappiness      int64    `json:"memory_swappiness,omitempty"`
	Name                  string   `json:"container_name,omitempty"`
	NoHealthCheck         bool     `json:"no_health_check,omitempty"`
	OOMKillDisable        bool     `json:"oom_kill_disable,omitempty"`
	OOMScoreAdj           *int     `json:"oom_score_adj,omitempty"`
	Arch                  string   `json:"arch,omitempty"`
	OS                    string   `json:"os,omitempty"`
	Variant               string   `json:"variant,omitempty"`
	PID                   string   `json:"pid,omitempty"`
	PIDsLimit             *int64   `json:"pi_ds_limit,omitempty"`
	Platform              string   `json:"platform,omitempty"`
	Pod                   string   `json:"pod,omitempty"`
	PodIDFile             string   `json:"pod_id_file,omitempty"`
	Personality           string   `json:"personality,omitempty"`
	PreserveFDs           uint     `json:"preserve_f_ds,omitempty"`
	PreserveFD            []uint   `json:"preserve_fd,omitempty"`
	Privileged            bool     `json:"privileged,omitempty"`
	PublishAll            bool     `json:"publish_all,omitempty"`
	Pull                  string   `json:"pull,omitempty"`
	Quiet                 bool     `json:"quiet,omitempty"`
	ReadOnly              bool     `json:"read_only,omitempty"`
	ReadWriteTmpFS        bool     `json:"read_write_tmp_fs,omitempty"`
	Restart               string   `json:"restart,omitempty"`
	Replace               bool     `json:"replace,omitempty"`
	Requires              []string `json:"requires,omitempty"`
	Retry                 *uint    `json:"retry,omitempty"`
	RetryDelay            string   `json:"retry_delay,omitempty"`
	Rm                    bool     `json:"rm,omitempty"`
	RootFS                bool     `json:"root_fs,omitempty"`
	Secrets               []string `json:"secrets,omitempty"`
	SecurityOpt           []string `json:"security_opt,omitempty"`
	SdNotifyMode          string   `json:"sd_notify_mode,omitempty"`
	ShmSize               string   `json:"shm_size,omitempty"`
	ShmSizeSystemd        string   `json:"shm_size_systemd,omitempty"`
	SignaturePolicy       string   `json:"signature_policy,omitempty"`
	StopSignal            string   `json:"stop_signal,omitempty"`
	StopTimeout           uint     `json:"stop_timeout,omitempty"`
	StorageOpts           []string `json:"storage_opts,omitempty"`
	SubGIDName            string   `json:"sub_gid_name,omitempty"`
	SubUIDName            string   `json:"sub_uid_name,omitempty"`
	Sysctl                []string `json:"sysctl,omitempty"`
	Systemd               string   `json:"systemd,omitempty"`
	Timeout               uint     `json:"timeout,omitempty"`
	TLSVerify             bool     `json:"tls_verify,omitempty"`
	TmpFS                 []string `json:"tmp_fs,omitempty"`
	TTY                   bool     `json:"tty,omitempty"`
	Timezone              string   `json:"timezone,omitempty"`
	Umask                 string   `json:"umask,omitempty"`
	EnvMerge              []string `json:"env_merge,omitempty"`
	UnsetEnv              []string `json:"unset_env,omitempty"`
	UnsetEnvAll           bool     `json:"unset_env_all,omitempty"`
	UIDMap                []string `json:"uid_map,omitempty"`
	Ulimit                []string `json:"ulimit,omitempty"`
	User                  string   `json:"user,omitempty"`
	UserNS                string   `json:"-,omitempty"`
	UTS                   string   `json:"uts,omitempty"`
	Mount                 []string `json:"mount,omitempty"`
	Volume                []string `json:"volume,omitempty"`
	VolumesFrom           []string `json:"volumes_from,omitempty"`
	Workdir               string   `json:"workdir,omitempty"`
	SeccompPolicy         string   `json:"seccomp_policy,omitempty"`
	PidFile               string   `json:"pid_file,omitempty"`
	ChrootDirs            []string `json:"chroot_dirs,omitempty"`
	IsInfra               bool     `json:"is_infra,omitempty"`
	IsClone               bool     `json:"is_clone,omitempty"`
	DecryptionKeys        []string `json:"decryption_keys,omitempty"`
	CgroupConf            []string `json:"cgroup_conf,omitempty"`
	GroupEntry            string   `json:"group_entry,omitempty"`
	PasswdEntry           string   `json:"passwd_entry,omitempty"`

	// Network options
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

	Command string   `json:"command,omitempty"`
	Args    []string `json:"args,omitempty"`

	// If Scrub is set, the container element will automatically be removed from its anchor
	// when the process exits.
	Scrub bool
}

func (c *ContainerCreateOptions) CmdLine(name string) []string {
	args := append([]string{}, "container", "create")

	// Add name
	args = append(args, "--name", c.Name)

	args = appendSA(args, "--add-host", c.AddHosts)
	args = appendSA(args, "--annotation", c.Annotation)
	args = appendS(args, "--arch", c.Arch)
	args = appendSA(args, "--attach", c.Attach)
	args = appendS(args, "--authfile", c.Authfile)
	args = appendS(args, "--blkio-weight", c.BlkIOWeight)
	args = appendSA(args, "--blkio-weight-device", c.BlkIOWeightDevice)
	args = appendSA(args, "--cap-add", c.CapAdd)
	args = appendSA(args, "--cap-drop", c.CapDrop)
	args = appendSA(args, "--cgroup-conf", c.CgroupConf)
	args = appendS(args, "--cgroupns", c.CgroupNS)
	args = appendS(args, "--cgroups", c.CgroupsMode)
	args = appendS(args, "--cgroup-parent", c.CgroupParent)
	args = appendSA(args, "--chrootdirs", c.ChrootDirs)
	args = appendS(args, "--cidfile", c.CIDFile)
	args = appendS(args, "--conmon-pidfile", c.ConmonPIDFile)
	args = appendI(args, "--cpu-period", c.CPUPeriod)
	args = appendI(args, "--cpu-quota", c.CPUQuota)
	args = appendI(args, "--cpu-rt-period", c.CPURTPeriod)
	args = appendI(args, "--cpu-rt-runtime", c.CPURTRuntime)
	args = appendI(args, "--cpu-shares", c.CPUShares)
	args = appendI(args, "--cpus", c.CPUS)
	args = appendS(args, "--cpuset-cpus", c.CPUSetCPUs)
	args = appendS(args, "--cpuset-mems", c.CPUSetMems)
	args = appendSA(args, "--decryption-key", c.DecryptionKeys)
	args = appendSA(args, "--device", c.Devices)
	args = appendSA(args, "--device-cgroup-rule", c.DeviceCgroupRule)
	args = appendSA(args, "--device-read-bps", c.DeviceReadBPs)
	args = appendSA(args, "--device-read-iops", c.DeviceReadIOPs)
	args = appendSA(args, "--device-write-bps", c.DeviceWriteBPs)
	args = appendSA(args, "--device-write-iops", c.DeviceWriteIOPs)
	args = appendSA(args, "--dns-option", c.DNSOptions)
	args = appendSA(args, "--dns-search", c.DNSSearch)
	args = appendSA(args, "--dns", c.DNSServers)
	if c.Entrypoint != nil {
		args = appendS(args, "--entry-point", *c.Entrypoint)
	}
	args = appendSA(args, "--env", c.Env)
	args = appendB(args, "--env-host", c.EnvHost)
	args = appendSA(args, "--env-file", c.EnvFile)
	args = appendSA(args, "--env-merge", c.EnvMerge)
	args = appendSA(args, "--expose", c.Expose)
	args = appendSA(args, "--gidmap", c.GIDMap)
	//	args = appendSA(args, "--gpus", c.GPUs)
	args = appendSA(args, "--group-add", c.GroupAdd)
	args = appendS(args, "--group-entry", c.GroupEntry)
	args = appendS(args, "--health-cmd", c.HealthCmd)
	args = appendS(args, "--health-interval", c.HealthInterval)
	args = appendS(args, "--health-on-failure", c.HealthOnFailure)
	args = appendI(args, "--health-retries", c.HealthRetries)
	//args = appendS(args, "--", c.HealthLogDestination string   `json:"health_log_destination,omitempty"`
	args = appendS(args, "--health-start-period", c.HealthStartPeriod)
	args = appendS(args, "--health-startup-cmd", c.HealthStartupCmd)
	args = appendS(args, "--health-startup-interval", c.HealthStartupInterval)
	args = appendI(args, "--health-startup-retries", c.HealthStartupRetries)
	args = appendI(args, "--health-startup-success", c.HealthStartupSuccess)
	//	args = appendS(args, "--", c.HealthMaxLogCount    uint     `json:"health_max_log_count,omitempty"`
	//	args = appendS(args, "--", c.HealthMaxLogSize     uint     `json:"health_max_log_size,omitempty"`
	args = appendS(args, "--health-startup-timeout", c.HealthStartupTimeout)
	args = appendS(args, "--health-timeout", c.HealthTimeout)
	args = appendS(args, "--hostname", c.Hostname)
	args = appendSA(args, "--host-user", c.HostUsers)
	args = appendB(args, "--http-proxy", c.HTTPProxy)
	args = appendS(args, "--image-volume", c.ImageVolume)
	args = appendB(args, "--interactive", c.Interactive)
	args = appendB(args, "--init", c.Init)
	args = appendS(args, "--init-ctr", c.InitContainerType)
	args = appendS(args, "--init-path", c.InitPath)
	args = appendS(args, "--ip", c.StaticIP)
	args = appendS(args, "--ip6", c.StaticIP6)
	args = appendS(args, "--ipc", c.IPC)
	args = appendSA(args, "--label", c.Label)
	args = appendSA(args, "--label-file", c.LabelFile)
	args = appendS(args, "--log-driver", c.LogDriver)
	args = appendSA(args, "--log-opt", c.LogOptions)
	args = appendS(args, "--mac-address", c.StaticMAC)
	args = appendS(args, "--memory", c.Memory)
	args = appendS(args, "--memory-reservation", c.MemoryReservation)
	args = appendS(args, "--memory-swap", c.MemorySwap)
	args = appendI(args, "--memory-swappiness", c.MemorySwappiness)
	args = appendSA(args, "--mount", c.Mount)
	args = appendSA(args, "--network-alias", c.Aliases)
	args = appendSA(args, "--network", c.Networks)
	args = appendB(args, "--no-healthcheck", c.NoHealthCheck)
	args = appendB(args, "--no-hosts", c.NoHosts)
	args = appendB(args, "--oom-kill-disable", c.OOMKillDisable)
	if c.OOMScoreAdj != nil {
		args = appendI(args, "--oom-score-adj", *c.OOMScoreAdj)
	}
	args = appendS(args, "--os", c.OS)
	args = appendS(args, "--passwd-entry", c.PasswdEntry)
	args = appendS(args, "--personality", c.Personality)
	args = appendS(args, "--pid", c.PID)
	args = appendI(args, "--pids-limit", *c.PIDsLimit)
	args = appendS(args, "--pidfile", c.PidFile)
	args = appendS(args, "--platform", c.Platform)
	args = appendS(args, "--pod", c.Pod)
	args = appendS(args, "--pod-id-file", c.PodIDFile)
	args = appendB(args, "--privileged", c.Privileged)
	args = appendSA(args, "--publish", c.PublishPorts)
	args = appendB(args, "--publish-all", c.PublishAll)
	args = appendS(args, "--pull", c.Pull)
	args = appendB(args, "--quiet", c.Quiet)
	args = appendS(args, "--rdt-class", c.IntelRdtClosID)
	args = appendB(args, "--read-only", c.ReadOnly)
	args = appendB(args, "--read-only-tmpfs", c.ReadWriteTmpFS)
	args = appendB(args, "--replace", c.Replace)
	args = appendSA(args, "--requires", c.Requires)
	args = appendS(args, "--restart", c.Restart)
	//args = appendI(args, "--retry", *c.Retry))
	//args = appendS(args, "--retry-delay", c.RetryDelay)
	args = appendB(args, "--rm", c.Rm)
	args = appendB(args, "--rootfs", c.RootFS)
	args = appendS(args, "--sd-notify-mode", c.SdNotifyMode)
	args = appendS(args, "--secomp-policy", c.SeccompPolicy)
	args = appendSA(args, "--secret", c.Secrets)
	args = appendSA(args, "--security-opt", c.SecurityOpt)
	args = appendS(args, "--shm-size", c.ShmSize)
	args = appendS(args, "--shm-size-systemd", c.ShmSizeSystemd)
	args = appendS(args, "--stop-signal", c.StopSignal)
	args = appendI(args, "--stop-timeout", c.StopTimeout)
	args = appendS(args, "--subgidname", c.SubGIDName)
	args = appendS(args, "--subuidname", c.SubUIDName)
	args = appendSA(args, "--sysctl", c.Sysctl)
	args = appendS(args, "--systemd", c.Systemd)
	args = appendI(args, "--timeout", c.Timeout)
	if c.TLSVerify {
		args = append(args, "--tls-verify")
	} else {
		args = append(args, "--tls-verify=false")
	}
	args = appendSA(args, "--tmpfs", c.TmpFS)
	args = appendB(args, "--tty", c.TTY)
	args = appendS(args, "--tz", c.Timezone)
	args = appendSA(args, "--uidmap", c.UIDMap)
	args = appendSA(args, "--ulimit", c.Ulimit)
	args = appendS(args, "--umask", c.Umask)
	args = appendSA(args, "--unsetenv", c.UnsetEnv)
	args = appendB(args, "--unsetenv-all", c.UnsetEnvAll)
	args = appendS(args, "--user", c.User)
	args = appendS(args, "--userns", c.UserNS)
	args = appendS(args, "--uts", c.UTS)
	//args = appendS(args, "--signature-policy", c.SignaturePolicy)
	args = appendS(args, "--variant", c.Variant)
	args = appendSA(args, "--volume", c.Volume)
	args = appendSA(args, "--volumes-from", c.VolumesFrom)
	args = appendS(args, "--workdir", c.Workdir)

	//args = appendSA(args, "--storage-opt", c.StorageOpts)
	// args = appendB(args, "--infra", c.IsInfra)
	//args = appendB(args, "--clone", c.IsClone)

	args = append(args, c.Image)
	if c.Command != "" {
		args = append(args, c.Command) // command path
		args = append(args, c.Args...)
	}

	return args

}

type ContainerCheckpointOptions struct {
	//entities.ContainerCheckpointOptions
	Compress      string `json:"compress,omitempty"`
	CreateImage   string `json:"create_image,omitempty"`
	Export        string `json:"export,omitempty"`
	FileLocks     bool   `json:"file_locks,omitempty"`
	IgnoreRootfs  bool   `json:"ignore_rootfs,omitempty"`
	IgnoreVolumes bool   `json:"ignore_volumes,omitempty"`
	LeaveRunning  bool   `json:"leave_running,omitempty"`
}

func (o *ContainerCheckpointOptions) CmdLine(name string) []string {
	args := append([]string{}, "container", "checkpoint")

	args = appendS(args, "--compress", o.Compress)
	args = appendS(args, "--create-image", o.CreateImage)
	args = appendS(args, "--export", o.Export)
	args = appendB(args, "--file-locks", o.FileLocks)
	args = appendB(args, "--ignore-rootfs", o.IgnoreRootfs)
	args = appendB(args, "--ignore-volumes", o.IgnoreVolumes)
	args = appendB(args, "--leave-running", o.LeaveRunning)

	args = append(args, name)
	return args
}

type ContainerRestoreOptions struct {
	FileLocks       bool   `json:"file_locks,omitempty"`
	IgnoreRootfs    bool   `json:"ignore_rootfs,omitempty"`
	IgnoreStaticIp  bool   `json:"ignore_static_ip,omitempty"`
	IgnoreStaticMac bool   `json:"ignore_static_mac,omitempty"`
	IgnoreVolumes   bool   `json:"ignore_volumes,omitempty"`
	Import          string `json:"import,omitempty"`
	ImportPrevious  string `json:"import_previous,omitempty"`
	Keep            bool   `json:"keep,omitempty"`
	Name            string `json:"name,omitempty"`
	Pod             string `json:"pod,omitempty"`
	Publish         string `json:"publish,omitempty"`
	TcpEstablished  bool   `json:"tcp_established,omitempty"`
}

func (o *ContainerRestoreOptions) CmdLine(name string) []string {
	args := append([]string{}, "container", "restore")
	args = appendB(args, "--file-locks", o.FileLocks)
	args = appendB(args, "--ignore-rootfs", o.IgnoreRootfs)
	args = appendB(args, "--ignore-static-ip", o.IgnoreStaticIp)
	args = appendB(args, "--ignore-static-mac", o.IgnoreStaticMac)
	args = appendB(args, "--ignore-volumes", o.IgnoreVolumes)
	args = appendS(args, "--import", o.Import)
	args = appendS(args, "--import-previous", o.ImportPrevious)
	args = appendB(args, "--keep", o.Keep)
	args = appendS(args, "--name", o.Name)
	args = appendS(args, "--pod", o.Pod)
	args = appendS(args, "--publish", o.Publish)
	args = appendB(args, "--tcp-established", o.TcpEstablished)

	args = append(args, name)

	return args
}

type ContainerCloneOptions struct {
	BlkioWeight       string  `json:"blkio_weight,omitempty"`
	BlkioWeightDevice string  `json:"blkio_weight_device,omitempty"` // format DeviceName:Weight
	CpuPeriod         uint    `json:"cpu_period,omitempty"`
	CPUQuota          int     `json:"cpu_quota,omitempty"`
	CpuRtPeriod       uint    `json:"cpu_rt_period,omitempty"`
	CpuRtRuntime      int     `json:"cpu_rt_runtime,omitempty"`
	CpuShares         int     `json:"cpu_shares,omitempty"`
	Cpus              float32 `json:"cpus,omitempty"`
	CpusetCpus        string  `json:"cpuset_cpus,omitempty"`
	CpusetMems        string  `json:"cpuset_mems,omitempty"`
	Destroy           bool    `json:"destroy,omitempty"`
	DeviceReadBps     string  `json:"device_read_bps,omitempty"`
	DeviceWriteBps    string  `json:"device_write_bps,omitempty"`
	Force             bool    `json:"force,omitempty"`
	Memory            string  `json:"memory,omitempty"`
	MemoryReservation string  `json:"memory_reservation,omitempty"`
	MemorySwap        string  `json:"memory_swap,omitempty"`
	MemorySwappiness  int     `json:"memory_swappiness,omitempty"`
	Name              string  `json:"name,omitempty"`
	Pod               string  `json:"pod,omitempty"`
	Run               bool    `json:"run,omitempty"`
}

func (o *ContainerCloneOptions) CmdLine(name string, newName string, image string) []string {
	args := append([]string{}, "container", "clone")
	args = append(args, "--name", o.Name)
	args = appendS(args, "--blkio-weight", o.BlkioWeight)
	args = appendS(args, "--blkio-weight-device", o.BlkioWeightDevice)
	args = appendI(args, "--cpu-period", o.CpuPeriod)
	args = appendI(args, "--cpu-quota", o.CPUQuota)
	args = appendI(args, "--cpu-rt-period", o.CpuRtPeriod)
	args = appendI(args, "--cpu-rt-runtime", o.CpuRtRuntime)
	args = appendI(args, "--cpu-shares", o.CpuShares)
	args = append(args, "--cpus", fmt.Sprintf("%f", o.Cpus))
	args = appendS(args, "--cpuset-cpus", o.CpusetCpus)
	args = appendS(args, "--cpuset-mems", o.CpusetMems)
	args = appendB(args, "--destroy", o.Destroy)
	args = appendS(args, "--device-read-bps", o.DeviceReadBps)
	args = appendS(args, "--device-write-bps", o.DeviceWriteBps)
	args = appendB(args, "--force", o.Force)
	args = appendS(args, "--memory", o.Memory)
	args = appendS(args, "--memory-reservation", o.MemoryReservation)
	args = appendS(args, "--memory-swap", o.MemorySwap)
	args = appendI(args, "--memory-swappiness", o.MemorySwappiness)
	args = appendS(args, "--pod", o.Pod)
	args = appendB(args, "--run", o.Run)

	return args
}

type ContainerExecOptions struct {
	Detach      bool     `json:"detach,omitempty"`
	Env         []string `json:"env,omitempty"`
	PreserveFds uint     `json:"preserve_fds,omitempty"`
	Privileged  bool     `json:"privileged,omitempty"`
	User        string   `json:"user,omitempty"`
	WorkDir     string   `json:"work_dir,omitempty"`
	Command     string   `json:"command,omitempty"`
	Args        []string `json:"args,omitempty"`
}

func (e *ContainerExecOptions) CmdLine(name string) []string {
	args := append([]string{}, "container", "exec")

	args = appendB(args, "--detach", e.Detach)
	args = appendSA(args, "--env", e.Env)
	args = appendI(args, "--preserve-fds", e.PreserveFds)
	args = appendB(args, "--privileged", e.Privileged)
	args = appendS(args, "--workdir", e.WorkDir)

	args = append(args, name)

	if e.Command != "" {
		args = append(args, e.Command)
		if len(e.Args) > 0 {
			args = append(args, e.Args...)
		}
	}

	return args
}

type ContainerStopOpts struct {
	Ignore bool
	Time   int
}

func (o *ContainerStopOpts) CmdLine(name string) []string {
	args := append([]string{}, "container", "stop")
	args = appendB(args, "--ignore", o.Ignore)
	args = appendI(args, "--time", o.Time)

	args = append(args, name)
	return args
}
