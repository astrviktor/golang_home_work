package logger

import (
	"fmt"
	"os"
)

type Logger struct {
	level      int
	timeFormat string
}

// INFO — логирование ошибок, предупреждений и сообщений.
// DEBUG — логирование всех событий при отладке.
// WARN — логирование ошибок и предупреждений.
// ERROR — логирование ошибок.

const (
	INFO  = 1
	DEBUG = 2
	WARN  = 3
	ERROR = 4
)

func New(level int, timeFormat string) *Logger {
	return &Logger{level, timeFormat}
}

func (l *Logger) Info(msg string) {
	if l.level <= INFO {
		fmt.Printf("INFO: %s\n", msg)
	}
}

func (l *Logger) Debug(msg string) {
	if l.level <= DEBUG {
		fmt.Printf("DEBUG: %s\n", msg)
	}
}

func (l *Logger) Warn(msg string) {
	if l.level <= WARN {
		fmt.Printf("WARN: %s", msg)
	}
}

func (l *Logger) Error(msg string) {
	if l.level <= ERROR {
		fmt.Printf("ERROR: %s\n", msg)
	}
}

func (l *Logger) Fatal(msg string) {
	if l.level <= ERROR {
		fmt.Printf("ERROR: %s", msg)
		os.Exit(1)
	}
}

func (l *Logger) GetTimeFormat() string {
	return l.timeFormat
}
