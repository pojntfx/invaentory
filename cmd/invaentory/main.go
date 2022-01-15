package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/pojntfx/invaentory/pkg/enumerate"
	"github.com/pojntfx/invaentory/pkg/ping"
	"github.com/schollz/progressbar/v3"
)

func main() {
	parallel := flag.Int("parallel", 100, "Amount of pings to run in parallel")
	multicastTimeout := flag.Int("multicast-timeout", 2000, "Time in milliseconds to wait for responses for multicast (IPv6) pings")
	unicastTimeout := flag.Int("unicast-timeout", 500, "Time in milliseconds to wait for responses for unicast (IPv4) pings")

	ipv6 := flag.Bool("6", true, "Ping using ICMPv6")
	ipv4 := flag.Bool("4", true, "Ping using ICMPv4")
	exclude := flag.String("exclude", "", "Regex of addresses to exclude")

	progress := flag.Bool("progress", true, "Show progress bar on STDERR")
	verbose := flag.Bool("verbose", false, "Enable verbose logging to STDERR")

	flag.Parse()

	var filter *regexp.Regexp
	if *exclude != "" {
		filter = regexp.MustCompile(*exclude)
	}

	hosts, err := enumerate.EnumerateLocalHosts(filter, *ipv6, *ipv4)
	if err != nil {
		panic(err)
	}

	var bar *progressbar.ProgressBar
	if *progress {
		bar = progressbar.NewOptions(
			len(hosts),
			progressbar.OptionSetDescription("Pinging"),
			progressbar.OptionSetItsString("host"),
			progressbar.OptionSetWriter(os.Stderr),
			progressbar.OptionThrottle(100*time.Millisecond),
			progressbar.OptionShowCount(),
			progressbar.OptionShowIts(),
			progressbar.OptionFullWidth(),
			// VT-100 compatibility
			progressbar.OptionUseANSICodes(true),
			progressbar.OptionSetTheme(progressbar.Theme{
				Saucer:        "=",
				SaucerHead:    ">",
				SaucerPadding: " ",
				BarStart:      "[",
				BarEnd:        "]",
			}),
		)
	}

	if err := ping.PingHosts(
		context.Background(),

		hosts,

		int64(*parallel),
		time.Duration(*multicastTimeout)*time.Millisecond,
		time.Duration(*unicastTimeout)*time.Millisecond,

		func(ip string) {
			if bar != nil {
				if err := bar.Add(1); err != nil {
					panic(err)
				}
			}

			if *verbose {
				if bar != nil {
					if err := bar.Clear(); err != nil {
						panic(err)
					}
				}

				log.Println("Starting to ping host", ip)
			}
		},
		func(ip string) {
			if bar != nil {
				if err := bar.Clear(); err != nil {
					panic(err)
				}
			}

			fmt.Println(ip)
		},
	); err != nil {
		panic(err)
	}

	if bar != nil {
		bar.Clear()
	}
}
