package app

import (
	"context"
	"time"

	"github.com/astrviktor/golang_home_work/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	logger  Logger
	storage Storage
	// TODO
}

type Logger interface { // TODO
}

type Storage interface {
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
	Create(event storage.Event) (string, error)
	Update(id string, event storage.Event) (bool, error)
	Delete(id string) (bool, error)
	Get(id string) (storage.Event, bool, error)
	EventListStartEnd(start time.Time, end time.Time) ([]storage.Event, error)
	EventListDay(date time.Time) ([]storage.Event, error)
	EventListWeek(date time.Time) ([]storage.Event, error)
	EventListMonth(date time.Time) ([]storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{logger, storage}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
