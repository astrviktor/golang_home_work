package memorystorage

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/astrviktor/golang_home_work/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
)

type Storage struct {
	events map[string]storage.Event
	mutex  *sync.RWMutex
}

func New() *Storage {
	mutex := sync.RWMutex{}
	events := make(map[string]storage.Event)
	return &Storage{events, &mutex}
}

// Connect - соединение.
func (s *Storage) Connect(ctx context.Context) error {
	return nil
}

// Close - закрытие соединения.
func (s *Storage) Close(ctx context.Context) error {
	return nil
}

// Create - Создать (событие).
func (s *Storage) Create(event storage.Event) (string, error) {
	dateStart := event.DateStart
	dateEnd := event.DateEnd
	events, err := s.EventListStartEnd(dateStart, dateEnd)
	if err != nil {
		return "", err
	}

	if len(events) != 0 {
		return "", storage.ErrDateTimeBusy
	}

	ID := uuid.New().String()
	newEvent := storage.Event{
		ID:                 ID,
		Title:              event.Title,
		DateStart:          event.DateStart,
		DateEnd:            event.DateEnd,
		Description:        event.Description,
		UserID:             event.UserID,
		TimeToNotification: event.TimeToNotification,
		Notified:           event.Notified,
	}

	s.mutex.Lock()
	s.events[ID] = newEvent
	s.mutex.Unlock()

	return ID, nil
}

// Update - Обновить (событие).
func (s *Storage) Update(event storage.Event) (bool, error) {
	id := event.ID

	s.mutex.Lock()
	s.events[id] = event
	s.mutex.Unlock()

	return true, nil
}

// Delete - Удалить (ID события).
func (s *Storage) Delete(id string) (bool, error) {
	s.mutex.Lock()
	_, ok := s.events[id]
	if ok {
		delete(s.events, id)
	}
	s.mutex.Unlock()
	return ok, nil
}

// Get - Получить событие (ID события).
func (s *Storage) Get(id string) (storage.Event, bool, error) {
	s.mutex.Lock()
	event, ok := s.events[id]
	s.mutex.Unlock()

	return event, ok, nil
}

// EventListStartEnd - Список событий со старта (дата) по окончание (дата).
func (s *Storage) EventListStartEnd(start time.Time, end time.Time) ([]storage.Event, error) {
	events := make([]storage.Event, 0)

	s.mutex.Lock()
	for _, event := range s.events {
		if (event.DateStart.Equal(start) || (event.DateStart.After(start)) && event.DateStart.Before(end)) ||
			(event.DateEnd.After(start) && event.DateEnd.Before(end)) {
			events = append(events, event)
		}
	}
	s.mutex.Unlock()

	sort.SliceStable(events, func(i, j int) bool {
		return events[i].DateStart.Before(events[j].DateStart)
	})

	return events, nil
}

// EventListDay - СписокСобытийНаДень (дата).
func (s *Storage) EventListDay(date time.Time) ([]storage.Event, error) {
	dateStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	dateEnd := dateStart.Add(24 * time.Hour)

	return s.EventListStartEnd(dateStart, dateEnd)
}

// EventListWeek - СписокСобытийНаНеделю (дата начала недели).
func (s *Storage) EventListWeek(date time.Time) ([]storage.Event, error) {
	dateStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	dateEnd := dateStart.Add(7 * 24 * time.Hour)

	return s.EventListStartEnd(dateStart, dateEnd)
}

// EventListMonth - СписокСобытийНaМесяц (дата начала месяца).
func (s *Storage) EventListMonth(date time.Time) ([]storage.Event, error) {
	dateStart := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
	dateEnd := dateStart.AddDate(0, 1, 0)

	return s.EventListStartEnd(dateStart, dateEnd)
}

func (s *Storage) Notified(id string) error {
	s.mutex.Lock()
	event, ok := s.events[id]
	if ok {
		event.Notified = "yes"
		s.events[id] = event
	}
	s.mutex.Unlock()

	return nil
}

func (s *Storage) GetForNotification(date time.Time) ([]storage.Notification, error) {
	notifications := make([]storage.Notification, 0)

	s.mutex.Lock()
	for _, event := range s.events {
		duration := time.Duration(event.TimeToNotification) * time.Minute
		dateToNotification := event.DateStart.Add(-duration)
		if event.Notified == "no" && dateToNotification.Before(date) {
			notification := storage.Notification{
				ID:        event.ID,
				Title:     event.Title,
				DateStart: event.DateStart,
				UserID:    event.UserID,
			}

			notifications = append(notifications, notification)
		}
	}
	s.mutex.Unlock()

	sort.SliceStable(notifications, func(i, j int) bool {
		return notifications[i].DateStart.Before(notifications[j].DateStart)
	})

	return notifications, nil
}

func (s *Storage) DeleteOlder(date time.Time) error {
	s.mutex.Lock()
	for _, event := range s.events {
		if date.After(event.DateStart) {
			delete(s.events, event.ID)
		}
	}
	s.mutex.Unlock()

	return nil
}
