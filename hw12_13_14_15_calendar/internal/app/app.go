package app

import (
	"github.com/astrviktor/golang_home_work/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	logger  Logger
	storage storage.Storage
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Fatal(msg string)
	GetTimeFormat() string
}

func New(logger Logger, storage storage.Storage) *App {
	return &App{logger, storage}
}
