package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/andywow/golang-lessons/lesson1/ntpclient"
)

const (
	timeFormat = "2006-01-02T15:04:05-0700"
)

var (
	serverAddress = flag.String("server.address", "ru.pool.ntp.org", "Address of NTP server")
	queryTimeout  = flag.Int("query.timeout", 30, "Query timeout")
	queryTTL      = flag.Int("query.ttl", 100, "Query TTL")
)

func main() {
	flag.Parse()
	fmt.Printf("Quering %s for time\n", *serverAddress)
	ntpTime, err := ntpclient.GetTime(*serverAddress, *queryTimeout, *queryTTL, 1, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error, while query server: %s\n", err)
		os.Exit(2)
	}
	fmt.Printf("Local system time: %s\n", time.Now().Format(timeFormat))
	fmt.Printf("NTP Server time: %s\n", ntpTime.Format(timeFormat))
}
