package main

import (
	"fmt"
	"time"

	"github.com/AR1011/slog"
)

func main() {
	// Create a channel to receive logs
	logChan := make(chan *slog.Log, 100)

	// Create a new logger with the channel writer
	logger, err := slog.NewLogger(slog.WithToChanWriter(&slog.ToChanWriterOptions{
		Ch:    logChan,
		Level: slog.InfoLevel,
	}))
	if err != nil {
		panic(err)
	}
	slog.SetLogger(logger)

	// Start a goroutine to process logs
	done := make(chan bool)
	go func() {
		for log := range logChan {
			// Custom processing of log entries
			fmt.Printf("Custom handler: [%s] %s\n", log.Type, log.Msg)
			if log.Args != nil {
				fmt.Printf("  Args: %v\n", log.Args)
			}
		}
		done <- true
	}()

	// Log some messages
	slog.Info("Processing started")
	slog.Warn("Resource usage high", "cpu", 85, "memory", "75%")
	slog.Error("System error detected", "code", 500, "component", "api")

	// Wait a moment to ensure logs are processed
	time.Sleep(100 * time.Millisecond)

	// Properly close everything
	slog.Close() // This will close the channel writer
	<-done       // Wait for the processor to finish

	fmt.Println("Example completed successfully")
}
