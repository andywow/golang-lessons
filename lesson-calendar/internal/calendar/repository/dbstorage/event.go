package dbstorage

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/andywow/golang-lessons/lesson-calendar/internal/calendar/repository"
	"github.com/andywow/golang-lessons/lesson-calendar/pkg/eventapi"
	"github.com/pkg/errors"

	// init sql driver
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
)

// EventDatabase event database
type EventDatabase struct {
	Database *sqlx.DB
}

// NewDatabaseStorage constructor
func NewDatabaseStorage(ctx context.Context, dbHost string, dbPort int,
	dbName, dbUser, dbPassword string) (repository.EventRepository, error) {
	db, err := sqlx.Connect("pgx", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName))
	if err != nil {
		return nil, errors.Wrap(err, "could not initialize db")
	}
	_, err = db.Query("SELECT 1")
	if err != nil {
		return nil, errors.Wrap(err, "could not execute test query")
	}
	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(10)
	// map json fields from protobuf annotations
	db.Mapper = reflectx.NewMapperFunc("json", func(str string) string {
		return str
	})
	eventDatabase := EventDatabase{
		Database: db,
	}

	go func() {
		select {
		case <-ctx.Done():
			db.Close()
		}
	}()

	return &eventDatabase, nil
}

// CheckIfTimeIsBusy check if event time is busy
func (s *EventDatabase) CheckIfTimeIsBusy(ctx context.Context, event *eventapi.Event) error {
	if err := s.verifyConnection(ctx); err != nil {
		return err
	}
	var err error
	intUUID := 0
	if event.Uuid != "" {
		intUUID, err = strconv.Atoi(event.Uuid)
		if err != nil {
			return errors.Wrap(repository.ErrInvalidData, "invalid uuid data")
		}
	}
	rs, err := s.Database.NamedQueryContext(ctx,
		`select count(*) as count from calendar.event where
			uuid!=:uuid and username=:username and deleted=false 
			and (
				(:start_time>=start_time and :start_time<start_time + duration * interval '1 minute') 
				or 
				(:end_time>start_time and :end_time<=start_time + duration * interval '1 minute')
			)`,
		map[string]interface{}{
			"uuid":       intUUID,
			"username":   event.Username,
			"start_time": event.StartTime,
			"end_time":   event.StartTime.Add(time.Duration(event.Duration) * time.Minute),
		})
	if err != nil {
		return errors.Wrap(err, "could not execute check time statement")
	}
	if rs.Next() {
		var count int64
		err := rs.Scan(&count)
		if err != nil {
			return errors.Wrap(err, "could not parse select result")
		}
		if count != 0 {
			return repository.ErrDateBusy
		}
		return nil
	}
	return errors.New("could not get select result")
}

// Close connection
func (s *EventDatabase) Close() error {
	if err := s.Database.Close(); err != nil {
		return errors.Wrap(err, "fail to close connection")
	}
	return nil
}

// CreateEvent create event
func (s *EventDatabase) CreateEvent(ctx context.Context, event *eventapi.Event) error {
	if err := s.verifyConnection(ctx); err != nil {
		return err
	}
	rs, err := s.Database.NamedQueryContext(ctx,
		`insert into calendar.event (start_time, duration, header, description, username, notification_period)
			values(:start_time, :duration, :header, :description, :username, :notification_period) 
			returning uuid`,
		event)
	if err != nil {
		return errors.Wrap(err, "could not create new event")
	}
	if rs.Next() {
		err := rs.Scan(&event.Uuid)
		if err != nil {
			return errors.Wrap(err, "could not get id for new event")
		}
		return nil
	}
	return repository.ErrGetQueryResult
}

// GetEventsForDate get events for 1 day
func (s *EventDatabase) GetEventsForDate(ctx context.Context, date time.Time) ([]*eventapi.Event, error) {
	startTime := date.Truncate(24 * time.Hour)
	endTime := date.Truncate(24 * time.Hour).Add(24 * time.Hour)
	return s.getEvents(ctx, startTime, endTime)
}

// GetEventsForWeek get events for week
func (s *EventDatabase) GetEventsForWeek(ctx context.Context, date time.Time) ([]*eventapi.Event, error) {
	startTime := date.Truncate(24 * time.Hour)
	endTime := date.Truncate(24 * time.Hour).Add(24 * 7 * time.Hour)
	return s.getEvents(ctx, startTime, endTime)
}

// GetEventsForMonth get events for month
func (s *EventDatabase) GetEventsForMonth(ctx context.Context, date time.Time) ([]*eventapi.Event, error) {
	startTime := date.Truncate(24 * time.Hour)
	endTime := date.Truncate(24*time.Hour).AddDate(0, 1, 0)
	return s.getEvents(ctx, startTime, endTime)
}

// GetEventsForNotification get events for send notifications
func (s *EventDatabase) GetEventsForNotification(ctx context.Context, date time.Time) ([]*eventapi.Event, error) {
	if err := s.verifyConnection(ctx); err != nil {
		return nil, err
	}
	rs, err := s.Database.NamedQueryContext(ctx,
		`select uuid, start_time, duration, header, description, username, notification_period from calendar.event where 
			deleted=false and notification_period!=0 and 
			start_time - :current_time = + notification_period * interval '1 minute'`,
		map[string]interface{}{
			"current_time": date,
		})
	if err != nil {
		return nil, errors.Wrap(err, "could not execute get event for notification statement")
	}
	events := []*eventapi.Event{}
	for rs.Next() {
		var event eventapi.Event
		if err := rs.StructScan(&event); err != nil {
			return nil, errors.Wrap(err, "could not parse select result")
		}
		events = append(events, &event)
	}
	return events, nil
}

// DeleteEvent delete event
func (s *EventDatabase) DeleteEvent(ctx context.Context, uuid string) error {
	if err := s.verifyConnection(ctx); err != nil {
		return err
	}
	rs, err := s.Database.ExecContext(ctx,
		`update calendar.event set deleted=true where uuid=$1`,
		uuid)
	if err != nil {
		return errors.Wrap(err, "could not execute delete query")
	}
	rowCount, err := rs.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "could not get deleted rows")
	}
	if rowCount == 0 {
		return repository.ErrEventNotFound
	}
	return nil
}

// UpdateEvent update event
func (s *EventDatabase) UpdateEvent(ctx context.Context, event *eventapi.Event) error {
	if err := s.verifyConnection(ctx); err != nil {
		return err
	}
	query := "update calendar.event set "
	fields := []string{}
	if event.StartTime != nil {
		fields = append(fields, "start_time=:start_time")
	}
	if event.Duration != 0 {
		fields = append(fields, "duration=:duration")
	}
	if event.Header != "" {
		fields = append(fields, "header=:header")
	}
	if event.Description != "" {
		fields = append(fields, "description=:description")
	}
	if event.NotificationPeriod != 0 {
		fields = append(fields, "notification_period=:notification_period")
	}
	if len(fields) == 0 {
		return errors.Wrap(repository.ErrInvalidData, "no fields to update")
	}
	query = query + strings.Join(fields, ",")
	query = query + " where uuid=:uuid returning uuid"
	rs, err := s.Database.NamedQueryContext(ctx, query, event)
	if err != nil {
		return errors.Wrap(err, "could not execute update query")
	}
	if rs.Next() {
		if err := rs.Scan(&event.Uuid); err != nil {
			return errors.Wrap(err, "could not get id for updated event")
		}
		return nil
	}
	return repository.ErrEventNotFound
}

// check connection
func (s *EventDatabase) verifyConnection(ctx context.Context) error {
	if err := s.Database.PingContext(ctx); err != nil {
		return errors.Wrap(repository.ErrStorageUnavailable, err.Error())
	}
	return nil
}

func (s *EventDatabase) getEvents(ctx context.Context, startDate, endDate time.Time) ([]*eventapi.Event, error) {
	if err := s.verifyConnection(ctx); err != nil {
		return nil, err
	}
	rs, err := s.Database.NamedQueryContext(ctx,
		`select uuid, start_time, duration, header, description, username, notification_period from calendar.event where 
		  deleted=false and 
			start_time>=:start_time and 
			start_time<=:end_time`,
		map[string]interface{}{
			"start_time": startDate,
			"end_time":   endDate,
		})
	if err != nil {
		return nil, errors.Wrap(err, "could not execute get event statement")
	}
	events := []*eventapi.Event{}
	for rs.Next() {
		var event eventapi.Event
		if err := rs.StructScan(&event); err != nil {
			return nil, errors.Wrap(err, "could not parse select result")
		}
		events = append(events, &event)
	}
	return events, nil
}
