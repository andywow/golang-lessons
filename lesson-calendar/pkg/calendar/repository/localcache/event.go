package localcache

import (
	"context"
	"fmt"
	"sync"

	"github.com/andywow/golang-lessons/lesson-calendar/pkg/calendar/model"
	"github.com/andywow/golang-lessons/lesson-calendar/pkg/calendar/repository"
)

// EventLocalStorage local memory storage
type EventLocalStorage struct {
	events    map[string]*model.Event
	mutex     *sync.Mutex
	currentId int
}

// NewEventLocalStorage constructor
func NewEventLocalStorage() *EventLocalStorage {
	return &EventLocalStorage{
		events: make(map[string]*model.Event),
		mutex:  new(sync.Mutex),
	}
}

// CreateEvent create event
func (s *EventLocalStorage) CreateEvent(ctx context.Context, event *model.Event) error {

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.checkIfEventBusy(event) {
		return repository.ErrDateBusy
	}

	event.ID = string(s.currentId + 1)
	// add new structure or user can modify event in storage explicity, not through interface
	s.events[event.ID] = &model.Event{
		ID:   event.ID,
		Time: event.Time,
	}

	return nil
}

// GetEvents get events
func (s *EventLocalStorage) GetEvents(ctx context.Context) []*model.Event {
	events := make([]*model.Event, len(s.events))

	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, event := range s.events {
		events = append(events, event)
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
func (s *EventLocalStorage) UpdateEvent(ctx context.Context, event *model.Event) error {

	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, ok := s.events[event.ID]
	if !ok {
		return repository.ErrEventNotFound
	}

	if s.checkIfEventBusy(event) {
		return repository.ErrDateBusy
	}

	s.events[event.ID].Time = event.Time

	return nil
}

func (s *EventLocalStorage) checkIfEventBusy(event *model.Event) bool {
	for _, storageEvent := range s.events {
		if event.Time.Year() == storageEvent.Time.Year() &&
			event.Time.Month() == storageEvent.Time.Month() &&
			event.Time.Day() == storageEvent.Time.Day() {
			fmt.Println("ev: ", event)
			fmt.Println("sev: ", storageEvent)
			return true
		}
	}
	return false
}
