package slog

import (
	er "errors"
	"sync"
	"time"

	"github.com/mattn/go-runewidth"
)

func joinError(m string, err error) error {
	if err == nil {
		return nil
	}
	m = m + ": "
	return er.New(m + err.Error())
}

const (
	ColorDarkBLue   = "\u001b[38;2;0;51;204m"
	ColorDarkGreen  = "\033[38;2;0;153;0m"
	ColorDarkPurple = "\033[38;2;102;0;153m"
	ColorGreen      = "\033[32m"
	ColorBlue       = "\033[34m"
	ColorRed        = "\033[31m"
	ColorOrange     = "\033[33m"
	ColorPurple     = "\033[35m"
	ColorYellow     = "\033[93m"
	ColorPink       = "\033[95m"
	ColorWhite      = "\033[0m"
	FontBold        = "\u001b[1m"
	FontNormal      = "\u001b[0m"
)

const (
	PadWidth = 60
	Truncate = false
)

func pad(s string) string {
	paddedString := " " + s
	width := runewidth.StringWidth(paddedString)

	if width > PadWidth {
		if Truncate {
			return runewidth.Truncate(paddedString, PadWidth, "")
		}
		return paddedString
	}

	for width < PadWidth {
		paddedString += " "
		width = runewidth.StringWidth(paddedString)
	}

	return paddedString
}

func GetTime() string {
	return time.Now().Format("2006-01-02 15:04:05.000")
}

type LogBuffer struct {
	mu  sync.RWMutex
	buf []Log
	max int
}

func NewLogBuffer(max int) *LogBuffer {
	return &LogBuffer{
		buf: make([]Log, 0, max),
		max: max,
	}
}

func (b *LogBuffer) Add(log *Log) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if len(b.buf) >= b.max {
		b.buf = b.buf[1:]
	}

	b.buf = append(b.buf, log.Copy())
}

func (b *LogBuffer) GetLogs(limit int64) []Log {
	b.mu.RLock()
	defer b.mu.RUnlock()

	var logs []Log

	if limit <= 0 || int64(len(b.buf)) < limit {
		logs = b.buf
	} else {
		logs = b.buf[len(b.buf)-int(limit):]
	}

	logsCopy := make([]Log, len(logs))
	for i := range logs {
		logsCopy[i] = logs[i].Copy()
	}

	return logsCopy
}

func (b *LogBuffer) GetLogsInterface(limit int64, format LogFormat) ([]interface{}, int64) {
	logs := b.GetLogs(limit)
	if len(logs) == 0 {
		return []interface{}{}, 0
	}

	logsInterface := make([]interface{}, 0, len(logs))

	for _, log := range logs {
		encoded := EncodeLogToInterface(log, format)
		logsInterface = append(logsInterface, encoded)
	}

	return logsInterface, int64(len(logsInterface))
}

func (b *LogBuffer) Clear() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.buf = make([]Log, 0, b.max)
}
