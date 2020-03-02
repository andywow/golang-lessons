package repository

import (
	"context"
	"time"

	"github.com/andywow/golang-lessons/lesson-calendar/internal/calendar"
)

// EventRepository work with data store
type EventRepository interface {
	CreateEvent(ctx context.Context, event *calendar.Event) error
	GetEventsForDate(ctx context.Context, date time.Time) ([]*calendar.Event, error)
	GetEventsForWeek(ctx context.Context, date time.Time) ([]*calendar.Event, error)
	GetEventsForMonth(ctx context.Context, date time.Time) ([]*calendar.Event, error)
	DeleteEvent(ctx context.Context, uuid string) error
	UpdateEvent(ctx context.Context, event *calendar.Event) error
}
