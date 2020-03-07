package localcache

import (
	"context"
	"testing"
	"time"

	"github.com/andywow/golang-lessons/lesson-calendar/internal/calendar/repository"
	"github.com/andywow/golang-lessons/lesson-calendar/pkg/eventapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestEvent(t *testing.T, date time.Time) eventapi.Event {
	t.Helper()
	return eventapi.Event{
		StartTime:   &date,
		Duration:    60,
		Header:      "test",
		Description: "test",
		Username:    "test",
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
	assert.Equal(t, repository.ErrDateBusy, err)

}

func TestDeleteEvent(t *testing.T) {

	s := NewEventLocalStorage()

	event := createTestEvent(t, time.Now())

	// create event
	err := s.CreateEvent(context.Background(), &event)
	require.NoError(t, err)

	// delete event
	err = s.DeleteEvent(context.Background(), event.Uuid)
	require.NoError(t, err)

	// delete non existing event
	err = s.DeleteEvent(context.Background(), event.Uuid)
	require.Error(t, err)
	assert.Equal(t, repository.ErrEventNotFound, err)

}

func TestGetEventsForDate(t *testing.T) {

	s := NewEventLocalStorage()

	event := createTestEvent(t, time.Now())
	event2 := createTestEvent(t, time.Now().Truncate(24*time.Hour).Add(time.Hour*4))
	event3 := createTestEvent(t, time.Now().Add(time.Hour*48))

	// create events
	err := s.CreateEvent(context.Background(), &event)
	require.NoError(t, err)
	err = s.CreateEvent(context.Background(), &event2)
	require.NoError(t, err)
	err = s.CreateEvent(context.Background(), &event3)
	require.NoError(t, err)

	// get events
	events, err := s.GetEventsForDate(context.Background(), time.Now())
	require.NoError(t, err)
	assert.Equal(t, 2, len(events))

}

func TestGetEventsForWeek(t *testing.T) {

	s := NewEventLocalStorage()

	startTime := time.Date(2020, time.March, 2, 12, 12, 12, 12, time.UTC)
	event := createTestEvent(t, startTime)
	event2 := createTestEvent(t, startTime.Truncate(24*time.Hour).Add(time.Hour*3*24))
	event3 := createTestEvent(t, startTime.Add(time.Hour*8*24))

	// create events
	err := s.CreateEvent(context.Background(), &event)
	require.NoError(t, err)
	err = s.CreateEvent(context.Background(), &event2)
	require.NoError(t, err)
	err = s.CreateEvent(context.Background(), &event3)
	require.NoError(t, err)

	// get events
	events, err := s.GetEventsForWeek(context.Background(), startTime)
	require.NoError(t, err)
	assert.Equal(t, 2, len(events))

}

func TestGetEventsForMonth(t *testing.T) {

	s := NewEventLocalStorage()

	startTime := time.Date(2020, time.February, 1, 12, 12, 12, 12, time.UTC)
	event := createTestEvent(t, startTime)
	event2 := createTestEvent(t, startTime.Truncate(24*time.Hour).Add(time.Hour*7*24))
	event3 := createTestEvent(t, startTime.Add(time.Hour*29*24))

	// create events
	err := s.CreateEvent(context.Background(), &event)
	require.NoError(t, err)
	err = s.CreateEvent(context.Background(), &event2)
	require.NoError(t, err)
	err = s.CreateEvent(context.Background(), &event3)
	require.NoError(t, err)

	// get events
	events, err := s.GetEventsForMonth(context.Background(), startTime)
	require.NoError(t, err)
	assert.Equal(t, 2, len(events))

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
	newTime := event.StartTime.Add(time.Duration(60*24) * time.Minute)
	event.StartTime = &newTime
	err = s.UpdateEvent(context.Background(), &event)
	require.NoError(t, err)

}
