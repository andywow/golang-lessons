package sender

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/andywow/golang-lessons/lesson-calendar/internal/calendar"
	"github.com/andywow/golang-lessons/lesson-calendar/internal/calendar/msgsystem"

	"go.uber.org/zap"
)

// options
type options struct {
	logger    *zap.Logger
	msgSystem *msgsystem.MsgSystem
}

// Option server options
type Option interface {
	apply(*options)
}

// logger option
type loggerOption struct {
	Log *zap.Logger
}

// msg system option
type msgSystemOption struct {
	MsgSystem *msgsystem.MsgSystem
}

type senderMetrics struct {
	events prometheus.Counter
}

// Sender scheduler for send events
type Sender struct {
	logger        *zap.SugaredLogger
	messageSystem msgsystem.MsgSystem
	metrics       *senderMetrics
}

// Start scheduler
func (s Sender) Start(ctx context.Context, opts ...Option) {
	options := options{
		logger: zap.NewNop(),
	}
	for _, o := range opts {
		o.apply(&options)
	}
	s.logger = options.logger.Sugar()
	s.messageSystem = *options.msgSystem

	s.metrics = &senderMetrics{
		events: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "sender",
			Name:      "events_total",
			Help:      "sent events total count",
		}),
	}

	prometheus.MustRegister(s.metrics.events)

	s.logger.Info("listening for messages")

	if err := s.messageSystem.ReceiveMessages(ctx,
		func(internalCtx context.Context, message []byte) error {

			event, err := calendar.ConvertFromJSON(message)
			if err != nil {
				s.logger.Error("failed to parse message, cause: %v", err)
			}
			s.logger.Infof("recevied message with uuid %s for user %s", event.Uuid, event.Username)
			s.metrics.events.Inc()
			return nil

		}); err != nil {
		s.logger.Error("error receiveing messages: %v", err)
	}

}

// WithLogger apply logger
func WithLogger(log *zap.Logger) Option {
	return loggerOption{Log: log}
}

// WithMsgSystem apply msg system
func WithMsgSystem(msgSystem *msgsystem.MsgSystem) Option {
	return msgSystemOption{MsgSystem: msgSystem}
}

func (o loggerOption) apply(opts *options) {
	opts.logger = o.Log
}

func (o msgSystemOption) apply(opts *options) {
	opts.msgSystem = o.MsgSystem
}
