package calendar

// CheckEventData check event data
func CheckEventData(e *Event) error {
	if e.StartTime == nil || e.Duration <= 0 || e.Header == "" || e.Description == "" || e.User == "" {
		return ErrIncorrectEvent
	}
	return nil
}
