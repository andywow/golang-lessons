package ntpclient

import (
	"fmt"
	"os"
	"time"

	"github.com/beevik/ntp"
)

// GetTime get time from ntp server
func GetTime(serverAddress string, queryTimeout, queryTTL, retries, retryTimeout int) (time.Time, error) {
	queryOptions := ntp.QueryOptions{Timeout: time.Duration(queryTimeout) * time.Second, TTL: queryTTL}
	var (
		err      error
		response *ntp.Response
	)
	for i := 0; i < retries; i++ {
		response, err = ntp.QueryWithOptions(serverAddress, queryOptions)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while query NTP Server: %s\n", err)
			time.Sleep(time.Duration(retryTimeout) * time.Second)
			continue
		}
		return response.Time, nil
	}
	return time.Time{}, err
}
