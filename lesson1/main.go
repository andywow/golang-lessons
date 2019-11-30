package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/beevik/ntp"
)

const (
	timeFormat = "2006-01-02T15:04:05-0700"
)

var (
	serverAddress = flag.String("server.address", "ru.pool.ntp.org", "Address of NTP server")
	queryTimeout  = flag.Int("query.timeout", 30, "Query timeout")
	queryTTL      = flag.Int("query.ttl", 100, "Query TTL")
)

// GetTime get time from ntp server
func GetTime(serverAddress string, queryTimeout, queryTTL, retries int) (*time.Time, error) {
	queryOptions := ntp.QueryOptions{Timeout: time.Duration(queryTimeout) * time.Second, TTL: queryTTL}
	var (
		err      error
		response *ntp.Response
	)
	for i := 0; i < retries; i++ {
		response, err = ntp.QueryWithOptions(serverAddress, queryOptions)
		if err != nil {
			continue
		}
		return &response.Time, nil
	}
	return nil, err
}

func main() {
	flag.Parse()
	fmt.Printf("Quering %s for time\n", *serverAddress)
	ntptime, err := GetTime(*serverAddress, *queryTimeout, *queryTTL, 1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error, while query server: %s\n", err)
		os.Exit(2)
	}
	fmt.Printf("Local system time: %s\n", time.Now().Format(timeFormat))
	fmt.Printf("NTP Server time: %s\n", ntptime.Format(timeFormat))
}
