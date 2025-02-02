package main

import (
	"time"

	"github.com/AR1011/slog"
)

func main() {
	// Configure a file writer for application logs
	slog.AddWriter(slog.WithToFileWriter(&slog.ToFileWriterOptions{
		FileName: "advanced.example.log",
		Format:   slog.FormatJson,
		Level:    slog.DebugLevel,
	}))

	// Configure console output for immediate feedback
	slog.AddWriter(slog.WithStdIoWriter(&slog.ToStdStreamWriterOptions{
		Stream: slog.StdOut,
		Format: slog.FormatAnsi,
		Level:  slog.DebugLevel,
	}))

	// Log with various formats and contexts
	slog.Info("Application started", "version", "1.0.0", "env", "production")

	// Log different levels with context
	slog.Debug("Configuration loaded", "config_file", "config.yaml")
	slog.Warn("High memory usage", "memory_mb", 1500, "threshold_mb", 1000)

	// Simulating some work
	time.Sleep(time.Second)

	// Error with context
	slog.Error("Database connection failed",
		"error", "connection refused",
		"retry_count", 3,
		"last_attempt", time.Now().Format(time.RFC3339),
	)

	// Don't forget to close the logger
	defer slog.Close()
}
