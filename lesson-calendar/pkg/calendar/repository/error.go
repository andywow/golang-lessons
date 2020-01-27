package repository

import "errors"

var (
	// ErrDateBusy busy date
	ErrDateBusy = errors.New("time is already taken by another event")
	// ErrEventNotFound event not found
	ErrEventNotFound = errors.New("event not found")
)
