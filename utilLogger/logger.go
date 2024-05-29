package utilLogger

import (
	"fmt"
	"github.com/hilaoyu/go-utils/utilTime"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path"
)

type Logger struct {
	logger  *zerolog.Logger
	writers []io.Writer

	timeFormat string
}

var defaultLogger *Logger

func SetDefaultLogger(logger *Logger) {
	defaultLogger = logger
}

func Trace(msg interface{}) {
	if nil == defaultLogger {
		return
	}
	defaultLogger.Trace(msg)
}
func TraceF(format string, a ...any) {
	if nil == defaultLogger {
		return
	}
	defaultLogger.TraceF(format, a...)
}
func Debug(msg interface{}) {
	if nil == defaultLogger {
		return
	}
	defaultLogger.Debug(msg)
}
func DebugF(format string, a ...any) {
	if nil == defaultLogger {
		return
	}
	defaultLogger.DebugF(format, a...)
}
func Info(msg interface{}) {
	if nil == defaultLogger {
		return
	}
	defaultLogger.Info(msg)
}
func InfoF(format string, a ...any) {
	if nil == defaultLogger {
		return
	}
	defaultLogger.InfoF(format, a...)
}
func Warn(msg interface{}) {
	if nil == defaultLogger {
		return
	}
	defaultLogger.Warn(msg)
}
func WarnF(format string, a ...any) {
	if nil == defaultLogger {
		return
	}
	defaultLogger.WarnF(format, a...)
}
func Error(msg interface{}) {
	if nil == defaultLogger {
		return
	}
	defaultLogger.Error(msg)
}
func ErrorF(format string, a ...any) {
	if nil == defaultLogger {
		return
	}
	defaultLogger.ErrorF(format, a...)
}
func Fatal(msg interface{}) {
	if nil == defaultLogger {
		return
	}
	defaultLogger.Fatal(msg)
}
func FatalF(format string, a ...any) {
	if nil == defaultLogger {
		return
	}
	defaultLogger.FatalF(format, a...)
}
func Panic(msg interface{}) {
	if nil == defaultLogger {
		return
	}
	defaultLogger.Panic(msg)
}
func PanicF(format string, a ...any) {
	if nil == defaultLogger {
		return
	}
	defaultLogger.PanicF(format, a...)
}

func NewLogger() *Logger {

	logger := &Logger{}

	return logger
}

func (l *Logger) SetTimeFormat(format string) (err error) {
	l.timeFormat = format
	return
}
func (l *Logger) AddConsoleWriter() (err error) {
	l.writers = append(l.writers, zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.TimeFormat = l.timeFormat
	}))
	return
}

func (l *Logger) AddFileWriter(dest string) (err error) {
	w, err := os.OpenFile(dest, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	l.writers = append(l.writers, w)
	return
}
func (l *Logger) AddFileRotationWriter(dir string, name string, maxBackups int, maxSize int, maxAge int) (err error) {
	w, err := NewFileRotationWriter(dir, name, maxBackups, maxSize, maxAge)
	if err != nil {
		return
	}
	l.writers = append(l.writers, w)
	return
}

func (l *Logger) Init(force bool) *Logger {
	if "" == l.timeFormat {
		l.timeFormat = utilTime.GetTimeFormat()
	}
	zerolog.TimeFieldFormat = l.timeFormat
	if force || nil == l.logger {
		mw := zerolog.MultiLevelWriter(l.writers...)
		logger := zerolog.New(mw).With().
			Timestamp().
			Logger()

		l.logger = &logger
	}

	return l
}

func (l *Logger) Trace(msg interface{}) {
	l.Init(false)
	l.logger.Trace().Msg(fmt.Sprintf("%+v", msg))
}
func (l *Logger) TraceF(format string, a ...any) {
	l.Trace(fmt.Sprintf(format, a...))
}
func (l *Logger) Debug(msg interface{}) {
	l.Init(false)
	l.logger.Debug().Msg(fmt.Sprintf("%+v", msg))
}
func (l *Logger) DebugF(format string, a ...any) {
	l.Debug(fmt.Sprintf(format, a...))
}
func (l *Logger) Info(msg interface{}) {
	l.Init(false)
	l.logger.Info().Msg(fmt.Sprintf("%+v", msg))
}
func (l *Logger) InfoF(format string, a ...any) {
	l.Info(fmt.Sprintf(format, a...))
}
func (l *Logger) Warn(msg interface{}) {
	l.Init(false)
	l.logger.Warn().Msg(fmt.Sprintf("%+v", msg))
}
func (l *Logger) WarnF(format string, a ...any) {
	l.Warn(fmt.Sprintf(format, a...))
}
func (l *Logger) Error(msg interface{}) {
	l.Init(false)
	l.logger.Error().Msg(fmt.Sprintf("%+v", msg))
}
func (l *Logger) ErrorF(format string, a ...any) {
	l.Error(fmt.Sprintf(format, a...))
}
func (l *Logger) Fatal(msg interface{}) {
	l.Init(false)
	l.logger.Fatal().Msg(fmt.Sprintf("%+v", msg))
}
func (l *Logger) FatalF(format string, a ...any) {
	l.Fatal(fmt.Sprintf(format, a...))
}
func (l *Logger) Panic(msg interface{}) {
	l.Init(false)
	l.logger.Panic().Msg(fmt.Sprintf("%+v", msg))
}
func (l *Logger) PanicF(format string, a ...any) {
	l.Panic(fmt.Sprintf(format, a...))
}

func NewFileRotationWriter(dir string, name string, maxBackups int, maxSize int, maxAge int) (w io.Writer, err error) {
	err = os.MkdirAll(dir, 0744)
	if err != nil {
		return
	}
	w = &lumberjack.Logger{
		Filename:   path.Join(dir, name),
		MaxBackups: maxBackups, // files
		MaxSize:    maxSize,    // megabytes
		MaxAge:     maxAge,     // days
	}

	return
}
