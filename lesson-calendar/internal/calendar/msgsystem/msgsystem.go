package msgsystem

import (
	"context"
)

// MsgSystem message system
type MsgSystem interface {
	Close() error
	SendMessage(ctx context.Context, message []byte) error
	ReceiveMessages(ctx context.Context, processFunc func(ctx context.Context, message []byte) error) error
}
