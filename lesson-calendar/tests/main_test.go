package main

import (
	"os"
	"testing"

	"github.com/cucumber/godog"
)

func FeatureContext(s *godog.Suite) {
	createTest := new(createEventTest)

	s.Step(`^I create event$`, createTest.iCreateEvent)
	s.Step(`^There are no create errors$`, createTest.thereAreNoCreateErrors)
	s.Step(`^Event uuid should not be empty$`, createTest.eventUuidShouldNotBeEmpty)
	s.Step(`^I create event on busy time$`, createTest.iCreateEventOnBusyTime)
	s.Step(`^I receive date already busy error$`, createTest.iReceiveDateAlreadyBusyError)

	listTest := new(listEventTest)

	s.Step(`^I list events for date$`, listTest.iListEventsForDate)
	s.Step(`^There are no list errors$`, listTest.thereAreNoListErrors)
	s.Step(`^Event list should not be empty$`, listTest.eventListShouldNotBeEmpty)
	s.Step(`^I list events for week$`, listTest.iListEventsForWeek)
	s.Step(`^There are no list errors$`, listTest.thereAreNoListErrors)
	s.Step(`^Event list should not be empty$`, listTest.eventListShouldNotBeEmpty)
	s.Step(`^I list events for month$`, listTest.iListEventsForMonth)
	s.Step(`^There are no list errors$`, listTest.thereAreNoListErrors)
	s.Step(`^Event list should not be empty$`, listTest.eventListShouldNotBeEmpty)

}

func TestMain(m *testing.M) {
	status := godog.RunWithOptions("integration", func(s *godog.Suite) {
		FeatureContext(s)
	}, godog.Options{
		Format:    "pretty",
		Paths:     []string{"features"},
		Randomize: 0,
	})

	if st := m.Run(); st > status {
		status = st
	}

	os.Exit(status)
}
