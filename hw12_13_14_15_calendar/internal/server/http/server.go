package internalhttp

import (
	"context"
	"net"
	"net/http"
)

type Server struct {
	logger Logger
	app    Application
	addr   string
}

type Logger interface {
	Debug(msg ...interface{})
	Info(msg ...interface{})
	Warn(msg ...interface{})
	Error(msg ...interface{})
	Fatal(msg ...interface{})
	WithLogging(h http.HandlerFunc) http.HandlerFunc
}

type Application interface { // TODO
}

func NewServer(logger Logger, app Application, host string, port string) *Server {
	return &Server{logger, app, net.JoinHostPort(host, port)}
}

func (s *Server) Start(ctx context.Context) {
	mux := http.NewServeMux()

	logger := s.logger.WithLogging(s.HelloWorld)
	mux.HandleFunc("/hello", logger)

	server := &http.Server{
		Addr:    s.addr,
		Handler: mux,
	}

	s.logger.Info("server starting on address " + s.addr)

	go func() {
		s.logger.Fatal(server.ListenAndServe())
	}()

	<-ctx.Done()
}

func (s *Server) Stop(ctx context.Context) error {
	// TODO
	return nil
}

func (s *Server) HelloWorld(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Hello World!"))
}

// TODO
