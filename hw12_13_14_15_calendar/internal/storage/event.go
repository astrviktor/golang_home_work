package storage

import (
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
}

var ErrDateTimeBusy = errors.New("это время занято другим событием")
