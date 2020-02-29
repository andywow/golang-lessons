package calendar

import (
	"time"
)

// Event of calendar
type Event struct {
	UUID               string
	StartTime, EndTime time.Time
	Header             string
	Description        string
	User               string
}

// CheckEventData check event data
func (e Event) CheckEventData() error {
	if e.StartTime.After(e.EndTime) || e.Header == "" || e.Description == "" ||
		e.User == "" {
		return ErrIncorrectEvent
	}
	return nil
}
