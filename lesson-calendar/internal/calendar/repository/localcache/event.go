package localcache

import (
	"context"
	"strconv"
	"sync"

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
func NewEventLocalStorage() *EventLocalStorage {
	return &EventLocalStorage{
		events: make(map[string]*calendar.Event),
	}
}

// CreateEvent create event
func (s *EventLocalStorage) CreateEvent(ctx context.Context, event *calendar.Event) error {

	if err := event.CheckEventData(); err != nil {
		return err
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.checkIfEventBusy(event) {
		return repository.ErrDateBusy
	}

	s.currentUUID++
	event.UUID = strconv.Itoa(s.currentUUID)

	// add new structure or user can modify event in storage explicity, not through interface
	s.events[event.UUID] = &calendar.Event{
		UUID:        event.UUID,
		StartTime:   event.StartTime,
		EndTime:     event.EndTime,
		Header:      event.Header,
		Description: event.Description,
		User:        event.User,
	}

	return nil
}

// GetEvents get events
func (s *EventLocalStorage) GetEvents(ctx context.Context) []calendar.Event {
	events := make([]calendar.Event, len(s.events))

	s.mutex.Lock()
	defer s.mutex.Unlock()

	index := 0
	for _, event := range s.events {
		events[index] = *event
		index++
	}

	return events
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

	if err := event.CheckEventData(); err != nil {
		return err
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, ok := s.events[event.UUID]; !ok {
		return repository.ErrEventNotFound
	}

	if s.checkIfEventBusy(event) {
		return repository.ErrDateBusy
	}

	s.events[event.UUID].StartTime = event.StartTime
	s.events[event.UUID].EndTime = event.EndTime
	s.events[event.UUID].Header = event.Header
	s.events[event.UUID].Description = event.Description
	s.events[event.UUID].User = event.User

	return nil
}

func (s *EventLocalStorage) checkIfEventBusy(event *calendar.Event) bool {
	for _, storageEvent := range s.events {
		if event.StartTime.Equal(storageEvent.StartTime) || event.EndTime.Equal(storageEvent.EndTime) ||
			(event.StartTime.After(storageEvent.StartTime) && event.StartTime.Before(storageEvent.EndTime)) ||
			(event.EndTime.After(storageEvent.StartTime) && event.EndTime.Before(storageEvent.EndTime)) {
			return true
		}
	}
	return false
}
