package sandbox

// DefaultVsockPort is the guest-side port the host bridges into. It must match
// the control plane's SandboxEnvdVsockPort.
const DefaultVsockPort = 49983

const defaultGuestUser = "runner"

type Options struct {
	// VsockPort is the AF_VSOCK port to accept on.
	VsockPort int `json:"vsock_port"`
	// ListenAddr, when set, serves plain TCP instead of vsock. For local runs
	// and tests only; production always uses vsock.
	ListenAddr string `json:"listen_addr"`
	// ControlToken is required on every route but /health.
	ControlToken string `json:"control_token"`
	// GuestUser owns spawned processes and resolves ~ in paths.
	GuestUser string `json:"guest_user"`
	// DataVolume is the mount point /pause-prepare releases.
	DataVolume string `json:"data_volume"`
}

func (o *Options) applyDefaults() {
	if o.VsockPort == 0 {
		o.VsockPort = DefaultVsockPort
	}
	if o.GuestUser == "" {
		o.GuestUser = defaultGuestUser
	}
}
