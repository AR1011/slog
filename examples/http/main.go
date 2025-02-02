package main

import (
	"github.com/AR1011/slog"
)

func main() {
	// Configure HTTP writer to send logs to a remote endpoint
	slog.AddWriter(slog.WithToHttpWriter(&slog.ToHttpWriterOptions{
		URL:    "http://localhost:8080/logs",
		Method: "POST",
		Format: slog.FormatJson,
		Level:  slog.InfoLevel,
		APIKey: "your-api-key-here",
	}))

	// Log some messages
	slog.Info("System startup")
	slog.Error("Connection failed", "service", "database", "error", "timeout")

	// Make sure to close the logger to flush any remaining logs
	defer slog.Close()
}
