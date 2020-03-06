package grpcserver

import (
	"context"
	"net"
	"time"

	"github.com/andywow/golang-lessons/lesson-calendar/internal/calendar/repository"
	"github.com/andywow/golang-lessons/lesson-calendar/pkg/eventapi"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// options
type options struct {
	logger  *zap.Logger
	storage *repository.EventRepository
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

// APIServer api server struct for grpc client interface
type APIServer struct {
	logger       *zap.SugaredLogger
	eventStorage repository.EventRepository
}

// CreateEvent create event
func (s APIServer) CreateEvent(ctx context.Context, in *eventapi.Event) (*eventapi.Event, error) {
	err := s.eventStorage.CreateEvent(ctx, in)
	s.logger.Info("create event processed")
	return in, err
}

// UpdateEvent update event
func (s APIServer) UpdateEvent(ctx context.Context, in *eventapi.Event) (*eventapi.Event, error) {
	err := s.eventStorage.UpdateEvent(ctx, in)
	s.logger.Infof("update event processed for uuid: %s\n", in.Uuid)
	return in, err
}

// DeleteEvent delete event
func (s APIServer) DeleteEvent(ctx context.Context, in *eventapi.Event) (*eventapi.EventDeleteStatus, error) {
	err := s.eventStorage.DeleteEvent(ctx, in.Uuid)
	s.logger.Infof("delete event processed for uuid: %s\n", in.Uuid)
	return &eventapi.EventDeleteStatus{Deleted: err == nil}, err
}

// GetEventsForDate get events for date
func (s APIServer) GetEventsForDate(ctx context.Context, in *eventapi.EventDate) (*eventapi.EventList, error) {
	eventList, err := s.eventStorage.GetEventsForDate(ctx, time.Unix(in.GetDate().GetSeconds(), 0))
	s.logger.Info("get event list for date processed")
	return &eventapi.EventList{Events: eventList}, err
}

// GetEventsForWeek get events for wwek
func (s APIServer) GetEventsForWeek(ctx context.Context, in *eventapi.EventDate) (*eventapi.EventList, error) {
	eventList, err := s.eventStorage.GetEventsForWeek(ctx, time.Unix(in.GetDate().GetSeconds(), 0))
	s.logger.Info("get event list for week processed")
	return &eventapi.EventList{Events: eventList}, err
}

// GetEventsForMonth get events for month
func (s APIServer) GetEventsForMonth(ctx context.Context, in *eventapi.EventDate) (*eventapi.EventList, error) {
	eventList, err := s.eventStorage.GetEventsForMonth(ctx, time.Unix(in.GetDate().GetSeconds(), 0))
	s.logger.Info("get event list for month processed")
	return &eventapi.EventList{Events: eventList}, err
}

// StartServer start http server
func (s APIServer) StartServer(address string, opts ...Option) error {
	options := options{
		logger: zap.NewNop(),
	}
	for _, o := range opts {
		o.apply(&options)
	}
	s.logger = options.logger.Sugar()
	s.eventStorage = *options.storage

	listener, err := net.Listen("tcp", address)
	if err != nil {
		s.logger.Fatal("Failed to listen %v", err)
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	eventapi.RegisterApiServerServer(grpcServer, s)

	return grpcServer.Serve(listener)
}

// WithLogger apply logger
func WithLogger(log *zap.Logger) Option {
	return loggerOption{Log: log}
}

// WithRepository apply storage
func WithRepository(repository *repository.EventRepository) Option {
	return repositoryOption{EventStorage: repository}
}

func (o loggerOption) apply(opts *options) {
	opts.logger = o.Log
}

func (o repositoryOption) apply(opts *options) {
	opts.storage = o.EventStorage
}
