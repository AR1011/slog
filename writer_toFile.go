package slog

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const writeCheckInterval int = 200

type ToFileWriter struct {
	opt    ToFileWriterOptions
	logCh  chan []byte
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	f      *os.File
}

func (w *ToFileWriter) Format() LogFormat {
	return w.opt.Format
}

func (w *ToFileWriter) Level() LogLevel {
	return w.opt.Level
}

func (w *ToFileWriter) FileName() string {
	return w.opt.FileName
}

type ToFileWriterOptions struct {
	FileName   string
	Format     LogFormat
	Level      LogLevel
	RotateSize int64 // kB
}

func (o *ToFileWriterOptions) toDefaultIfEmpty() {
	if o.FileName == "" {
		o.FileName = "logs/xllm.log"
	}

	if o.Format == "" {
		o.Format = FormatJson
	}

	if o.Level == 0 {
		o.Level = InfoLevel
	}

	if o.RotateSize == 0 {
		o.RotateSize = 100 * 1024 // kB
	}
}

func WithToFileWriter(opt *ToFileWriterOptions) *ToFileWriter {
	if opt == nil {
		return nil
	}

	opt.toDefaultIfEmpty()
	ctx, cancel := context.WithCancel(context.Background())
	w := &ToFileWriter{
		opt:    *opt,
		logCh:  make(chan []byte, 100),
		ctx:    ctx,
		cancel: cancel,
	}
	w.wg.Add(1)
	go w.run()
	return w
}

func (w *ToFileWriter) run() {
	defer w.wg.Done()

	if err := w.openFile(); err != nil {
		fmt.Printf("ToFileWriter.run(): failed to open file: %v\n", err)
		return
	}

	var writeCounter int = 0

	for {
		select {
		case <-w.ctx.Done():
			for {
				select {
				case logEntry := <-w.logCh:
					w.writeLog(logEntry)
					writeCounter++
					if writeCounter >= writeCheckInterval {
						w.rotate()
						writeCounter = 0
					}
				default:
					w.closeFile()
					return
				}
			}

		case logEntry := <-w.logCh:
			w.writeLog(logEntry)
			writeCounter++
			if writeCounter >= writeCheckInterval {
				w.rotate()
				writeCounter = 0
			}
		}
	}
}

func (w *ToFileWriter) openFile() error {
	dir := filepath.Dir(w.FileName())
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Printf("ToFileWriter.openFile(): creating directory %s\n", dir)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return joinError("ToFileWriter.openFile(): failed to create directory", err)
		}
	}

	f, err := os.OpenFile(w.FileName(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return joinError("ToFileWriter.openFile(): failed to open file", err)
	}
	w.f = f
	return nil
}

func (w *ToFileWriter) rotate() {
	if w.f == nil {
		if err := w.openFile(); err != nil {
			fmt.Printf("ToFileWriter.rotate(): failed to open log file: %v\n", err)
		}
		return
	}

	stat, err := w.f.Stat()
	if err != nil {
		fmt.Printf("ToFileWriter.rotate(): failed to get file stats: %v\n", err)
		return
	}

	if stat.Size() < w.opt.RotateSize*1024 {
		return
	}

	if err := w.f.Close(); err != nil {
		fmt.Printf("ToFileWriter.rotate(): failed to close log file: %v\n", err)
		return
	}
	w.f = nil

	timestamp := time.Now().Format("20060102150405")
	newFileName := fmt.Sprintf("%s.%s.log", w.FileName(), timestamp)
	if err := os.Rename(w.FileName(), newFileName); err != nil {
		fmt.Printf("ToFileWriter.rotate(): failed to rename log file: %v\n", err)
		return
	}

	if err := w.openFile(); err != nil {
		fmt.Printf("ToFileWriter.rotate(): failed to open log file: %v\n", err)
	}

	fmt.Printf("ToFileWriter.rotate(): rotated log file [filename=%s prevFileSize=%dkB]\n", newFileName, stat.Size()/1024)
}

func (w *ToFileWriter) writeLog(b []byte) {
	if w.f == nil {
		if err := w.openFile(); err != nil {
			fmt.Printf("ToFileWriter.writeLog(): failed to open file: %v\n", err)
			return
		}
	}

	if _, err := w.f.Write(b); err != nil {
		fmt.Printf("ToFileWriter.writeLog(): failed to write to file: %v\n", err)
	}
}

func (w *ToFileWriter) Write(l *Log) error {
	if l == nil {
		return nil
	}

	b, err := encodeLog(l, w.Format())
	if err != nil {
		return joinError("ToFileWriter.Write(): failed to encode log", err)
	}

	select {
	case <-w.ctx.Done():
		fmt.Printf("ToFileWriter.Write(): failed to write log entry to file: context canceled\n")
	case w.logCh <- b:
	default:
		fmt.Printf("ToFileWriter.Write(): failed to write log entry to file: buffer full\n")
	}

	return nil
}

func (w *ToFileWriter) Close() {
	w.cancel()
	w.wg.Wait()
}

func (w *ToFileWriter) closeFile() {
	if w.f != nil {
		w.f.Close()
		w.f = nil
	}
}

var _ Writer = &ToFileWriter{}
