package calendar

import "github.com/andywow/golang-lessons/lesson-calendar/pkg/eventapi"

// CheckEventData check event data
func CheckEventData(e *eventapi.Event) error {
	if e.StartTime == nil || e.Duration <= 0 || e.Header == "" || e.Description == "" || e.Username == "" {
		return ErrIncorrectEvent
	}
	return nil
}
