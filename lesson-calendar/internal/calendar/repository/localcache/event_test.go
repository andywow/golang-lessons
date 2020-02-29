package localcache

import (
	"context"
	"testing"
	"time"

	"github.com/andywow/golang-lessons/lesson-calendar/internal/calendar"
	"github.com/andywow/golang-lessons/lesson-calendar/internal/calendar/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestEvent(t *testing.T, date time.Time) calendar.Event {
	t.Helper()
	return calendar.Event{
		StartTime:   date,
		EndTime:     date.Add(time.Hour),
		Header:      "test",
		Description: "test",
		User:        "test",
	}
}

func TestCreateEvent(t *testing.T) {

	s := NewEventLocalStorage()

	now := time.Now()

	event := createTestEvent(t, now)
	err := s.CreateEvent(context.Background(), &event)
	require.NoError(t, err)

	// check event
	event2 := createTestEvent(t, now.Add(time.Hour*24))
	err = s.CreateEvent(context.Background(), &event2)
	require.NoError(t, err)

	// check for event on same date
	eventSameTime := createTestEvent(t, now)
	err = s.CreateEvent(context.Background(), &eventSameTime)
	require.Error(t, err)
	assert.Equal(t, err, repository.ErrDateBusy)

}

func TestDeleteEvent(t *testing.T) {

	s := NewEventLocalStorage()

	event := createTestEvent(t, time.Now())

	// create event
	err := s.CreateEvent(context.Background(), &event)
	require.NoError(t, err)

	// delete event
	err = s.DeleteEvent(context.Background(), event.UUID)
	require.NoError(t, err)

	// delete non existing event
	err = s.DeleteEvent(context.Background(), event.UUID)
	require.Error(t, err)
	assert.Equal(t, err, repository.ErrEventNotFound)

}

func TestGetEvents(t *testing.T) {

	s := NewEventLocalStorage()

	event := createTestEvent(t, time.Now())
	event2 := createTestEvent(t, time.Now().Add(time.Hour*4))

	// create events
	err := s.CreateEvent(context.Background(), &event)
	require.NoError(t, err)
	err = s.CreateEvent(context.Background(), &event2)
	require.NoError(t, err)

	// get events
	events := s.GetEvents(context.Background())
	require.NoError(t, err)
	assert.Equal(t, len(events), 2)

}

func TestUpdateEvent(t *testing.T) {

	s := NewEventLocalStorage()

	event := createTestEvent(t, time.Now())

	// update non existing event
	err := s.UpdateEvent(context.Background(), &event)
	require.Error(t, err)
	assert.Equal(t, err, repository.ErrEventNotFound)

	// create event
	err = s.CreateEvent(context.Background(), &event)
	require.NoError(t, err)

	// update event
	event.StartTime = time.Now().Add(time.Hour * time.Duration(48))
	event.EndTime = time.Now().Add(time.Hour * time.Duration(50))
	err = s.UpdateEvent(context.Background(), &event)
	require.NoError(t, err)

}
