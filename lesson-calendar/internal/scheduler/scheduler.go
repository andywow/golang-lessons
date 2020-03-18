package scheduler

import (
	"context"
	"time"

	"github.com/andywow/golang-lessons/lesson-calendar/internal/calendar"
	"github.com/andywow/golang-lessons/lesson-calendar/internal/calendar/msgsystem"

	"github.com/andywow/golang-lessons/lesson-calendar/internal/calendar/repository"
	"go.uber.org/zap"
)

// options
type options struct {
	logger    *zap.Logger
	storage   *repository.EventRepository
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

// storage option
type repositoryOption struct {
	EventStorage *repository.EventRepository
}

// storage option
type msgSystemOption struct {
	MsgSystem *msgsystem.MsgSystem
}

// Scheduler scheduler for send events
type Scheduler struct {
	logger        *zap.SugaredLogger
	eventStorage  repository.EventRepository
	messageSystem msgsystem.MsgSystem
}

func (s Scheduler) sendEvents(ctx context.Context, date time.Time) error {
	events, err := s.eventStorage.GetEventsForNotification(ctx, date)
	if err != nil {
		s.logger.Error(err)
		return err
	}
	s.logger.Infof("received %d events for send notifications", len(events))
	for _, event := range events {
		jsonMessage, err := calendar.ConvertToJSON(event)
		if err != nil {
			s.logger.Error(err)
			continue
		}
		if err := s.messageSystem.SendMessage(ctx, jsonMessage); err != nil {
			s.logger.Errorf("error, while sending notification about uuid: %s, err: %v",
				event.Uuid, err)
		} else {
			s.logger.Infof("notification about event %s was sent to queue", event.Uuid)
		}
	}
	return err
}

// Start scheduler
func (s Scheduler) Start(ctx context.Context, opts ...Option) {
	options := options{
		logger: zap.NewNop(),
	}
	for _, o := range opts {
		o.apply(&options)
	}
	s.logger = options.logger.Sugar()
	s.eventStorage = *options.storage
	s.messageSystem = *options.msgSystem

	nextTime := time.Now().Truncate(time.Minute).Add(time.Minute)
	s.logger.Infof("first event scan will be at %s", nextTime.Format("2006-01-02 15:04:05"))

	for {
		time.Sleep(time.Until(nextTime))
		select {
		case <-ctx.Done():
			return
		default:
			go s.sendEvents(ctx, nextTime)
		}
		nextTime = nextTime.Add(time.Minute)
	}
}

// WithLogger apply logger
func WithLogger(log *zap.Logger) Option {
	return loggerOption{Log: log}
}

// WithRepository apply storage
func WithRepository(repository *repository.EventRepository) Option {
	return repositoryOption{EventStorage: repository}
}

// WithMsgSystem apply msg system
func WithMsgSystem(msgSystem *msgsystem.MsgSystem) Option {
	return msgSystemOption{MsgSystem: msgSystem}
}

func (o loggerOption) apply(opts *options) {
	opts.logger = o.Log
}

func (o repositoryOption) apply(opts *options) {
	opts.storage = o.EventStorage
}

func (o msgSystemOption) apply(opts *options) {
	opts.msgSystem = o.MsgSystem
}
