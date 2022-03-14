package producer

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/astrviktor/golang_home_work/hw12_13_14_15_calendar/internal/storage"
	"github.com/streadway/amqp"
)

type Producer struct {
	name     string
	channel  *amqp.Channel
	exchange string
}

func NewProducer(name, uri, exchange string) (*Producer, error) {
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, fmt.Errorf("producer dial %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("producer channel %w", err)
	}

	producer := Producer{
		name:     name,
		channel:  channel,
		exchange: exchange,
	}

	return &producer, nil
}

func (p *Producer) SendMessage(msg []byte) error {
	if err := p.channel.Publish(
		p.exchange, // publish to an exchange
		"",         // routingKey - routing to 0 or more queues
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            msg,
			DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,              // 0-9
			// a bunch of application/implementation-specific fields
		},
	); err != nil {
		return fmt.Errorf("exchange publish: %w", err)
	}

	return nil
}

func (p *Producer) Produce(storage storage.Storage) error {
	if err := storage.DeleteOlder(time.Now().AddDate(-1, 0, 0)); err != nil {
		return fmt.Errorf("error deleting events older then year %w", err)
	}

	notifications, err := storage.GetForNotification(time.Now())
	if err != nil {
		return fmt.Errorf("error get notifications from storage %w", err)
	}

	if len(notifications) == 0 {
		log.Printf("nothing for produce")
		return nil
	}

	for _, notification := range notifications {
		msg, err := json.Marshal(notification)
		if err != nil {
			log.Println("error marshal notification:", err.Error())
			break
		}
		log.Println("sending message:", string(msg))

		if err = p.SendMessage(msg); err != nil {
			log.Println("error sending message:", err.Error())
			break
		}

		if err = storage.Notified(notification.ID); err != nil {
			log.Println("error notified message:", err.Error())
			break
		}
	}

	return nil
}

func (p *Producer) Close() error {
	if err := p.channel.Close(); err != nil {
		return fmt.Errorf("error on closing channel %w", err)
	}

	return nil
}
