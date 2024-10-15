package logger

import (
	"log/slog"
	"os"
)

type Logger struct {
	logger *slog.Logger
}

func NewLogger() *Logger {

	customHandler := slog.NewJSONHandler(os.Stdout, nil)
	l := slog.New(customHandler)
	slog.SetDefault(l)

	return &Logger{
		logger: l,
	}
}
