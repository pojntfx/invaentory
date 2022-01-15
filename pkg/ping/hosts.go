package ping

import (
	"context"
	"math"
	"net"
	"strings"
	"sync"
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
	onProgress func(percentage float64),
) error {
	// As we can't have a semaphore with a zero weight, return
	if len(hosts) == 0 {
		return nil
	}

	// Limit the amount of concurrent scans
	sem := semaphore.NewWeighted(parallel)
	done := make(chan error)

	// Take note of the current progress
	curr := 0
	var currLock sync.Mutex
	total := len(hosts)

	for i, host := range hosts {
		i, host := i, host

		if i == 0 && onProgress != nil {
			onProgress(0)
		}

		go func(t config.HostPingConfig, i int, curr *int, currLock *sync.Mutex, total int) {
			if err := sem.Acquire(ctx, 1); err != nil {
				done <- err

				return
			}
			defer func() {
				// Call the `onProgress` callback if progress has incremented
				if onProgress != nil {
					currLock.Lock()

					oldCurr := math.Floor(((float64(*curr) / float64(total)) * 100))
					*curr++
					newCurr := math.Floor(((float64(*curr) / float64(total)) * 100))

					if newCurr > oldCurr {
						onProgress(newCurr)
					}

					currLock.Unlock()
				}

				sem.Release(1)

				// Done
				if i >= total-1 {
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
		}(host, i, &curr, &currLock, total)
	}

	return <-done
}
