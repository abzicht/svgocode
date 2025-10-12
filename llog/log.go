package llog

import (
	"fmt"
	"io"
	"os"
)

const (
	LPanic = -1
	LFatal = 0
	LError = 1
	LWarn  = 2
	LInfo  = 3
	LDebug = 4
)

type Log struct {
	level  int
	prefix string
	writer io.Writer
}

var logDefault Log = Log{level: LInfo, prefix: "SVGOCODE", writer: os.Stderr}

var log *Log = &logDefault

func SetLogger(l *Log) {
	log = l
}

func SetLevel(level int) {
	if level > LDebug {
		level = LDebug
	}
	if level < LFatal {
		level = LFatal
	}
	log.level = level
}

func SetWriter(w io.Writer) {
	log.writer = w
}

func level2str(level int) string {
	switch level {
	case LPanic:
		return "Panic"
	case LFatal:
		return "Fatal"
	case LError:
		return "Error"
	case LWarn:
		return "Warn"
	case LInfo:
		return "Info"
	case LDebug:
		return "Debug"
	}
	return ""
}

func (l *Log) tag(level int) string {
	return fmt.Sprintf("%s[%s]: ", l.prefix, level2str(level))
}

// Return log string
func (l *Log) logs(level int, v ...any) string {
	return l.tag(level) + fmt.Sprint(v...)
}

// Return log (formatted string)
func (l *Log) logsf(level int, format string, v ...any) string {
	return l.tag(level) + fmt.Sprintf(format, v...)
}

// Log string to file, if level requirement is met
func (l *Log) log(level int, v ...any) {
	if l.level >= level {
		fmt.Fprint(l.writer, l.logs(level, v...))
	}
}

// Log (formatted string) to file, if level requirement is met
func (l *Log) logf(level int, format string, v ...any) {
	if l.level >= level {
		fmt.Fprint(l.writer, l.logsf(level, format, v...))
	}
}

func (l *Log) Panic(v ...any) {
	panic(l.logs(LPanic, v...))
}

func (l *Log) Fatal(v ...any) {
	l.log(LFatal, v...)
	os.Exit(1)
}

func (l *Log) Error(v ...any) {
	l.log(LError, v...)
}

func (l *Log) Warn(v ...any) {
	l.log(LWarn, v...)
}

func (l *Log) Info(v ...any) {
	l.log(LInfo, v...)
}

func (l *Log) Debug(v ...any) {
	l.log(LDebug, v...)
}

func (l *Log) Panicf(format string, v ...any) {
	panic(l.logsf(LPanic, format, v...))
}

func (l *Log) Fatalf(format string, v ...any) {
	l.logf(LFatal, format, v...)
	os.Exit(1)
}

func (l *Log) Errorf(format string, v ...any) {
	l.logf(LError, format, v...)
}

func (l *Log) Warnf(format string, v ...any) {
	l.logf(LWarn, format, v...)
}

func (l *Log) Infof(format string, v ...any) {
	l.logf(LInfo, format, v...)
}

func (l *Log) Debugf(format string, v ...any) {
	l.logf(LDebug, format, v...)
}

func Panic(v ...any) {
	log.Panic(v...)
}

func Fatal(v ...any) {
	log.Fatal(v...)
}

func Error(v ...any) {
	log.Error(v...)
}

func Warn(v ...any) {
	log.Warn(v...)
}

func Info(v ...any) {
	log.Info(v...)
}

func Debug(v ...any) {
	log.Debug(v...)
}

func Panicf(format string, v ...any) {
	log.Panicf(format, v...)
}

func Fatalf(format string, v ...any) {
	log.Fatalf(format, v...)
}

func Errorf(format string, v ...any) {
	log.Errorf(format, v...)
}

func Warnf(format string, v ...any) {
	log.Warnf(format, v...)
}

func Infof(format string, v ...any) {
	log.Infof(format, v...)
}

func Debugf(format string, v ...any) {
	log.Debugf(format, v...)
}
