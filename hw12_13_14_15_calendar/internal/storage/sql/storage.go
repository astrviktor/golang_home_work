package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/astrviktor/golang_home_work/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/stdlib" //nolint
)

type Storage struct {
	dsn string
	db  *sql.DB
}

func New(dsn string) *Storage {
	return &Storage{dsn, nil}
}

func (s *Storage) Connect(ctx context.Context) error {
	db, err := sql.Open("pgx", s.dsn)
	if err != nil {
		return err
	}

	err = db.PingContext(ctx)
	if err != nil {
		return err
	}

	s.db = db
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	return s.db.Close()
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

	tx, err := s.db.Begin()
	if err != nil {
		return "", err
	}

	sqlStatement := `INSERT INTO calendar.event
    (id, title, date_start, date_end, description, user_id, time_to_notification)
	VALUES ($1, $2, $3, $4, $5, $6, $7);`

	_, err = tx.Exec(sqlStatement, ID, event.Title, event.DateStart, event.DateEnd,
		event.Description, event.UserID, event.TimeToNotification)
	if err != nil {
		return "", err
	}

	err = tx.Commit()
	if err != nil {
		return "", err
	}

	return ID, nil
}

// Update - Обновить (событие).
func (s *Storage) Update(event storage.Event) (bool, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return false, err
	}

	sqlStatement := `UPDATE calendar.event 
	SET title=$2, date_start=$3, date_end=$4, description=$5, user_id=$6, time_to_notification=$7
	WHERE id = $1;`

	_, err = tx.Exec(sqlStatement, event.ID, event.Title, event.DateStart, event.DateEnd,
		event.Description, event.UserID, event.TimeToNotification)
	if err != nil {
		return false, err
	}

	err = tx.Commit()
	if err != nil {
		return false, err
	}

	return true, nil
}

// Delete - Удалить (ID события).
func (s *Storage) Delete(id string) (bool, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return false, err
	}

	sqlStatement := `DELETE FROM calendar.event WHERE id = $1;`

	_, err = tx.Exec(sqlStatement, id)
	if err != nil {
		return false, err
	}

	err = tx.Commit()
	if err != nil {
		return false, err
	}

	return true, nil
}

// Clear - очистка всех событий.
func (s *Storage) Clear() error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	sqlStatement := `DELETE FROM calendar.event;`

	_, err = tx.Exec(sqlStatement)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// Get - Получить событие (ID события).
func (s *Storage) Get(id string) (storage.Event, bool, error) {
	var event storage.Event

	sqlStatement := `SELECT id, title, date_start, date_end, description, user_id, time_to_notification
	FROM calendar.event 
	WHERE id = $1;`

	rows, err := s.db.Query(sqlStatement, id)
	if err != nil {
		return event, false, err
	}
	defer rows.Close()

	rows.Next()
	err = rows.Scan(&event.ID, &event.Title, &event.DateStart, &event.DateEnd, &event.Description,
		&event.UserID, &event.TimeToNotification)
	if errors.Is(err, sql.ErrNoRows) {
		return event, false, nil
	}

	if err != nil {
		return event, false, err
	}

	if rows.Err() != nil {
		return event, false, err
	}

	return event, true, nil
}

// EventListStartEnd - Список событий со старта (дата) по окончание (дата).
func (s *Storage) EventListStartEnd(start time.Time, end time.Time) ([]storage.Event, error) {
	events := make([]storage.Event, 0)

	sqlStatement := `SELECT id, title, date_start, date_end, description, user_id, time_to_notification
	FROM calendar.event 
	WHERE (date_start >= $1 AND date_start < $2) OR (date_end > $1 AND date_end < $2)
	ORDER BY date_start;`

	rows, err := s.db.Query(sqlStatement, start.Format(time.RFC3339), end.Format(time.RFC3339))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var event storage.Event

		err := rows.Scan(&event.ID, &event.Title, &event.DateStart, &event.DateEnd, &event.Description,
			&event.UserID, &event.TimeToNotification)

		if errors.Is(err, sql.ErrNoRows) {
			return events, nil
		}

		if err != nil {
			return events, err
		}

		if rows.Err() != nil {
			return events, err
		}

		events = append(events, event)
	}

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
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	sqlStatement := `UPDATE calendar.event 
	SET notified=true
	WHERE id = $1;`

	_, err = tx.Exec(sqlStatement, id)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetForNotification(date time.Time) ([]storage.Notification, error) {
	notifications := make([]storage.Notification, 0)

	sqlStatement := `SELECT id, title, date_start, user_id
	FROM calendar.event 
	WHERE notified = false AND date_start - time_to_notification * interval '1 minute' < $1
	ORDER BY date_start;`

	rows, err := s.db.Query(sqlStatement, date.Format(time.RFC3339))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var notification storage.Notification

		err := rows.Scan(&notification.ID, &notification.Title, &notification.DateStart, &notification.UserID)

		if errors.Is(err, sql.ErrNoRows) {
			return notifications, nil
		}

		if err != nil {
			return notifications, err
		}

		if rows.Err() != nil {
			return notifications, err
		}

		notifications = append(notifications, notification)
	}

	return notifications, nil
}
