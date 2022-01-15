package ping

import (
	"context"
	"net"
	"strings"
	"time"

	"github.com/digineo/go-ping"
	"github.com/pojntfx/invaentory/pkg/config"
	"golang.org/x/sync/semaphore"
)

func PingHosts(
	ctx context.Context,

	hosts []config.HostPingConfig,

	parallel int64,
	multicastTimeout time.Duration,
	unicastTimeout time.Duration,

	onStartPing func(ip string),
	onSuccessfulPing func(ip string),
) error {
	// As we can't have a semaphore with a zero weight, return
	if len(hosts) == 0 {
		return nil
	}

	// Limit the amount of concurrent scans
	sem := semaphore.NewWeighted(parallel)
	done := make(chan error)

	for i, host := range hosts {
		i, host := i, host

		go func(t config.HostPingConfig, i int, hosts []config.HostPingConfig) {
			if err := sem.Acquire(ctx, 1); err != nil {
				done <- err

				return
			}
			defer func() {
				sem.Release(1)

				// Done
				if i >= len(hosts)-1 {
					done <- nil
				}
			}()

			if onStartPing != nil {
				onStartPing(t.Dst.String())
			}

			// Determine whether the source IP is IPv6 or IPv4
			src4 := ""
			src6 := ""
			if ip := net.ParseIP(t.Src); ip.To4() == nil {
				src6 = t.Src
			} else {
				src4 = t.Src
			}

			pinger, err := ping.New(src4, src6)
			if err != nil {
				done <- err

				return
			}
			defer pinger.Close()

			timeout, cancel := context.WithTimeout(ctx, func() time.Duration {
				if src4 == "" {
					return multicastTimeout
				}

				return unicastTimeout
			}())
			defer cancel()

			// Start the ping
			replies, err := pinger.PingMulticastContext(timeout, t.Dst)
			if err != nil {
				if strings.HasSuffix(err.Error(), "sendto: invalid argument") {
					return
				}

				done <- err

				return
			}

			// Format the rplies
			for reply := range replies {
				if onSuccessfulPing != nil {
					addr := ""
					if reply.Address.IsLinkLocalUnicast() || reply.Address.IsLinkLocalMulticast() {
						addr = reply.Address.String() + "%" + t.Dev
					} else {
						addr = reply.Address.String()
					}

					onSuccessfulPing(addr)
				}
			}
		}(host, i, hosts)
	}

	return <-done
}
