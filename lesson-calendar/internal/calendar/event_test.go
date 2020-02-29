package calendar

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCheckEventData(t *testing.T) {
	event := Event{}
	err := event.CheckEventData()
	require.Error(t, err)
	event = Event{
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(time.Hour),
		Header:      "Test event",
		Description: "Test event description",
		User:        "Bob",
	}
	err = event.CheckEventData()
	require.NoError(t, err)
}
