package logger

import (
	"log/slog"
	"os"
)

// NewLogger создает экземпляр логгера с указанным уровнем логирования
func NewLogger(level string) *slog.Logger {
	var handler slog.Handler
	switch level {
	case "debug":
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	case "prod":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	}
	return slog.New(handler)
}
