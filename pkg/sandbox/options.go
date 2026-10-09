package sandbox

const DefaultVsockPort = 49983

const defaultGuestUser = "runner"

type Options struct {
	VsockPort    int    `json:"vsock_port"`
	ListenAddr   string `json:"listen_addr"`
	ControlToken string `json:"control_token"`
	GuestUser    string `json:"guest_user"`
	DataVolume   string `json:"data_volume"`
}

func (o *Options) applyDefaults() {
	if o.VsockPort == 0 {
		o.VsockPort = DefaultVsockPort
	}
	if o.GuestUser == "" {
		o.GuestUser = defaultGuestUser
	}
}
