package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"regexp"
	"runtime"
	"time"

	"github.com/pojntfx/invaentory/pkg/enumerate"
	"github.com/pojntfx/invaentory/pkg/ping"
)

func main() {
	parallel := flag.Int("parallel", runtime.NumCPU(), "Amount of pings to run in parallel")
	multicastTimeout := flag.Int("multicast-timeout", 2000, "Time in milliseconds to wait for responses for multicast (IPv6) pings")
	unicastTimeout := flag.Int("unicast-timeout", 500, "Time in milliseconds to wait for responses for unicast (IPv4) pings")

	exclude := flag.String("exclude", "", "Regex of addresses to exclude")

	progress := flag.Bool("progress", false, "Log progress to STDERR")
	verbose := flag.Bool("verbose", false, "Enable verbose logging to STDERR")

	flag.Parse()

	var filter *regexp.Regexp
	if *exclude != "" {
		filter = regexp.MustCompile(*exclude)
	}

	hosts, err := enumerate.EnumerateLocalHosts(filter)
	if err != nil {
		panic(err)
	}

	if err := ping.PingHosts(
		context.Background(),

		hosts,

		int64(*parallel),
		time.Duration(*multicastTimeout)*time.Millisecond,
		time.Duration(*unicastTimeout)*time.Millisecond,

		func(ip string) {
			if *verbose {
				log.Println("Starting to ping host", ip)
			}
		},
		func(ip string) {
			fmt.Println(ip)
		},
		func(percentage float64) {
			if *progress {
				log.Printf("%v%%", percentage)
			}
		},
	); err != nil {
		panic(err)
	}
}
