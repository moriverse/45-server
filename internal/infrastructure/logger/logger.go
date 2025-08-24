package logger

import (
	"io"
	"log/slog"
	"os"

	"github.com/moriverse/45-server/internal/infrastructure/config"
)

// NewLogger creates a new slog.Logger based on the application configuration.
func NewLogger(cfg config.LogConfig) *slog.Logger {
	var handler slog.Handler
	var level slog.Level

	switch cfg.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	var output io.Writer = os.Stdout

	switch cfg.Format {
	case "json":
		handler = slog.NewJSONHandler(output, opts)
	case "text":
		fallthrough
	default:
		handler = slog.NewTextHandler(output, opts)
	}

	return slog.New(handler)
}
