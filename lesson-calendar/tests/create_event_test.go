package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/andywow/golang-lessons/lesson-calendar/internal/calendar/repository"
	"github.com/andywow/golang-lessons/lesson-calendar/pkg/eventapi"
)

type createEventTest struct {
	apiClient apiServerTestClient
	err       error
	event     *eventapi.Event
}

func (t *createEventTest) iCreateEvent() error {
	t.apiClient.create()
	defer t.apiClient.close()

	eventTime := time.Now()
	t.event = &eventapi.Event{
		Description: "godog test",
		Duration:    1,
		Header:      "godog test",
		StartTime:   &eventTime,
		Username:    fmt.Sprintf("godog user %d", time.Now().Unix()),
	}
	t.event, t.err = t.apiClient.client.CreateEvent(t.apiClient.callCtx, t.event)
	return nil
}

func (t *createEventTest) thereAreNoCreateErrors() error {
	return t.err
}

func (t *createEventTest) eventUuidShouldNotBeEmpty() error {
	if t.event.Uuid == "" {
		return errors.New("uuid not returned")
	}
	return nil
}

func (t *createEventTest) iCreateEventOnBusyTime() error {
	t.apiClient.create()
	defer t.apiClient.close()

	t.event, t.err = t.apiClient.client.CreateEvent(t.apiClient.callCtx, t.event)
	return nil
}

func (t *createEventTest) iReceiveDateAlreadyBusyError() error {
	if t.err != repository.ErrDateBusy {
		return t.err
	}
	return nil
}
