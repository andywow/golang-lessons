package config

import (
	"time"
)

// RequestTimeout for GRPC calls
const RequestTimeout = 10 * time.Second

// ClientOptions global client options
type ClientOptions struct {
	GRPCHost string
	GRPCPort int64
}
