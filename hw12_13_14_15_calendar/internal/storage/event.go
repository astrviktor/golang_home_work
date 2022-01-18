package storage

import (
	"errors"
	"time"
)

type Event struct {
	ID                 string    // ID - уникальный идентификатор события (можно воспользоваться UUID)
	Title              string    // Заголовок - короткий
	DateStart          time.Time // Дата и время начала события
	DateEnd            time.Time // Дата и время окончания события
	Description        string    // Описание события - длинный текст
	UserID             int       // ID пользователя, владельца события
	TimeToNotification int       // За сколько минут высылать уведомление
}

var ErrDateTimeBusy = errors.New("это время занято другим событием")
