package slog

import "sync"

type ToChanWriterOptions struct {
	Level LogLevel
	Ch    chan *Log
}

type toChanWriter struct {
	mu    sync.RWMutex
	level LogLevel
	ch    chan *Log
}

func (w *toChanWriter) Level() LogLevel {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.level
}

func WithToChanWriter(opt *ToChanWriterOptions) *toChanWriter {
	if opt == nil {
		// will be skipped
		return nil
	}

	return &toChanWriter{
		level: opt.Level,
		ch:    opt.Ch,
	}
}

func (w *toChanWriter) write(l *Log) error {
	// non blocking write
	select {
	case w.ch <- l:
	default:
	}
	return nil
}

func (w *toChanWriter) Write(l *Log) error {
	if l == nil {
		return nil
	}

	if l.Level < w.Level() {
		return nil
	}

	log := l.Copy()
	w.write(&log)
	return nil
}

func (w *toChanWriter) Close() {
	close(w.ch)
}

var _ Writer = &toChanWriter{}
