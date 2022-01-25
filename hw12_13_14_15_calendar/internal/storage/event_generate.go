package storage

import (
	"math/rand"
	"time"
)

func getRandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789") //nolint:gofumpt

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))] //nolint:gosec
	}
	return string(s)
}

func GenerateEvent() Event {
	date := time.Now().UTC()
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
