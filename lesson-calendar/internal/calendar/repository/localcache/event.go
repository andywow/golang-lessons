package localcache

import (
	"context"
	"sync"

	"github.com/andywow/golang-lessons/lesson-calendar/internal/calendar"
	"github.com/andywow/golang-lessons/lesson-calendar/internal/calendar/repository"
)

// EventLocalStorage local memory storage
type EventLocalStorage struct {
	events    map[string]*calendar.Event
	mutex     sync.Mutex
	currentID int
}

// NewEventLocalStorage constructor
func NewEventLocalStorage() *EventLocalStorage {
	return &EventLocalStorage{
		events: make(map[string]*calendar.Event),
	}
}

// CreateEvent create event
func (s *EventLocalStorage) CreateEvent(ctx context.Context, event *calendar.Event) error {

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.checkIfEventBusy(event) {
		return repository.ErrDateBusy
	}

	event.ID = string(s.currentID + 1)
	// add new structure or user can modify event in storage explicity, not through interface
	s.events[event.ID] = &calendar.Event{
		ID:   event.ID,
		Time: event.Time,
	}

	return nil
}

// GetEvents get events
func (s *EventLocalStorage) GetEvents(ctx context.Context) []calendar.Event {
	events := make([]calendar.Event, len(s.events))

	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, event := range s.events {
		events = append(events, *event)
	}

	return events
}

// DeleteEvent delete event
func (s *EventLocalStorage) DeleteEvent(ctx context.Context, id string) error {

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, ok := s.events[id]; ok {
		delete(s.events, id)
		return nil
	}

	return repository.ErrEventNotFound
}

// UpdateEvent update event
func (s *EventLocalStorage) UpdateEvent(ctx context.Context, event *calendar.Event) error {

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, ok := s.events[event.ID]; !ok {
		return repository.ErrEventNotFound
	}

	if s.checkIfEventBusy(event) {
		return repository.ErrDateBusy
	}

	s.events[event.ID].Time = event.Time

	return nil
}

func (s *EventLocalStorage) checkIfEventBusy(event *calendar.Event) bool {
	for _, storageEvent := range s.events {
		if event.Time.Year() == storageEvent.Time.Year() &&
			event.Time.Month() == storageEvent.Time.Month() &&
			event.Time.Day() == storageEvent.Time.Day() {
			return true
		}
	}
	return false
}
