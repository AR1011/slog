package slog

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type ToHttpWriterOptions struct {
	Level  LogLevel
	Format LogFormat
	URL    string
	Method string
	APIKey string
}

type ToHttpWriter struct {
	mu     sync.RWMutex
	level  LogLevel
	format LogFormat
	url    string
	method string
	apiKey string

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	logCh  chan *Log
	client *http.Client
}

func validateToHttpWriterOptions(opt *ToHttpWriterOptions) error {
	if opt == nil {
		return fmt.Errorf("options are nil")
	}
	if opt.URL == "" {
		return fmt.Errorf("URL is empty")
	}
	if _, err := url.ParseRequestURI(opt.URL); err != nil {
		return fmt.Errorf("invalid URL: %v", err)
	}

	if opt.Method == "" {
		opt.Method = http.MethodPut
	}
	validMethods := map[string]bool{
		http.MethodGet:     true,
		http.MethodPost:    true,
		http.MethodPut:     true,
		http.MethodDelete:  true,
		http.MethodPatch:   true,
		http.MethodOptions: true,
		http.MethodHead:    true,
	}
	if !validMethods[opt.Method] {
		return fmt.Errorf("invalid HTTP method: %s", opt.Method)
	}

	if opt.Level == 0 {
		opt.Level = InfoLevel
	}
	if opt.Format == "" {
		opt.Format = FormatJson
	}

	return nil
}

func WithToHttpWriter(opt *ToHttpWriterOptions) *ToHttpWriter {
	err := validateToHttpWriterOptions(opt)
	if err != nil {
		fmt.Printf("WithToHttpWriter(): %v\n", err)
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	w := &ToHttpWriter{
		level:  opt.Level,
		format: opt.Format,
		url:    opt.URL,
		method: opt.Method,
		apiKey: opt.APIKey,
		ctx:    ctx,
		cancel: cancel,
		logCh:  make(chan *Log, 100),
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}

	w.wg.Add(1)
	go w.run()
	return w
}

func (w *ToHttpWriter) Level() LogLevel {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.level
}

func (w *ToHttpWriter) Write(l *Log) error {
	if l == nil {
		return nil
	}

	if l.Level < w.Level() {
		return nil
	}

	select {
	case <-w.ctx.Done():
		fmt.Printf("ToHttpWriter.Write(): context canceled, can't write log\n")
	default:
		select {
		case w.logCh <- l:
		default:
			fmt.Printf("ToHttpWriter.Write(): log channel full, dropping log\n")
		}
	}

	return nil
}

func (w *ToHttpWriter) Close() {
	w.cancel()
	w.wg.Wait()
}

func (w *ToHttpWriter) run() {
	defer w.wg.Done()
	for {
		select {
		case <-w.ctx.Done():
			for {
				select {
				case l := <-w.logCh:
					w.sendLog(l)
				default:
					return
				}
			}
		case l := <-w.logCh:
			w.sendLog(l)
		}
	}
}

func (w *ToHttpWriter) sendLog(l *Log) {
	b, err := encodeLog(l, w.format)
	if err != nil {
		fmt.Printf("ToHttpWriter.sendLog(): failed to encode log: %v\n", err)
		return
	}

	req, err := http.NewRequest(w.method, w.url, bytes.NewReader(b))
	if err != nil {
		fmt.Printf("ToHttpWriter.sendLog(): failed to create request: %v\n", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	if w.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+w.apiKey)
	}

	resp, err := w.client.Do(req)
	if err != nil {
		fmt.Printf("ToHttpWriter.sendLog(): failed to send log: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		fmt.Printf("ToHttpWriter.sendLog(): received non-2xx status code: %d\n", resp.StatusCode)
		return
	}
}
