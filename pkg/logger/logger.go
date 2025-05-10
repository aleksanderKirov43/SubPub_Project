package logger

import "log"

type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

type StandardLogger struct {
}

func (l *StandardLogger) Info(msg string, args ...interface{}) {
	log.Printf("INFO: "+msg, args...)
}

func (l *StandardLogger) Error(msg string, args ...interface{}) {
	log.Printf("ERROR: "+msg, args...)
}

func NewLogger() Logger {
	return &StandardLogger{}
}
