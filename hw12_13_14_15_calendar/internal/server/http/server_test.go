package internalhttp

import (
	"context"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/astrviktor/golang_home_work/hw12_13_14_15_calendar/internal/app"
	"github.com/astrviktor/golang_home_work/hw12_13_14_15_calendar/internal/logger"
	generate "github.com/astrviktor/golang_home_work/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/astrviktor/golang_home_work/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
)

func TestHTTPServer(t *testing.T) {
	logg := logger.New(4, "2006-01-02T15:04:05Z07:00")

	storage := memorystorage.New()
	calendar := app.New(logg, storage)

	httpServer := NewServer(logg, calendar, storage, "127.0.0.1", "7777")

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go httpServer.Start(ctx)
	time.Sleep(2 * time.Second)

	client := NewClient("127.0.0.1", "7777", time.Second)

	date, err := time.Parse("2006-01-02", "2022-01-01")
	require.NoError(t, err)

	t.Run("http test create event and get event", func(t *testing.T) {
		newEvent := generate.GenerateEventDate(date, date.Add(time.Second))

		id, err := client.CreateEvent(newEvent)
		require.NoError(t, err)
		require.Len(t, id, 36)
		newEvent.ID = id

		getEvent, err := client.GetEvent(id)
		require.NoError(t, err)

		require.Equal(t, newEvent, getEvent)
	})

	t.Run("http test create event and update event", func(t *testing.T) {
		newEvent := generate.GenerateEventDate(date.Add(2*time.Second), date.Add(3*time.Second))

		id, err := client.CreateEvent(newEvent)
		require.NoError(t, err)
		require.Len(t, id, 36)
		newEvent.ID = id

		updEvent := generate.GenerateEventDate(date.Add(4*time.Second), date.Add(5*time.Second))
		updEvent.ID = id
		ok, err := client.UpdateEvent(updEvent)
		require.NoError(t, err)
		require.True(t, ok)

		getEvent, err := client.GetEvent(id)
		require.NoError(t, err)

		require.Equal(t, updEvent, getEvent)
	})

	t.Run("http test create event and delete event", func(t *testing.T) {
		newEvent := generate.GenerateEventDate(date.Add(6*time.Second), date.Add(7*time.Second))

		id, err := client.CreateEvent(newEvent)
		require.NoError(t, err)
		require.Len(t, id, 36)
		newEvent.ID = id

		ok, err := client.DeleteEvent(id)
		require.NoError(t, err)
		require.True(t, ok)

		ok, err = client.DeleteEvent(id)
		require.Error(t, err)
		require.False(t, ok)
	})

	t.Run("http test create and list", func(t *testing.T) {
		dateA, err := time.Parse("2006-01-02", "2022-02-07")
		require.NoError(t, err)
		dateB, err := time.Parse("2006-01-02", "2022-02-08")
		require.NoError(t, err)
		dateC, err := time.Parse("2006-01-02", "2022-02-14")
		require.NoError(t, err)

		eventA := generate.GenerateEventDate(dateA, dateA.Add(time.Hour))
		id, err := client.CreateEvent(eventA)
		require.NoError(t, err)
		require.Len(t, id, 36)
		eventA.ID = id

		eventB := generate.GenerateEventDate(dateB, dateB.Add(time.Hour))
		id, err = client.CreateEvent(eventB)
		require.NoError(t, err)
		require.Len(t, id, 36)
		eventB.ID = id

		eventC := generate.GenerateEventDate(dateC, dateC.Add(time.Hour))
		id, err = client.CreateEvent(eventC)
		require.NoError(t, err)
		require.Len(t, id, 36)
		eventC.ID = id

		events, err := client.GetListDay("2022-02-08")
		require.NoError(t, err)
		require.Equal(t, 1, len(events))

		events, err = client.GetListWeek("2022-02-07")
		require.NoError(t, err)
		require.Equal(t, 2, len(events))

		events, err = client.GetListMonth("2022-02-01")
		require.NoError(t, err)
		require.Equal(t, 3, len(events))
	})
}
