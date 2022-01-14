package config

import "net"

type HostPingConfig struct {
	Src string
	Dev string
	Dst *net.IPAddr
}
