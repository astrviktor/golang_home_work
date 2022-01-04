package storage

import (
	"errors"
	"math/rand"
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

func getRandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789") //nolint:gofumpt

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))] //nolint:gosec
	}
	return string(s)
}

func GenerateEvent() Event {
	date := time.Now()
	date = time.Date(date.Year(), date.Month(), date.Day(),
		date.Hour(), date.Minute(), date.Second(), 0, date.Location())

	return Event{
		ID:                 "",
		Title:              getRandomString(rand.Intn(10) + 10), //nolint:gosec
		DateStart:          date,
		DateEnd:            date.Add(time.Second),
		Description:        getRandomString(rand.Intn(100) + 100), //nolint:gosec
		UserID:             rand.Intn(100),                        //nolint:gosec
		TimeToNotification: rand.Intn(30) + 30,                    //nolint:gosec
	}
}
