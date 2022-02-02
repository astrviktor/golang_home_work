package internalhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
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

	t.Run("http test create event and get event", func(t *testing.T) {
		// create event
		eventCreate := generate.GenerateEvent()

		eventBytes, err := json.Marshal(eventCreate)
		require.NoError(t, err)

		body := bytes.NewReader(eventBytes)

		client := http.Client{Timeout: time.Second * 5}
		resp, err := client.Post("http://127.0.0.1:7777/event", "application/json", body) //nolint:noctx
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		buffer := make([]byte, 1024)
		n, err := resp.Body.Read(buffer)
		require.ErrorIs(t, err, io.EOF)
		bytes := buffer[:n]

		resultCreate := ResponseID{}
		err = json.Unmarshal(bytes, &resultCreate)
		require.NoError(t, err)
		eventCreate.ID = resultCreate.ID

		err = resp.Body.Close()
		require.NoError(t, err)

		// get event
		resp, err = client.Get("http://127.0.0.1:7777/event?id=" + resultCreate.ID) //nolint:noctx
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		n, err = resp.Body.Read(buffer)
		require.ErrorIs(t, err, io.EOF)
		bytes = buffer[:n]

		resultGet := ResponseEvent{}
		err = json.Unmarshal(bytes, &resultGet)
		require.NoError(t, err)

		eventGet := resultGet.Event
		require.Equal(t, eventCreate, eventGet)

		err = resp.Body.Close()
		require.NoError(t, err)
	})
}
