package calendar

import (
	"testing"

	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/require"
)

func TestCheckEventData(t *testing.T) {
	event := &Event{}
	err := CheckEventData(event)
	require.Error(t, err)

	event = &Event{
		StartTime:   ptypes.TimestampNow(),
		Duration:    60 * 60,
		Header:      "test",
		Description: "test",
		User:        "test",
	}
	err = CheckEventData(event)
	require.NoError(t, err)
}
