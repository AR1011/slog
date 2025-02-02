package slog

import (
	"encoding/json"
	"os"
	"sync"
)

type StdStream int

const (
	StdOut StdStream = iota
	StdErr
)

type toStdStreamWriter struct {
	mu     sync.RWMutex
	level  LogLevel
	format LogFormat
	Stream StdStream
}

type ToStdStreamWriterOptions struct {
	Level  LogLevel
	Format LogFormat
	Stream StdStream
}

func (w *toStdStreamWriter) Level() LogLevel {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.level
}

func WithStdIoWriter(opt *ToStdStreamWriterOptions) *toStdStreamWriter {
	if opt == nil {
		// will be skipped
		return nil
	}

	return &toStdStreamWriter{
		level:  opt.Level,
		format: opt.Format,
		Stream: opt.Stream,
	}
}

func (w *toStdStreamWriter) Write(l *Log) error {
	if l == nil {
		return nil
	}

	if l.Level < w.Level() {
		return nil
	}

	var b []byte
	var err error

	switch w.format {
	case FormatJson:
		b, err = json.Marshal(l)
		if err != nil {
			return err
		}
		b = append(b, '\n')

	case FormatText:
		b = removeAnsiBytes([]byte(l.Str))

	case FormatAnsi:
		b = []byte(l.Str)

	default:
		b = []byte(l.Str)

	}

	switch w.Stream {
	case StdOut:
		_, err := os.Stdout.Write(b)
		return err
	case StdErr:
		_, err := os.Stderr.Write(b)
		return err
	default:
		return nil
	}
}

func (w *toStdStreamWriter) Close() {}

var _ Writer = &toStdStreamWriter{}
