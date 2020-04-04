package httpserver

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
)

type options struct {
	//TODO add storage interface
	logger *zap.Logger
}

// Option server options
type Option interface {
	apply(*options)
}

type loggerOption struct {
	Log *zap.Logger
}

var (
	logger             *zap.SugaredLogger
	metricReqProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "calculator_request_processed_total",
		Help: "The total number of processed requests",
	})
)

// HandleHello handle /hello URL
func handleHello(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Hello"))
	metricReqProcessed.Inc()
	if err != nil {
		logger.Error("Error while writing response: %v", err)
	}
	logger.Infof("Request processed from %s", r.RemoteAddr)
}

// StartServer start http server
func StartServer(address string, opts ...Option) error {
	options := options{
		logger: zap.NewNop(),
	}
	for _, o := range opts {
		o.apply(&options)
	}
	logger = options.logger.Sugar()

	http.HandleFunc("/hello", handleHello)

	return http.ListenAndServe(address, nil)
}

// WithLogger apply logger
func WithLogger(log *zap.Logger) Option {
	return loggerOption{Log: log}
}

func (l loggerOption) apply(opts *options) {
	opts.logger = l.Log
}
