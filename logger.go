package slog

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"time"
)

// set default logger
var Slog *SLogger = &SLogger{
	Buffer: NewLogBuffer(1000),
	writers: []Writer{
		&toStdStreamWriter{
			Stream: StdOut,
			format: FormatAnsi,
			level:  InfoLevel,
		},
	},
}

type Log struct {
	Type      string                 `json:"level"`
	Level     LogLevel               `json:"-"`
	TypeColor string                 `json:"-"`
	Time      string                 `json:"time"`
	Timestamp time.Time              `json:"timestamp"`
	Msg       string                 `json:"msg"`
	MsgColor  string                 `json:"-"`
	Args      map[string]interface{} `json:"args"`
	Str       string                 `json:"-"`
}

func (l *Log) Copy() Log {
	newLog := *l
	if l.Args != nil {
		newArgs := make(map[string]interface{}, len(l.Args))
		for k, v := range l.Args {
			newArgs[k] = v
		}
		newLog.Args = newArgs
	}
	return newLog
}

type SLogger struct {
	writers []Writer
	Buffer  *LogBuffer
}

func NewLogger(writers ...Writer) (*SLogger, error) {
	if len(writers) == 0 {
		w := &toStdStreamWriter{
			Stream: StdOut,
			format: FormatAnsi,
			level:  InfoLevel,
		}
		writers = append(writers, w)
	}

	for _, w := range writers {
		if w == nil {
			fmt.Printf("NewLogger(): writer is nil\n")
			return nil, errors.New("writer is nil")
		}
	}

	return &SLogger{
		writers: writers,
		Buffer:  NewLogBuffer(1000),
	}, nil
}

func SetLogger(l *SLogger) {
	if l == nil {
		return
	}
	Slog = l
}

func (s *SLogger) AddWriter(w Writer) {
	if w == nil {
		return
	}
	s.writers = append(s.writers, w)
}

func AddWriter(w Writer) {
	if Slog == nil {
		return
	}
	if w == nil {
		return
	}

	Slog.AddWriter(w)
}

func (s *SLogger) debug(msg string, args ...interface{}) {
	s.writeLog(DebugLevel, "DBUG", ColorYellow, msg, args...)
}

func (s *SLogger) info(msg string, args ...interface{}) {
	s.writeLog(InfoLevel, "INFO", ColorBlue, msg, args...)
}

func (s *SLogger) warn(msg string, args ...interface{}) {
	s.writeLog(WarnLevel, "WARN", ColorOrange, msg, args...)
}

func (s *SLogger) error(msg string, args ...interface{}) {
	s.writeLog(ErrorLevel, "EROR", ColorRed, msg, args...)
}

func (s *SLogger) fatal(msg string, args ...interface{}) {
	s.writeLog(FatalLevel, "FTAL", ColorPurple, msg, args...)
	os.Exit(1)
}

func (s *SLogger) stat(msg string, args ...interface{}) {
	s.writeLog(InfoLevel, "STAT", ColorGreen, msg, args...)
}

func (s *SLogger) panic(msg string, args ...interface{}) {
	s.writeLog(PanicLevel, "PANC", ColorPink, msg, args...)
	panic(msg)
}

func (s *SLogger) debugF(format string, formatArgs []interface{}, args ...interface{}) {
	formattedMsg := toFormatStr(format, formatArgs)
	s.debug(formattedMsg, args...)
}

func (s *SLogger) infoF(format string, formatArgs []interface{}, args ...interface{}) {
	formattedMsg := toFormatStr(format, formatArgs)
	s.info(formattedMsg, args...)
}

func (s *SLogger) warnF(format string, formatArgs []interface{}, args ...interface{}) {
	formattedMsg := toFormatStr(format, formatArgs)
	s.warn(formattedMsg, args...)
}

func (s *SLogger) errorF(format string, formatArgs []interface{}, args ...interface{}) {
	formattedMsg := toFormatStr(format, formatArgs)
	s.error(formattedMsg, args...)
}

func (s *SLogger) fatalF(format string, formatArgs []interface{}, args ...interface{}) {
	formattedMsg := toFormatStr(format, formatArgs)
	s.fatal(formattedMsg, args...)
}

func (s *SLogger) panicF(format string, formatArgs []interface{}, args ...interface{}) {
	formattedMsg := toFormatStr(format, formatArgs)
	s.panic(formattedMsg, args...)
}

func (s *SLogger) writeLog(lvl LogLevel, t string, c string, msg string, args ...interface{}) {
	log := &Log{
		Level:     lvl,
		Type:      t,
		TypeColor: c,
		Time:      GetTime(),
		Timestamp: time.Now(),
		Msg:       msg,
		Args:      toArgsMap(args),
	}
	log.Str = s.toString(log)
	s.write(log)
}

func toArgsMap(args []interface{}) map[string]interface{} {
	m := make(map[string]interface{}, len(args)/2)
	keys := make([]string, 0, len(args)/2)

	for i := 0; i < len(args); i += 2 {
		key, ok := args[i].(string)
		if !ok {
			key = fmt.Sprintf("arg%d", i)
		}
		keys = append(keys, key)

		var value interface{}
		if i+1 < len(args) {
			value = args[i+1]
		} else {
			value = nil
		}
		m[key] = value
	}

	ordered := make(map[string]interface{}, len(keys))
	for _, key := range keys {
		ordered[key] = m[key]
	}

	return ordered
}

func (s *SLogger) write(l *Log) {
	if l == nil {
		return
	}

	s.Buffer.Add(l)
	for _, w := range s.writers {
		if w == nil {
			fmt.Printf("SLogger.write(): writer is nil\n")
			continue
		}

		if l.Level < w.Level() {
			continue
		}

		err := w.Write(l)
		if err != nil {
			fmt.Printf("SLogger.write(): failed to write log entry: %v\n", err)
		}
	}
}

func (s *SLogger) toString(log *Log) string {
	if log == nil {
		return ""
	}

	logStr := FontBold + log.TypeColor + "[" + log.Type + "]" + FontNormal
	logStr += ColorWhite + " [" + log.Time + "]"
	logStr += pad(log.Msg) + ColorWhite

	if len(log.Args) > 0 {
		keys := make([]string, 0, len(log.Args))
		for k := range log.Args {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			logStr += fmt.Sprintf(" %s=%v", k, log.Args[k])
		}
	}

	logStr += "\n"
	return logStr
}

func (s *SLogger) Close() {
	for _, w := range s.writers {
		if w == nil {
			continue
		}
		w.Close()
	}
}

func Close() {
	if Slog == nil {
		return
	}
	Slog.Close()
}

func toFormatStr(f string, args []interface{}) string {
	if len(args) == 0 {
		return f
	}
	return fmt.Sprintf(f, args...)
}

func Info(msg string, args ...interface{}) {
	if Slog == nil {
		return
	}

	Slog.info(msg, args...)
}

func Debug(msg string, args ...interface{}) {
	if Slog == nil {
		return
	}
	Slog.debug(msg, args...)
}

func Warn(msg string, args ...interface{}) {
	if Slog == nil {
		return
	}
	Slog.warn(msg, args...)
}

func Error(msg string, args ...interface{}) {
	if Slog == nil {
		return
	}
	Slog.error(msg, args...)
}

func Fatal(msg string, args ...interface{}) {
	if Slog == nil {
		return
	}
	Slog.fatal(msg, args...)
}

func Stat(msg string, args ...interface{}) {
	if Slog == nil {
		return
	}
	Slog.stat(msg, args...)
}

func Panic(msg string, args ...interface{}) {
	if Slog == nil {
		return
	}
	Slog.panic(msg, args...)
}

func DebugF(format string, formatArgs []interface{}, args ...interface{}) {
	if Slog == nil {
		return
	}
	Slog.debugF(format, formatArgs, args...)
}

func InfoF(format string, formatArgs []interface{}, args ...interface{}) {
	if Slog == nil {
		return
	}
	Slog.infoF(format, formatArgs, args...)
}

func WarnF(format string, formatArgs []interface{}, args ...interface{}) {
	if Slog == nil {
		return
	}
	Slog.warnF(format, formatArgs, args...)
}

func ErrorF(format string, formatArgs []interface{}, args ...interface{}) {
	if Slog == nil {
		return
	}
	Slog.errorF(format, formatArgs, args...)
}

func FatalF(format string, formatArgs []interface{}, args ...interface{}) {
	if Slog == nil {
		return
	}
	Slog.fatalF(format, formatArgs, args...)

}

func PanicF(format string, formatArgs []interface{}, args ...interface{}) {
	if Slog == nil {
		return
	}
	Slog.panicF(format, formatArgs, args...)
}
