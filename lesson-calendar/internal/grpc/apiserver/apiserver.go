package apiserver

import (
	"context"
	"errors"
	"net"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
	s.logger.Info("create event processing")
	if err := s.eventStorage.CheckIfTimeIsBusy(ctx, in); err != nil {
		return nil, s.getStatusFromError(err)
	}
	if err := s.eventStorage.CreateEvent(ctx, in); err != nil {
		return nil, s.getStatusFromError(err)
	}
	s.logger.Info("create event processed")
	return in, nil
}

// UpdateEvent update event
func (s APIServer) UpdateEvent(ctx context.Context, in *eventapi.Event) (*eventapi.Event, error) {
	s.logger.Infof("update event received for uuid: %s", in.Uuid)
	if in.StartTime != nil {
		if err := s.eventStorage.CheckIfTimeIsBusy(ctx, in); err != nil {
			return nil, s.getStatusFromError(err)
		}
	}
	if err := s.eventStorage.UpdateEvent(ctx, in); err != nil {
		return nil, s.getStatusFromError(err)
	}
	s.logger.Infof("update event processed for uuid: %s", in.Uuid)
	return in, nil
}

// DeleteEvent delete event
func (s APIServer) DeleteEvent(ctx context.Context, in *eventapi.EventDelete) (*eventapi.EventDeleteStatus, error) {
	s.logger.Infof("delete event received for uuid: %s", in.Uuid)
	if err := s.eventStorage.DeleteEvent(ctx, in.Uuid); err != nil {
		return nil, s.getStatusFromError(err)
	}
	s.logger.Infof("delete event processed for uuid: %s", in.Uuid)
	return &eventapi.EventDeleteStatus{}, nil
}

// GetEventsForDate get events for date
func (s APIServer) GetEventsForDate(ctx context.Context, in *eventapi.EventDate) (*eventapi.EventList, error) {
	s.logger.Info("get event list for date received")
	eventList, err := s.eventStorage.GetEventsForDate(ctx, *in.GetDate())
	if err != nil {
		return nil, s.getStatusFromError(err)
	}
	s.logger.Info("get event list for date processed")
	return &eventapi.EventList{Events: eventList}, err
}

// GetEventsForWeek get events for wwek
func (s APIServer) GetEventsForWeek(ctx context.Context, in *eventapi.EventDate) (*eventapi.EventList, error) {
	s.logger.Info("get event list for week received")
	eventList, err := s.eventStorage.GetEventsForWeek(ctx, *in.GetDate())
	if err != nil {
		return nil, s.getStatusFromError(err)
	}
	s.logger.Info("get event list for week date processed")
	return &eventapi.EventList{Events: eventList}, err
}

// GetEventsForMonth get events for month
func (s APIServer) GetEventsForMonth(ctx context.Context, in *eventapi.EventDate) (*eventapi.EventList, error) {
	s.logger.Info("get event list for month received")
	eventList, err := s.eventStorage.GetEventsForMonth(ctx, *in.GetDate())
	if err != nil {
		return nil, s.getStatusFromError(err)
	}
	s.logger.Info("get event list for month processed")
	return &eventapi.EventList{Events: eventList}, err
}

// StartServer start http server
func (s APIServer) StartServer(ctx context.Context, address string, opts ...Option) error {
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

	go func() {
		select {
		case <-ctx.Done():
			grpcServer.GracefulStop()
		}
	}()

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

func (s APIServer) getStatusFromError(err error) error {
	if err != nil {
		s.logger.Error(err)
		if errors.Is(err, repository.ErrDateBusy) {
			return status.Error(codes.AlreadyExists, repository.ErrDateBusy.Error())
		}
		if errors.Is(err, repository.ErrInvalidData) {
			return status.Error(codes.InvalidArgument, repository.ErrInvalidData.Error())
		}
		if errors.Is(err, repository.ErrStorageUnavailable) {
			return status.Error(codes.Unavailable, repository.ErrStorageUnavailable.Error())
		}
		return status.Error(codes.Internal, "internal error")
	}
	return nil
}
