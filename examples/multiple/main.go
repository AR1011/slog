package main

import (
	"github.com/AR1011/slog"
)

func main() {
	// Create writers for both console and file output
	consoleWriter := slog.WithStdIoWriter(&slog.ToStdStreamWriterOptions{
		Stream: slog.StdOut,
		Format: slog.FormatAnsi,
		Level:  slog.DebugLevel,
	})

	fileWriter := slog.WithToFileWriter(&slog.ToFileWriterOptions{
		FileName:   "multiple.example.log",
		Format:     slog.FormatJson,
		Level:      slog.InfoLevel,
		RotateSize: 1024 * 1024 * 10, // 10MB
	})

	// Create a new logger with both writers
	logger, err := slog.NewLogger(consoleWriter, fileWriter)
	if err != nil {
		panic(err)
	}

	// Set as the default logger
	slog.SetLogger(logger)

	// Log messages will go to both writers based on their levels
	slog.Debug("This only goes to console")          // Only appears in console
	slog.Info("This goes to both console and file")  // Appears in both
	slog.Error("Critical system error", "code", 500) // Appears in both with extra context

	// Don't forget to close the logger
	defer slog.Close()
}
