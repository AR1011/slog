package main

import (
	"github.com/AR1011/slog"
)

func main() {
	// Configure file writer with JSON format and rotation
	slog.AddWriter(slog.WithToFileWriter(&slog.ToFileWriterOptions{
		FileName:   "file.example.log",
		Format:     slog.FormatJson,
		Level:      slog.InfoLevel,
		RotateSize: 1024 * 1024 * 100, // 100MB
	}))

	// Log some messages
	slog.Info("Application started")
	slog.Error("Connection failed", "service", "database", "error", "timeout")

	// Don't forget to close the logger
	defer slog.Close()
}
