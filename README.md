# slog

A simple, flexible, and feature-rich structured logging library for Go.

## Features

- Multiple log levels (Debug, Info, Warn, Error, Fatal, Panic, Stat)
- Built-in default stdout logger with ANSI colors
- Support for multiple writers (stdout, file, HTTP, and custom channels)
- ANSI color formatting for terminal output
- JSON support for structured logging
- File rotation support
- Buffer system for log management
- Flexible argument handling with both standard and formatted logging
- HTTP logging with API key support
- Channel-based logging for custom log processing
- Concurrent-safe logging

## Installation

```bash
go get github.com/AR1011/slog
```

## Basic Usage

The library comes with a default stdout logger, so you can start logging immediately:

```go
package main

import "github.com/AR1011/slog"

func main() {
    // Basic logging
    slog.Debug("Debug message")
    slog.Info("Info message")
    slog.Warn("Warning message")
    slog.Error("Error message")
    
    // Logging with context
    slog.Info("User logged in", "user_id", 123, "ip", "192.168.1.1")
    
    // Format string logging
    slog.InfoF("User %s logged in from %s", []interface{}{"john", "192.168.1.1"})
    
    // Don't forget to close the logger when you're done
    defer slog.Close()
}
```

## Log Levels

The library supports the following log levels:
- `Debug`: Detailed information for debugging
- `Info`: General information about program execution
- `Warn`: Warning messages for potentially harmful situations
- `Error`: Error messages for serious problems
- `Fatal`: Critical errors that terminate the program
- `Panic`: Critical errors that trigger a panic
- `Stat`: Statistical or metric information

## Writers

### Stdout Writer

The stdout writer is the default writer, but you can configure it explicitly:

```go
logger, _ := slog.NewLogger(slog.WithStdIoWriter(&slog.ToStdStreamWriterOptions{
    Stream: slog.StdOut,  // or slog.StdErr
    Format: slog.FormatAnsi,
    Level:  slog.DebugLevel,
}))
```

### File Writer

Write logs to a file with optional rotation:

```go
slog.AddWriter(slog.WithToFileWriter(&slog.ToFileWriterOptions{
    FileName:   "app.log",
    Format:     slog.FormatJson,
    Level:      slog.InfoLevel,
    RotateSize: 1024 * 1024 * 100, // 100MB
}))
```

### HTTP Writer

Send logs to a remote endpoint:

```go
slog.AddWriter(slog.WithToHttpWriter(&slog.ToHttpWriterOptions{
    URL:    "http://localhost:8080/logs",
    Method: "POST",
    Format: slog.FormatJson,
    Level:  slog.InfoLevel,
    APIKey: "your-api-key-here",
}))
```

### Channel Writer

Process logs using a custom channel handler:

```go
logChan := make(chan *slog.Log, 100)
logger, _ := slog.NewLogger(slog.WithToChanWriter(&slog.ToChanWriterOptions{
    Ch:    logChan,
    Level: slog.InfoLevel,
}))

// Process logs in a separate goroutine
go func() {
    for log := range logChan {
        // Custom log processing
        fmt.Printf("Custom handler: [%s] %s\n", log.Type, log.Msg)
    }
}()
```

## Multiple Writers

You can combine multiple writers to send logs to different destinations:

```go
logger, _ := slog.NewLogger(
    // Console output with ANSI colors
    slog.WithStdIoWriter(&slog.ToStdStreamWriterOptions{
        Stream: slog.StdOut,
        Format: slog.FormatAnsi,
        Level:  slog.DebugLevel,
    }),
    
    // JSON file output
    slog.WithToFileWriter(&slog.ToFileWriterOptions{
        FileName: "app.log",
        Format:  slog.FormatJson,
        Level:   slog.InfoLevel,
    }),
)

slog.SetLogger(logger)
```

## Formatting Options

The library supports two main formatting options:
- `FormatAnsi`: Colored output for terminal (default for stdout)
- `FormatJson`: JSON structured logging (ideal for file and HTTP writers)

## Examples

For more detailed examples, check out the `examples` directory:
- `examples/stdout/`: Basic logging with stdout
- `examples/file/`: File logging with rotation
- `examples/http/`: Remote logging via HTTP
- `examples/channel/`: Custom log processing with channels
- `examples/multiple/`: Using multiple writers
- `examples/advanced/`: Advanced usage patterns

## License

This project is licensed under the MIT License - see the LICENSE file for details. 