package logger

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

func Init(level, logFile string) {
	var logLevel slog.Level
	switch strings.ToLower(level) {
	case "debug":
		logLevel = slog.LevelDebug
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	var writers []io.Writer
	writers = append(writers, os.Stdout)

	if logFile != "" {
		dir := filepath.Dir(logFile)
		if err := os.MkdirAll(dir, 0755); err != nil {
			slog.Error("failed to create log directory", "error", err, "path", dir)
		} else {
			f, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				slog.Error("failed to open log file", "error", err, "path", logFile)
			} else {
				writers = append(writers, f)
			}
		}
	}

	mw := io.MultiWriter(writers...)
	handler := slog.NewJSONHandler(mw, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
