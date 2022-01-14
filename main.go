package main

import (
	"context"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"math"
	"net"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/digineo/go-ping"
	"golang.org/x/sync/semaphore"
)

const (
	multicastIP = "ff02::1"
)

type target struct {
	src4 string
	src6 string
	dev  string
	dst  *net.IPAddr
}

func main() {
	multicastTimeout := flag.Int("multicastTimeout", 2000, "Time in milliseconds to wait for responses for multicast (IPv6) pings")
	unicastTimeout := flag.Int("unicastTimeout", 500, "Time in milliseconds to wait for responses for unicast (IPv4) pings")
	verbose := flag.Bool("verbose", false, "Enable verbose logging")
	progress := flag.Bool("progress", false, "Show progress")
	parallel := flag.Int("parallel", runtime.NumCPU(), "Amount of pings to run in parallel")
	exclude := flag.String("exclude", "", "Regex of addresses to exclude")

	flag.Parse()

	multicastDst, err := net.ResolveIPAddr("ip", multicastIP)
	if err != nil {
		panic(err)
	}

	ifaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	targets := []target{}
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			panic(err)
		}

		for _, addr := range addrs {
			ip, network, err := net.ParseCIDR(addr.String())
			if err != nil {
				panic(err)
			}

			if !(ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast()) {
				newTargets := []target{}

				if ip.To4() == nil && ip.To16() != nil {
					newTargets = append(newTargets, target{"", ip.String(), iface.Name, multicastDst})
				} else if ip.To4() != nil {
					// See https://stackoverflow.com/a/60542265
					mask := binary.BigEndian.Uint32(network.Mask)
					start := binary.BigEndian.Uint32(network.IP)
					finish := (start & mask) | (mask ^ 0xffffffff)

					for i := start; i <= finish; i++ {
						rawDst := make(net.IP, 4)

						binary.BigEndian.PutUint32(rawDst, i)

						dst, err := net.ResolveIPAddr("ip", rawDst.String())
						if err != nil {
							panic(err)
						}

						newTargets = append(newTargets, target{ip.String(), "", iface.Name, dst})
					}
				}

				for _, target := range newTargets {
					skip := false
					if *exclude != "" {
						skip, err = regexp.MatchString(*exclude, target.dst.String())
						if err != nil {
							panic(err)
						}
					}

					if !skip {
						targets = append(targets, target)
					}
				}
			}
		}
	}

	if *verbose {
		log.Printf("Starting pings for %v targets", len(targets))
	}

	if len(targets) == 0 {
		return
	}

	sem := semaphore.NewWeighted(int64(*parallel))
	done := make(chan struct{})

	curr := 0
	var currLock sync.Mutex
	total := len(targets)
	out := []string{}
	var outLock sync.Mutex

	for i, t := range targets {
		i, t := i, t

		go func(t target, i int, verbose *bool, progress *bool, curr *int, currLock *sync.Mutex, total int, outLock *sync.Mutex) {
			if err := sem.Acquire(context.Background(), 1); err != nil {
				panic(err)
			}

			defer func() {
				currLock.Lock()

				oldCurr := math.Floor(((float64(*curr) / float64(total)) * 100))
				*curr++
				newCurr := math.Floor(((float64(*curr) / float64(total)) * 100))

				if newCurr > oldCurr && *progress {
					log.Printf("%v%%", newCurr)
				}

				currLock.Unlock()

				sem.Release(1)

				if i >= total-1 {
					done <- struct{}{}
				}
			}()

			if *verbose {
				src := t.src4
				if src == "" {
					src = t.src6
				}

				log.Printf("Pinging destination %v from source %v using device %v", t.dst, src, t.dev)
			}

			pinger, err := ping.New(t.src4, t.src6)
			if err != nil {
				panic(err)
			}
			defer pinger.Close()

			ctx, cancel := context.WithTimeout(context.Background(), func() time.Duration {
				if t.src4 == "" {
					return time.Duration(*multicastTimeout) * time.Millisecond
				}

				return time.Duration(*unicastTimeout) * time.Millisecond
			}())
			defer cancel()

			results, err := pinger.PingMulticastContext(ctx, t.dst)
			if err != nil {
				if strings.HasSuffix(err.Error(), "sendto: invalid argument") {
					return
				}

				panic(err)
			}

			for res := range results {
				addr := ""

				if res.Address.IsLinkLocalUnicast() || res.Address.IsLinkLocalMulticast() {
					addr = res.Address.String() + "%" + t.dev
				} else {
					addr = res.Address.String()
				}

				outLock.Lock()
				out = append(out, addr)
				outLock.Unlock()

				if *verbose {
					log.Println("Found reachable address", addr)
				}
			}
		}(t, i, verbose, progress, &curr, &currLock, total, &outLock)
	}

	<-done

	for _, found := range out {
		fmt.Println(found)
	}
}
