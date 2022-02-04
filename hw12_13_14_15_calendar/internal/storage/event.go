package storage

import (
	"context"
	"errors"
	"time"
)

type Event struct {
	ID                 string    `json:"id"`                 // ID - уникальный идентификатор события (UUID)
	Title              string    `json:"title"`              // Заголовок - короткий
	DateStart          time.Time `json:"dateStart"`          // Дата и время начала события
	DateEnd            time.Time `json:"dateEnd"`            // Дата и время окончания события
	Description        string    `json:"description"`        // Описание события - длинный текст
	UserID             int       `json:"usedId"`             // ID пользователя, владельца события
	TimeToNotification int       `json:"timeToNotification"` // За сколько минут высылать уведомление
	Notified           bool      `json:"notified"`           // Было ли отправлено уведомление
}

var ErrDateTimeBusy = errors.New("это время занято другим событием")

type Storage interface {
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
	Create(event Event) (string, error)
	Update(event Event) (bool, error)
	Delete(id string) (bool, error)
	Get(id string) (Event, bool, error)
	EventListStartEnd(start time.Time, end time.Time) ([]Event, error)
	EventListDay(date time.Time) ([]Event, error)
	EventListWeek(date time.Time) ([]Event, error)
	EventListMonth(date time.Time) ([]Event, error)
	Notified(id string) error
	GetForNotification(date time.Time) ([]Notification, error)
	DeleteOlder(date time.Time) error
}
