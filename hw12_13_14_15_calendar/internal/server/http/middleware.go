package internalhttp

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

func (r *StatusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

func (s *Server) WithLogging(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		now := "[" + time.Now().Format(s.logger.GetTimeFormat()) + "]"
		start := time.Now()
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		userAgent := r.UserAgent()

		recorder := &StatusRecorder{
			ResponseWriter: w,
			Status:         0,
		}

		h(recorder, r)
		s.logger.Info(fmt.Sprint(ip, now, r.Method, r.RequestURI,
			r.Proto, recorder.Status, time.Since(start), userAgent))
	}
}
