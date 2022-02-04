package consumer

import (
	"context"
	"fmt"
	"log"

	"github.com/astrviktor/golang_home_work/hw12_13_14_15_calendar/internal/queue"
	"github.com/streadway/amqp"
)

type Consumer struct {
	name    string
	channel *amqp.Channel
	queue   string
}

func NewConsumer(name string, uri, queue string) (*Consumer, error) {
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, fmt.Errorf("producer dial %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("producer channel %w", err)
	}

	consumer := Consumer{
		name:    name,
		channel: channel,
		queue:   queue,
	}

	return &consumer, nil
}

func (c *Consumer) Consume(ctx context.Context) (<-chan queue.Message, error) {
	messages := make(chan queue.Message)

	go func() {
		<-ctx.Done()
		if err := c.channel.Close(); err != nil {
			log.Println(err)
		}
	}()

	deliveries, err := c.channel.Consume(c.queue, c.name, false, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("start consuming %w", err)
	}

	go func() {
		defer func() {
			close(messages)
			log.Println("close messages channel")
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case del := <-deliveries:

				if err := del.Ack(false); err != nil {
					log.Println(err)
				}

				msg := queue.Message{
					Data: del.Body,
				}

				select {
				case <-ctx.Done():
					return
				case messages <- msg:
				}
			}
		}
	}()

	return messages, nil
}

func (c *Consumer) Close() error {
	if err := c.channel.Close(); err != nil {
		return fmt.Errorf("error on closing channel %w", err)
	}

	return nil
}
