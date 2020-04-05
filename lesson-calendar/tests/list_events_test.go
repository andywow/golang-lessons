package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/andywow/golang-lessons/lesson-calendar/pkg/eventapi"
)

type listEventTest struct {
	apiClient apiServerTestClient
	err       error
	eventList *eventapi.EventList
}

func (t *listEventTest) iListEventsForDate() error {
	t.apiClient.create()
	defer t.apiClient.close()

	eventTime := time.Now().Add(time.Hour)
	event := &eventapi.Event{
		Description: "godog test",
		Duration:    1,
		Header:      "godog test",
		StartTime:   &eventTime,
		Username:    fmt.Sprintf("godog list %d", time.Now().Unix()),
	}
	if _, err := t.apiClient.client.CreateEvent(t.apiClient.callCtx, event); err != nil {
		return err
	}

	eventDate := eventapi.EventDate{
		Date: &eventTime,
	}
	t.eventList, t.err = t.apiClient.client.GetEventsForDate(t.apiClient.callCtx, &eventDate)

	return nil
}

func (t *listEventTest) thereAreNoListErrors() error {
	return t.err
}

func (t *listEventTest) eventListShouldNotBeEmpty() error {
	if len(t.eventList.Events) == 0 {
		return errors.New("event list is empty")
	}
	return nil
}

func (t *listEventTest) iListEventsForWeek() error {
	t.apiClient.create()
	defer t.apiClient.close()

	var (
		eventTime time.Time
		err       error
	)

	if eventTime, err = time.Parse("2006.01.02", "2020.04.06"); err != nil {
		return err
	}

	event := &eventapi.Event{
		Description: "godog test",
		Duration:    1,
		Header:      "godog test",
		StartTime:   &eventTime,
		Username:    fmt.Sprintf("godog list %d", time.Now().Unix()),
	}

	if _, err := t.apiClient.client.CreateEvent(t.apiClient.callCtx, event); err != nil {
		return err
	}

	eventDate := eventapi.EventDate{
		Date: &eventTime,
	}
	t.eventList, t.err = t.apiClient.client.GetEventsForWeek(t.apiClient.callCtx, &eventDate)

	return nil
}

func (t *listEventTest) iListEventsForMonth() error {
	t.apiClient.create()
	defer t.apiClient.close()

	var (
		eventTime time.Time
		err       error
	)

	if eventTime, err = time.Parse("2006.01.02", "2020.05.01"); err != nil {
		return err
	}

	event := &eventapi.Event{
		Description: "godog test",
		Duration:    1,
		Header:      "godog test",
		StartTime:   &eventTime,
		Username:    fmt.Sprintf("godog list %d", time.Now().Unix()),
	}

	if _, err := t.apiClient.client.CreateEvent(t.apiClient.callCtx, event); err != nil {
		return err
	}

	eventDate := eventapi.EventDate{
		Date: &eventTime,
	}
	t.eventList, t.err = t.apiClient.client.GetEventsForWeek(t.apiClient.callCtx, &eventDate)

	return nil
}
