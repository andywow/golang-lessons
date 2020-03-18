package calendar

import (
	"encoding/json"

	"github.com/andywow/golang-lessons/lesson-calendar/pkg/eventapi"
	"github.com/pkg/errors"
)

// CheckEventData check event data
func CheckEventData(e *eventapi.Event) error {
	if e.StartTime == nil || e.Duration <= 0 || e.Header == "" || e.Description == "" || e.Username == "" {
		return ErrIncorrectEvent
	}
	return nil
}

// ConvertToJSON convert to json
func ConvertToJSON(e *eventapi.Event) ([]byte, error) {
	blob, err := json.Marshal(&e)
	if err != nil {
		return nil, errors.Wrap(err, "failed convert to json")
	}
	return blob, nil
}

// ConvertFromJSON convert from json
func ConvertFromJSON(blob []byte) (*eventapi.Event, error) {
	var e eventapi.Event
	err := json.Unmarshal(blob, &e)
	if err != nil {
		return nil, errors.Wrap(err, "failed convert from json")
	}
	return &e, nil
}
