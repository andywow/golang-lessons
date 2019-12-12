package ntpclient

import (
	"math/rand"
	"testing"
	"time"
)

const (
	badHost             = "127.0.0.1"
	testGoodHostTimeout = 5
	testGoodHostRetry   = 2
)

var (
	goodHosts = []string{"ru.pool.ntp.org", "time1.google.com", "time2.google.com", "time3.google.com"}
)

func TestGetTimeCorrectHost(t *testing.T) {
	rand.Seed(time.Now().Unix())
	_, err := GetTime(goodHosts[rand.Intn(len(goodHosts))], 10, 200, testGoodHostRetry, testGoodHostTimeout)
	if err != nil {
		t.Errorf("Received incorrect response from correct address: %s", err)
	}
}

// TestMain test test
func TestGetTimeInvalidHost(t *testing.T) {
	_, err := GetTime(badHost, 10, 100, 1, 0)
	if err == nil {
		t.Errorf("Received correct response from invalid address")
	}
}
