package repository

import (
	"context"

	"github.com/andywow/golang-lessons/lesson-calendar/pkg/calendar/model"
)

// EventRepository work with data store
type EventRepository interface {
	CreateEvent(ctx context.Context, event *model.Event) error
	GetEvents(ctx context.Context) []*model.Event
	DeleteEvent(ctx context.Context, id string) error
	UpdateEvent(ctx context.Context, event *model.Event) error
}
