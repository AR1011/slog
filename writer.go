package slog

import (
	"encoding/json"
	"regexp"
)

type LogFormat string

const (
	FormatJson LogFormat = "json"
	FormatText LogFormat = "text"
	FormatAnsi LogFormat = "ansi"
)

type (
	LogLevelString string
	LogLevel       int
)

const (
	DebugLevelString LogLevelString = "debug"
	InfoLevelString  LogLevelString = "info"
	WarnLevelString  LogLevelString = "warn"
	ErrorLevelString LogLevelString = "error"
	FatalLevelString LogLevelString = "fatal"
	PanicLevelString LogLevelString = "panic"
)

const (
	DebugLevel LogLevel = iota + 1
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	PanicLevel
)

func ToLogLevel(s string) LogLevel {
	switch s {
	case string(DebugLevelString):
		return DebugLevel

	case string(InfoLevelString):
		return InfoLevel

	case string(WarnLevelString):
		return WarnLevel

	case string(ErrorLevelString):
		return ErrorLevel

	case string(FatalLevelString):
		return FatalLevel

	case string(PanicLevelString):
		return PanicLevel

	default:
		return InfoLevel
	}
}

func ToLogFormat(s string) LogFormat {
	switch s {
	case string(FormatJson):
		return FormatJson

	case string(FormatText):
		return FormatText

	case string(FormatAnsi):
		return FormatAnsi

	default:
		return FormatText
	}
}

func IsValidLogFormat(s string) bool {
	formats := []LogFormat{FormatJson, FormatText, FormatAnsi}
	for _, f := range formats {
		if s == string(f) {
			return true
		}
	}
	return false
}

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func removeAnsiBytes(b []byte) []byte {
	return ansiRegex.ReplaceAll(b, []byte{})
}

func RemoveAnsiString(s string) string {
	return ansiRegex.ReplaceAllString(s, "")
}

func encodeLog(l *Log, f LogFormat) ([]byte, error) {
	var b []byte
	var err error

	switch f {

	case FormatJson:
		b, err = json.Marshal(l)
		if err != nil {
			return b, err
		}
		b = append(b, '\n')

	case FormatText:
		b = removeAnsiBytes([]byte(l.Str))

	case FormatAnsi:
		b = []byte(l.Str)

	default:
		b = removeAnsiBytes([]byte(l.Str))
	}

	return b, nil
}

func EncodeLogToInterface(l Log, f LogFormat) interface{} {
	switch f {

	case FormatJson:
		return l

	case FormatAnsi:
		return l.Str

	default:
		return RemoveAnsiString(l.Str)
	}
}

type Writer interface {
	Level() LogLevel
	Write(*Log) error
	Close()
}
