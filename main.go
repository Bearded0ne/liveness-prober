package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/sparrc/go-ping"
)

func check(host string, ticker time.Time) {
	// initialize pinger
	pinger, err := ping.NewPinger(host)
	if err != nil {
		panic(err)
	}

	// limit to only 1 ping request
	pinger.Count = 1
	pinger.Run() // blocks until finished
	stats := pinger.Statistics()

	// display results
	fmt.Printf("host: %s\npacket loss: %v%%\navg round-trip: %v\ntimestamp: %v\n--\n",
		stats.Addr, stats.PacketLoss, stats.AvgRtt, ticker)
}

func main() {
	// pull env vars
	var probeHosts = strings.Split(os.Getenv("PROBE_HOSTS"), ",")
	if (len(probeHosts) < 2 && probeHosts[0] == "") {
		probeHosts[0] = "www.google.com"
	}

	var probeInterval , _ = strconv.ParseInt(os.Getenv("PROBE_INTERVAL"), 0, 64)
	if (probeInterval == 0) {
		probeInterval = 1000
	}

	fmt.Printf("hosts : %v", len(probeHosts))
	fmt.Printf("interval : %v", probeInterval)

	// set ticker
	ticker := time.NewTicker(time.Duration(probeInterval) * time.Millisecond)

	fmt.Println("liveness prober :: started")
	fmt.Println("--")

	// loop through hosts and ping
	for t := range ticker.C {
		for _, host := range probeHosts {
			if (host != "") {
				check(host, t)
			}
		}
	}

	// gracefully shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		fmt.Println("liveness prober :: ended")
		os.Exit(1)
	}()
}
