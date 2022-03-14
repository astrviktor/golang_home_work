package queue

import (
	"context"

	"github.com/astrviktor/golang_home_work/hw12_13_14_15_calendar/internal/storage"
)

type Message struct {
	Data []byte
}

type Producer interface {
	SendMessage(msg []byte) error
	Produce(storage storage.Storage) error
	Close() error
}

type Consumer interface {
	Consume(ctx context.Context) (<-chan Message, error)
	Close() error
}
