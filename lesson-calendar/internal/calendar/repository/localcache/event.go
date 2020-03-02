package localcache

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/andywow/golang-lessons/lesson-calendar/internal/calendar"
	"github.com/andywow/golang-lessons/lesson-calendar/internal/calendar/repository"
)

// EventLocalStorage local memory storage
type EventLocalStorage struct {
	events      map[string]*calendar.Event
	mutex       sync.Mutex
	currentUUID int
}

// NewEventLocalStorage constructor
func NewEventLocalStorage() repository.EventRepository {
	return &EventLocalStorage{
		events: make(map[string]*calendar.Event),
	}
}

// CreateEvent create event
func (s *EventLocalStorage) CreateEvent(ctx context.Context, event *calendar.Event) error {

	if err := calendar.CheckEventData(event); err != nil {
		return err
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.checkIfEventBusy(event) {
		return repository.ErrDateBusy
	}

	s.currentUUID++
	event.Uuid = strconv.Itoa(s.currentUUID)

	// add new structure or user can modify event in storage explicity, not through interface
	s.events[event.Uuid] = &calendar.Event{
		Uuid:        event.Uuid,
		StartTime:   event.StartTime,
		Duration:    event.Duration,
		Header:      event.Header,
		Description: event.Description,
		User:        event.User,
	}

	return nil
}

// GetEventsForDate get events for 1 day
func (s *EventLocalStorage) GetEventsForDate(ctx context.Context, date time.Time) ([]*calendar.Event, error) {
	startTime := date.Truncate(24 * time.Hour)
	endTime := date.Truncate(24 * time.Hour).Add(24 * time.Hour)
	return s.getEvents(startTime, endTime), nil
}

// GetEventsForWeek get events for week
func (s *EventLocalStorage) GetEventsForWeek(ctx context.Context, date time.Time) ([]*calendar.Event, error) {
	startTime := date.Truncate(24 * time.Hour)
	endTime := date.Truncate(24 * time.Hour).Add(24 * 7 * time.Hour)
	return s.getEvents(startTime, endTime), nil
}

// GetEventsForMonth get events for month
func (s *EventLocalStorage) GetEventsForMonth(ctx context.Context, date time.Time) ([]*calendar.Event, error) {
	startTime := date.Truncate(24 * time.Hour)
	endTime := date.Truncate(24*time.Hour).AddDate(0, 1, 0)
	return s.getEvents(startTime, endTime), nil
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
func (s *EventLocalStorage) UpdateEvent(ctx context.Context, event *calendar.Event) error {

	if err := calendar.CheckEventData(event); err != nil {
		return err
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, ok := s.events[event.Uuid]; !ok {
		return repository.ErrEventNotFound
	}

	if s.checkIfEventBusy(event) {
		return repository.ErrDateBusy
	}

	s.events[event.Uuid].StartTime = event.StartTime
	s.events[event.Uuid].Duration = event.Duration
	s.events[event.Uuid].Header = event.Header
	s.events[event.Uuid].Description = event.Description
	s.events[event.Uuid].User = event.User

	return nil
}

func (s *EventLocalStorage) checkIfEventBusy(event *calendar.Event) bool {
	for _, storageEvent := range s.events {
		if event.Uuid != storageEvent.Uuid &&
			((event.StartTime.GetSeconds() >= storageEvent.StartTime.GetSeconds() &&
				event.StartTime.GetSeconds() <= storageEvent.StartTime.GetSeconds()+storageEvent.Duration) ||
				(event.StartTime.GetSeconds()+event.Duration >= storageEvent.StartTime.GetSeconds() &&
					event.StartTime.GetSeconds()+event.Duration <= storageEvent.StartTime.GetSeconds()+storageEvent.Duration)) {
			return true
		}
	}
	return false
}

func (s *EventLocalStorage) getEvents(startDate, endDate time.Time) []*calendar.Event {
	var events []*calendar.Event

	startTime := startDate.Unix()
	endTime := endDate.Unix()

	s.mutex.Lock()
	defer s.mutex.Unlock()

	index := 0
	for _, event := range s.events {
		if startTime <= event.StartTime.GetSeconds() && endTime >= event.StartTime.GetSeconds() {
			events = append(events, event)
			index++
		}
	}

	return events
}
