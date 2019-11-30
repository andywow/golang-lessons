package main

import (
	"flag"
	"math/rand"
	"os"
	"testing"
	"time"
)

const (
	badHost = "127.0.0.1"
)

var (
	goodHosts = []string{"ru.pool.ntp.org", "time1.google.com", "time2.google.com", "time3.google.com"}
)

func TestGetTimeCorrectHost(t *testing.T) {
	rand.Seed(time.Now().Unix())
	_, err := GetTime(goodHosts[rand.Intn(len(goodHosts))], 10, 200, 2)
	if err != nil {
		t.Errorf("Received incorrect response from correct address: %s", err)
	}
}

// TestMain test test
func TestGetTimeInvalidHost(t *testing.T) {
	_, err := GetTime(badHost, 10, 100, 2)
	if err == nil {
		t.Errorf("Received correct response from invalid address")
	}
}

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}
