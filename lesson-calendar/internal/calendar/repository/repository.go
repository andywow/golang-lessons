package repository

import (
	"context"
	"time"

	"github.com/andywow/golang-lessons/lesson-calendar/pkg/eventapi"
)

// EventRepository work with data store
type EventRepository interface {
	CreateEvent(ctx context.Context, event *eventapi.Event) error
	GetEventsForDate(ctx context.Context, date time.Time) ([]*eventapi.Event, error)
	GetEventsForWeek(ctx context.Context, date time.Time) ([]*eventapi.Event, error)
	GetEventsForMonth(ctx context.Context, date time.Time) ([]*eventapi.Event, error)
	DeleteEvent(ctx context.Context, uuid string) error
	UpdateEvent(ctx context.Context, event *eventapi.Event) error
}
