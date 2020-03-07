package repository

import "errors"

var (
	// ErrDateBusy busy date
	ErrDateBusy = errors.New("time is already taken by another event")
	// ErrEventNotFound event not found
	ErrEventNotFound = errors.New("event not found")
	// ErrStorageUnavailable storage unavailable
	ErrStorageUnavailable = errors.New("storage unavailable")
	// ErrInvalidData invalid input data
	ErrInvalidData = errors.New("invalid input data")
	// ErrGetQueryResult no result from query
	ErrGetQueryResult = errors.New("can not get result for query")
)
