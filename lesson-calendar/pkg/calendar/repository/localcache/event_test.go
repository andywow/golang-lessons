package localcache

import (
	"context"
	"testing"
	"time"

	"github.com/andywow/golang-lessons/lesson-calendar/pkg/calendar/model"
	"github.com/andywow/golang-lessons/lesson-calendar/pkg/calendar/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateEvent(t *testing.T) {

	s := NewEventLocalStorage()

	event := model.Event{
		Time: time.Now(),
	}

	err := s.CreateEvent(context.Background(), &event)
	require.NoError(t, err)

	eventSameTime := model.Event{
		Time: time.Now(),
	}

	// check for event on same date
	err = s.CreateEvent(context.Background(), &eventSameTime)
	require.Error(t, err)
	assert.Equal(t, err, repository.ErrDateBusy)

}

func TestDeleteEvent(t *testing.T) {

	s := NewEventLocalStorage()

	currentTime := time.Now()
	event := model.Event{
		Time: currentTime,
	}

	// create event
	err := s.CreateEvent(context.Background(), &event)
	require.NoError(t, err)

	// delete event
	err = s.DeleteEvent(context.Background(), event.ID)
	require.NoError(t, err)

	// delete non existing event
	err = s.DeleteEvent(context.Background(), event.ID)
	require.Error(t, err)
	assert.Equal(t, err, repository.ErrEventNotFound)

}

func TestGetEvents(t *testing.T) {

	s := NewEventLocalStorage()

	event := model.Event{
		Time: time.Now(),
	}
	event2 := model.Event{
		Time: time.Now().Add(time.Hour * time.Duration(48)),
	}

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

	event := model.Event{
		Time: time.Now(),
	}

	// update non existing event
	err := s.UpdateEvent(context.Background(), &event)
	require.Error(t, err)
	assert.Equal(t, err, repository.ErrEventNotFound)

	// create event
	err = s.CreateEvent(context.Background(), &event)
	require.NoError(t, err)

	// update event
	event.Time = time.Now().Add(time.Hour * time.Duration(48))
	err = s.UpdateEvent(context.Background(), &event)
	require.NoError(t, err)

}
