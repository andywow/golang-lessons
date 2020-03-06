package calendar

import (
	"testing"

	"github.com/andywow/golang-lessons/lesson-calendar/pkg/eventapi"
	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/require"
)

func TestCheckEventData(t *testing.T) {
	event := &eventapi.Event{}
	err := CheckEventData(event)
	require.Error(t, err)

	event = &eventapi.Event{
		StartTime:   ptypes.TimestampNow(),
		Duration:    60 * 60,
		Header:      "test",
		Description: "test",
		User:        "test",
	}
	err = CheckEventData(event)
	require.NoError(t, err)
}
