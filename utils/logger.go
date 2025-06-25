package utils

import (
	"log"
	"os"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
)

type Logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	level       string
}

func NewLogger(level string) *Logger {
	return &Logger{
		infoLogger:  log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger: log.New(os.Stderr, ColorRed+"[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile),
		level:       level,
	}
}

func (l *Logger) Info(msg string) {
	if l.level == "info" || l.level == "debug" {
		l.infoLogger.Output(2, msg)
	}
}

func (l *Logger) Error(msg string) {
	if l.level == "error" || l.level == "info" || l.level == "debug" {
		l.errorLogger.Output(2, msg+ColorReset)
	}
}