package repository

import (
	"context"

	"github.com/andywow/golang-lessons/lesson-calendar/internal/calendar"
)

// EventRepository work with data store
type EventRepository interface {
	CreateEvent(ctx context.Context, event *calendar.Event) error
	GetEvents(ctx context.Context) []calendar.Event
	DeleteEvent(ctx context.Context, uuid string) error
	UpdateEvent(ctx context.Context, event *calendar.Event) error
}
