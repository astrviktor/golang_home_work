package sqlstorage

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/astrviktor/golang_home_work/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

const dsn = "postgres://user:password123@localhost:5432/calendar"

func TestStorage(t *testing.T) {
	t.Skip() // Remove for SQL tests.

	t.Run("sql storage create and get", func(t *testing.T) {
		testStorage := New(dsn)
		ctx := context.Background()
		err := testStorage.Connect(ctx)
		require.NoError(t, err)

		err = testStorage.Clear()
		require.NoError(t, err)

		newEvent := storage.GenerateEvent()
		uuid, err := testStorage.Create(newEvent)
		require.NoError(t, err)
		newEvent.ID = uuid

		getEvent, ok, err := testStorage.Get(uuid)
		require.NoError(t, err)
		require.True(t, ok)

		require.Equal(t, newEvent, getEvent)

		err = testStorage.Close(ctx)
		require.NoError(t, err)
	})

	t.Run("sql storage create and update", func(t *testing.T) {
		testStorage := New(dsn)
		ctx := context.Background()
		err := testStorage.Connect(ctx)
		require.NoError(t, err)

		err = testStorage.Clear()
		require.NoError(t, err)

		newEvent := storage.GenerateEvent()
		uuid, err := testStorage.Create(newEvent)
		require.NoError(t, err)
		newEvent.ID = uuid

		updEvent := storage.GenerateEvent()
		updEvent.ID = uuid
		ok, err := testStorage.Update(updEvent)
		require.NoError(t, err)
		require.True(t, ok)

		getEvent, ok, err := testStorage.Get(uuid)
		require.NoError(t, err)
		require.True(t, ok)

		require.Equal(t, updEvent, getEvent)

		err = testStorage.Close(ctx)
		require.NoError(t, err)
	})

	t.Run("sql storage create and delete", func(t *testing.T) {
		testStorage := New(dsn)
		ctx := context.Background()
		err := testStorage.Connect(ctx)
		require.NoError(t, err)

		err = testStorage.Clear()
		require.NoError(t, err)

		newEvent := storage.GenerateEvent()
		uuid, err := testStorage.Create(newEvent)
		require.NoError(t, err)
		newEvent.ID = uuid

		ok, err := testStorage.Delete(uuid)
		require.NoError(t, err)
		require.True(t, ok)

		getEvent, ok, err := testStorage.Get(uuid)
		require.NoError(t, err)
		require.False(t, ok)

		require.Equal(t, getEvent, storage.Event{})

		err = testStorage.Close(ctx)
		require.NoError(t, err)
	})

	t.Run("sql storage create and list", func(t *testing.T) {
		testStorage := New(dsn)
		ctx := context.Background()
		err := testStorage.Connect(ctx)
		require.NoError(t, err)

		err = testStorage.Clear()
		require.NoError(t, err)

		var newEvents []storage.Event

		for i := 0; i < 10; i++ {
			newEvent := storage.GenerateEvent()

			newEvent.DateStart = newEvent.DateStart.Add(time.Duration(i) * time.Hour)
			newEvent.DateEnd = newEvent.DateStart.Add(10 * time.Minute)

			uuid, err := testStorage.Create(newEvent)
			require.NoError(t, err)
			newEvent.ID = uuid

			newEvents = append(newEvents, newEvent)
		}

		getEvents, err := testStorage.EventListStartEnd(time.Now().Add(-24*time.Hour), time.Now().Add(24*time.Hour))
		require.NoError(t, err)
		require.Equal(t, 10, len(getEvents))
		require.Equal(t, newEvents, getEvents)

		err = testStorage.Close(ctx)
		require.NoError(t, err)
	})

	t.Run("sql storage create and ErrDateTimeBusy error", func(t *testing.T) {
		testStorage := New(dsn)
		ctx := context.Background()
		err := testStorage.Connect(ctx)
		require.NoError(t, err)

		err = testStorage.Clear()
		require.NoError(t, err)

		newEvent := storage.GenerateEvent()

		newEvent.DateEnd = newEvent.DateStart.Add(10 * time.Minute)
		_, err = testStorage.Create(newEvent)
		require.NoError(t, err)

		newEvent.DateStart = newEvent.DateStart.Add(5 * time.Minute)
		newEvent.DateEnd = newEvent.DateStart.Add(10 * time.Minute)
		_, err = testStorage.Create(newEvent)
		require.Error(t, err)

		require.ErrorIs(t, err, storage.ErrDateTimeBusy)

		err = testStorage.Close(ctx)
		require.NoError(t, err)
	})
}

func TestStorageMultithreading(t *testing.T) {
	t.Skip() // Remove for SQL tests.

	testStorage := New(dsn)
	ctx := context.Background()
	err := testStorage.Connect(ctx)
	require.NoError(t, err)

	err = testStorage.Clear()
	require.NoError(t, err)

	var uuidSet1, uuidSet2 []string

	wg := &sync.WaitGroup{}

	// Create

	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			uuid, _ := testStorage.Create(storage.GenerateEvent())
			uuidSet1 = append(uuidSet1, uuid)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			uuid, _ := testStorage.Create(storage.GenerateEvent())
			uuidSet2 = append(uuidSet2, uuid)
		}
	}()

	wg.Wait()

	// Get

	wg.Add(2)

	go func() {
		defer wg.Done()
		for _, uuid := range uuidSet1 {
			_, _, _ = testStorage.Get(uuid)
		}
	}()

	go func() {
		defer wg.Done()
		for _, uuid := range uuidSet2 {
			_, _, _ = testStorage.Get(uuid)
		}
	}()

	wg.Wait()

	// Update

	wg.Add(2)

	go func() {
		defer wg.Done()
		for _, uuid := range uuidSet1 {
			event := storage.GenerateEvent()
			event.ID = uuid
			_, _ = testStorage.Update(event)
		}
	}()

	go func() {
		defer wg.Done()
		for _, uuid := range uuidSet2 {
			event := storage.GenerateEvent()
			event.ID = uuid
			_, _ = testStorage.Update(event)
		}
	}()

	wg.Wait()

	// Delete

	wg.Add(2)

	go func() {
		defer wg.Done()
		for _, uuid := range uuidSet1 {
			_, _ = testStorage.Delete(uuid)
		}
	}()

	go func() {
		defer wg.Done()
		for _, uuid := range uuidSet2 {
			_, _ = testStorage.Delete(uuid)
		}
	}()

	wg.Wait()
}
