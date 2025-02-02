package main

import (
	"github.com/AR1011/slog"
)

func main() {
	// the default logger is already configured to use the stdout writer
	// so if you would like to change the default logger, you can do so by calling slog.SetLogger()

	logger, err := slog.NewLogger(slog.WithStdIoWriter(&slog.ToStdStreamWriterOptions{
		Stream: slog.StdOut,
		Format: slog.FormatAnsi,
		Level:  slog.DebugLevel,
	}))
	if err != nil {
		panic(err)
	}

	slog.SetLogger(logger)

	// Log messages with different levels
	slog.Debug("This is a debug message")
	slog.Info("This is an info message")
	slog.Warn("This is a warning message")
	slog.Error("This is an error message")

	// This is a debug message, here is a number 15            foo=bar
	slog.DebugF("This is a %s messages, here is a number %d", []interface{}{"warning", 1}, "foo", "bar")
	slog.InfoF("This is an %s message, here is a number %d", []interface{}{"info", 2}, "foo", "bar")
	slog.WarnF("This is a %s message, here is a number %d", []interface{}{"warning", 3}, "foo", "bar")
	slog.ErrorF("This is an %s message, here is a number %d", []interface{}{"error", 4}, "foo", "bar")

	// this will call panic() after logging the message
	// slog.Panic("This is a panic message")

	// Log with additional context
	slog.Info("User action", "user_id", 123, "action", "login")
}
