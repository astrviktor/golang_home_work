package logger

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"time"
)

type Logger struct {
	level      string
	timeFormat string
}

// DEBUG — логирование всех событий при отладке.
// INFO — логирование ошибок, предупреждений и сообщений.
// WARN — логирование ошибок и предупреждений.
// ERROR — логирование всех ошибок.

const (
	DEBUG = "DEBUG"
	INFO  = "INFO"
	ERROR = "ERROR"
	WARN  = "WARN"
)

func New(level string, timeFormat string) *Logger {
	return &Logger{level, timeFormat}
}

func (l *Logger) Debug(msg ...interface{}) {
	if l.level == DEBUG {
		fmt.Print("DEBUG: ")
		fmt.Println(msg...)
	}
}

func (l *Logger) Info(msg ...interface{}) {
	if l.level == DEBUG || l.level == INFO || l.level == ERROR {
		fmt.Print("INFO: ")
		fmt.Println(msg...)
	}
}

func (l *Logger) Warn(msg ...interface{}) {
	if l.level == DEBUG || l.level == INFO || l.level == WARN || l.level == ERROR {
		fmt.Print("WARN: ")
		fmt.Println(msg...)
	}
}

func (l *Logger) Error(msg ...interface{}) {
	fmt.Print("ERROR: ")
	fmt.Println(msg...)
}

func (l *Logger) Fatal(msg ...interface{}) {
	fmt.Print("ERROR: ")
	fmt.Println(msg...)
	os.Exit(1)
}

type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

func (r *StatusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

func (l *Logger) WithLogging(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		now := "[" + time.Now().Format(l.timeFormat) + "]"
		start := time.Now()
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		userAgent := r.UserAgent()

		recorder := &StatusRecorder{
			ResponseWriter: w,
			Status:         0,
		}

		h(recorder, r)
		l.Info(ip, now, r.Method, r.RequestURI, r.Proto, recorder.Status, time.Since(start), userAgent)
	}
}
