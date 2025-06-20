package utils

import (
	"log"
	"os"
)

type Logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	level       string
}

// NewLogger creates a new Logger based on the provided level ("info", "error", "debug")
func NewLogger(level string) *Logger {
	return &Logger{
		infoLogger:  log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger: log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile),
		level:       level,
	}
}

// Info logs messages only if level is "info" or "debug"
func (l *Logger) Info(msg string) {
	if l.level == "info" || l.level == "debug" {
		l.infoLogger.Println(msg)
	}
}

// Error logs error messages (always if level is error/info/debug)
func (l *Logger) Error(msg string) {
	if l.level == "error" || l.level == "info" || l.level == "debug" {
		l.errorLogger.Println(msg)
	}
}
