package rabbitmq

import (
	"context"
	"fmt"
	"time"

	"github.com/andywow/golang-lessons/lesson-calendar/internal/calendar/config"
	"github.com/andywow/golang-lessons/lesson-calendar/internal/calendar/msgsystem"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

const msgContentType = "application/json"

// RabbitMQ rabbit message system
type RabbitMQ struct {
	conn  *amqp.Connection
	ch    *amqp.Channel
	queue *amqp.Queue
}

// NewRabbitMQ return new instance
func NewRabbitMQ(ctx context.Context, cfg config.RabbitMQConfig) (msgsystem.MsgSystem, error) {
	var (
		conn  *amqp.Connection
		ch    *amqp.Channel
		queue amqp.Queue
		err   error
	)
	for attempt := 0; attempt <= cfg.Retries; {
		time.Sleep(5 * time.Second)
		conn, err = amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d", cfg.User, cfg.Password, cfg.Host, cfg.Port))
		if err != nil {
			if attempt < cfg.Retries {
				continue
			}
			return nil, errors.Wrap(err, "failed connect to server")
		}
		ch, err = conn.Channel()
		if err != nil {
			if attempt < cfg.Retries {
				continue
			}
			return nil, errors.Wrap(err, "failed open channel")
		}
		queue, err = ch.QueueDeclare(cfg.Queue, true, false, false, false, nil)
		if err != nil {
			if attempt < cfg.Retries {
				continue
			}
			return nil, errors.Wrap(err, "failed top queue")
		}
		break
	}

	m := RabbitMQ{
		conn:  conn,
		ch:    ch,
		queue: &queue,
	}

	go func() {
		<-ctx.Done()
		m.Close()
	}()

	return &m, nil
}

// Close close connection
func (m *RabbitMQ) Close() error {
	if m.ch != nil {
		if err := m.ch.Close(); err != nil {
			return errors.Wrap(err, "failed to close amqp channel")
		}
	}
	if m.conn != nil {
		if err := m.conn.Close(); err != nil {
			return errors.Wrap(err, "failed to close amqp connection")
		}
	}
	return nil
}

// SendMessage send message to remote queue
func (m *RabbitMQ) SendMessage(ctx context.Context, message []byte) error {
	if err := m.ch.Publish("", m.queue.Name, false, false, amqp.Publishing{
		ContentType: msgContentType,
		Body:        message,
	}); err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}

// ReceiveMessages receive message from remote queue
func (m *RabbitMQ) ReceiveMessages(ctx context.Context,
	processFunc func(internalCtx context.Context, message []byte) error) error {
	msgChannel, err := m.ch.Consume(m.queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return errors.Wrap(err, "failed to start consuming messages")
	}

	for msg := range msgChannel {

		if msg.ContentType != msgContentType {
			if err := msg.Ack(true); err != nil {
				return errors.Wrap(err, "failed to send ACK for message")
			}
			continue
		}
		if err := processFunc(ctx, msg.Body); err != nil {
			return errors.Wrap(err, "error from process func")
		}
		if err := msg.Ack(true); err != nil {
			return errors.Wrap(err, "failed to send ACK for message")
		}
	}

	return nil
}
