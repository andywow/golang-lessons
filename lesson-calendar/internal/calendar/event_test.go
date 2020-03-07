package calendar

import (
	"testing"
	"time"

	"github.com/andywow/golang-lessons/lesson-calendar/pkg/eventapi"
	"github.com/stretchr/testify/require"
)

func TestCheckEventData(t *testing.T) {
	event := &eventapi.Event{}
	err := CheckEventData(event)
	require.Error(t, err)

	eventTime := time.Now()
	event = &eventapi.Event{
		StartTime:   &eventTime,
		Duration:    60,
		Header:      "test",
		Description: "test",
		Username:    "test",
	}
	err = CheckEventData(event)
	require.NoError(t, err)
}
