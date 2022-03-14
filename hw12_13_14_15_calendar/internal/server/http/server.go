package internalhttp

import (
	"context"
	"net"
	"net/http"

	"github.com/astrviktor/golang_home_work/hw12_13_14_15_calendar/internal/app"
	"github.com/astrviktor/golang_home_work/hw12_13_14_15_calendar/internal/storage"
)

type Server struct {
	logger  app.Logger
	app     Application
	storage storage.Storage
	addr    string
}

type Application interface { // TODO
}

func NewServer(logger app.Logger, app Application, storage storage.Storage, host string, port string) *Server {
	return &Server{logger, app, storage, net.JoinHostPort(host, port)}
}

// GET    /event?id=1     : возвращает event по ID
// POST   /event          : создаёт event из body
// PUT    /event          : обновляет event из body
// DELETE /event?id=1     : удаляет event по ID
// GET    /list/day?date=2021-01-01       : возвращает все event за день
// GET    /list/week?date=2021-01-01      : возвращает все event за неделю
// GET    /list/month?date=2021-01-01     : возвращает все event за месяц

func (s *Server) Start(ctx context.Context) {
	mux := http.NewServeMux()

	mux.HandleFunc("/event", s.WithLogging(s.Event))
	mux.HandleFunc("/list/day", s.WithLogging(s.GetListDay))
	mux.HandleFunc("/list/week", s.WithLogging(s.GetListWeek))
	mux.HandleFunc("/list/month", s.WithLogging(s.GetListMonth))

	server := &http.Server{
		Addr:    s.addr,
		Handler: mux,
	}

	s.logger.Info("http server starting on address " + s.addr)

	go func() {
		s.logger.Fatal(server.ListenAndServe().Error())
	}()

	<-ctx.Done()
}

func (s *Server) Stop(ctx context.Context) error {
	// TODO
	return nil
}
