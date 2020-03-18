package localcache

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/andywow/golang-lessons/lesson-calendar/internal/calendar"
	"github.com/andywow/golang-lessons/lesson-calendar/internal/calendar/repository"
	"github.com/andywow/golang-lessons/lesson-calendar/pkg/eventapi"
)

// EventLocalStorage local memory storage
type EventLocalStorage struct {
	events      map[string]*eventapi.Event
	mutex       sync.Mutex
	currentUUID int
}

// NewEventLocalStorage constructor
func NewEventLocalStorage() repository.EventRepository {
	return &EventLocalStorage{
		events: make(map[string]*eventapi.Event),
	}
}

// Close fake close
func (s *EventLocalStorage) Close() error {
	return nil
}

// CreateEvent create event
func (s *EventLocalStorage) CreateEvent(ctx context.Context, event *eventapi.Event) error {

	if err := calendar.CheckEventData(event); err != nil {
		return err
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if err := s.CheckIfTimeIsBusy(ctx, event); err != nil {
		return err
	}

	s.currentUUID++
	event.Uuid = strconv.Itoa(s.currentUUID)

	// add new structure or user can modify event in storage explicity, not through interface
	s.events[event.Uuid] = &eventapi.Event{
		Uuid:        event.Uuid,
		StartTime:   event.StartTime,
		Duration:    event.Duration,
		Header:      event.Header,
		Description: event.Description,
		Username:    event.Username,
	}

	return nil
}

// GetEventsForDate get events for 1 day
func (s *EventLocalStorage) GetEventsForDate(ctx context.Context, date time.Time) ([]*eventapi.Event, error) {
	startTime := date.Truncate(24 * time.Hour)
	endTime := date.Truncate(24 * time.Hour).Add(24 * time.Hour)
	return s.getEvents(startTime, endTime), nil
}

// GetEventsForWeek get events for week
func (s *EventLocalStorage) GetEventsForWeek(ctx context.Context, date time.Time) ([]*eventapi.Event, error) {
	startTime := date.Truncate(24 * time.Hour)
	endTime := date.Truncate(24 * time.Hour).Add(24 * 7 * time.Hour)
	return s.getEvents(startTime, endTime), nil
}

// GetEventsForMonth get events for month
func (s *EventLocalStorage) GetEventsForMonth(ctx context.Context, date time.Time) ([]*eventapi.Event, error) {
	startTime := date.Truncate(24 * time.Hour)
	endTime := date.Truncate(24*time.Hour).AddDate(0, 1, 0)
	return s.getEvents(startTime, endTime), nil
}

// GetEventsForNotification get events for send notifications
func (s *EventLocalStorage) GetEventsForNotification(ctx context.Context, date time.Time) ([]*eventapi.Event, error) {
	return nil, errors.New("unimplemented")
}

// DeleteEvent delete event
func (s *EventLocalStorage) DeleteEvent(ctx context.Context, uuid string) error {

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, ok := s.events[uuid]; ok {
		delete(s.events, uuid)
		return nil
	}

	return repository.ErrEventNotFound
}

// UpdateEvent update event
func (s *EventLocalStorage) UpdateEvent(ctx context.Context, event *eventapi.Event) error {

	if err := calendar.CheckEventData(event); err != nil {
		return err
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, ok := s.events[event.Uuid]; !ok {
		return repository.ErrEventNotFound
	}

	if err := s.CheckIfTimeIsBusy(ctx, event); err != nil {
		return err
	}

	s.events[event.Uuid].StartTime = event.StartTime
	s.events[event.Uuid].Duration = event.Duration
	s.events[event.Uuid].Header = event.Header
	s.events[event.Uuid].Description = event.Description
	s.events[event.Uuid].Username = event.Username

	return nil
}

// CheckIfTimeIsBusy check if time is busy
func (s *EventLocalStorage) CheckIfTimeIsBusy(ctx context.Context, event *eventapi.Event) error {
	for _, storageEvent := range s.events {
		if event.Uuid != storageEvent.Uuid &&
			((event.StartTime.Unix() >= storageEvent.StartTime.Unix() &&
				event.StartTime.Unix() <= storageEvent.StartTime.Add(time.Duration(storageEvent.Duration)*time.Minute).Unix()) ||
				(event.StartTime.Add(time.Duration(event.Duration)*time.Minute).Unix() >= storageEvent.StartTime.Unix() &&
					event.StartTime.Add(time.Duration(event.Duration)*time.Minute).Unix() <= storageEvent.StartTime.Add(time.Duration(storageEvent.Duration)*time.Minute).Unix())) {
			return repository.ErrDateBusy
		}
	}
	return nil
}

func (s *EventLocalStorage) getEvents(startDate, endDate time.Time) []*eventapi.Event {
	var events []*eventapi.Event

	startTime := startDate.Unix()
	endTime := endDate.Unix()

	s.mutex.Lock()
	defer s.mutex.Unlock()

	index := 0
	for _, event := range s.events {
		if startTime <= event.StartTime.Unix() && endTime >= event.StartTime.Unix() {
			events = append(events, event)
			index++
		}
	}

	return events
}
