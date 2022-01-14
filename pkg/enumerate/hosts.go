package enumerate

import (
	"encoding/binary"
	"net"
	"regexp"

	"github.com/pojntfx/invaentory/pkg/config"
)

const (
	ipv6MulticastIP = "ff02::1"
)

func EnumerateLocalHosts(
	exclude *regexp.Regexp,
) ([]config.HostPingConfig, error) {
	// Get the IPv6 multicast address
	multicastDst, err := net.ResolveIPAddr("ip", ipv6MulticastIP)
	if err != nil {
		return nil, err
	}

	hosts := []config.HostPingConfig{}

	// Iterate over interfaces and their addresses to gather all available IPs
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}

		for _, addr := range addrs {
			ip, network, err := net.ParseCIDR(addr.String())
			if err != nil {
				return nil, err
			}

			// Skip loopback and link-local addresses, which can't be bound to
			if !(ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast()) {
				newTargets := []config.HostPingConfig{}

				if ip.To4() == nil && ip.To16() != nil {
					// Add the IPv6 multicast address
					newTargets = append(
						newTargets,
						config.HostPingConfig{
							Src: ip.String(),
							Dev: iface.Name,
							Dst: multicastDst,
						},
					)
				} else if ip.To4() != nil {
					// Calculate and add the pingable IPv4 addresses
					// See https://stackoverflow.com/a/60542265
					mask := binary.BigEndian.Uint32(network.Mask)
					start := binary.BigEndian.Uint32(network.IP)
					finish := (start & mask) | (mask ^ 0xffffffff)

					for i := start; i <= finish; i++ {
						rawDst := make(net.IP, 4)

						binary.BigEndian.PutUint32(rawDst, i)

						dst, err := net.ResolveIPAddr("ip", rawDst.String())
						if err != nil {
							return nil, err
						}

						newTargets = append(
							newTargets,
							config.HostPingConfig{
								Src: ip.String(),
								Dev: iface.Name,
								Dst: dst,
							},
						)
					}
				}

				// Skip excluded targets
				for _, target := range newTargets {
					skip := false
					if exclude != nil {
						skip = exclude.Match([]byte(target.Dst.String()))
					}

					if !skip {
						hosts = append(hosts, target)
					}
				}
			}
		}
	}

	return hosts, nil
}
