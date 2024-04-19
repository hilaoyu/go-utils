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
func Debug(msg interface{}) {
	if nil == defaultLogger {
		return
	}
	defaultLogger.Debug(msg)
}
func Info(msg interface{}) {
	if nil == defaultLogger {
		return
	}
	defaultLogger.Info(msg)
}
func Warn(msg interface{}) {
	if nil == defaultLogger {
		return
	}
	defaultLogger.Warn(msg)
}
func Error(msg interface{}) {
	if nil == defaultLogger {
		return
	}
	defaultLogger.Error(msg)
}
func Fatal(msg interface{}) {
	if nil == defaultLogger {
		return
	}
	defaultLogger.Fatal(msg)
}
func Panic(msg interface{}) {
	if nil == defaultLogger {
		return
	}
	defaultLogger.Fatal(msg)
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
	err = os.MkdirAll(dir, 0744)
	if err != nil {
		return
	}

	w := &lumberjack.Logger{
		Filename:   path.Join(dir, name),
		MaxBackups: maxBackups, // files
		MaxSize:    maxSize,    // megabytes
		MaxAge:     maxAge,     // days
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
func (l *Logger) Debug(msg interface{}) {
	l.Init(false)
	l.logger.Debug().Msg(fmt.Sprintf("%+v", msg))
}
func (l *Logger) Info(msg interface{}) {
	l.Init(false)
	l.logger.Info().Msg(fmt.Sprintf("%+v", msg))
}
func (l *Logger) Warn(msg interface{}) {
	l.Init(false)
	l.logger.Warn().Msg(fmt.Sprintf("%+v", msg))
}
func (l *Logger) Error(msg interface{}) {
	l.Init(false)
	l.logger.Error().Msg(fmt.Sprintf("%+v", msg))
}
func (l *Logger) Fatal(msg interface{}) {
	l.Init(false)
	l.logger.Fatal().Msg(fmt.Sprintf("%+v", msg))
}
func (l *Logger) Panic(msg interface{}) {
	l.Init(false)
	l.logger.Fatal().Msg(fmt.Sprintf("%+v", msg))
}
