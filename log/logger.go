package log

import (
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
)

type Params struct {
	Level      string
	File       string
	MaxSize    int
	MaxBackups int
	MaxAge     int
}

type Logger struct {
	*Params
	*logrus.Logger
}

func NewLogger(params *Params) (l *Logger) {
	if params.MaxSize == 0 {
		params.MaxSize = 50
	}
	if params.MaxBackups == 0 {
		params.MaxBackups = 3
	}
	if params.MaxAge == 0 {
		params.MaxAge = 30
	}

	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			slash := strings.LastIndex(f.File, "/")
			filename := f.File[slash+1:]
			return "", "[" + filename + ":" + strconv.Itoa(f.Line) + "]"
		},
	})
	logger.SetReportCaller(true)

	l = &Logger{
		Params: params,
		Logger: logger,
	}

	level := logrus.InfoLevel
	if params.Level == "debug" {
		level = logrus.DebugLevel
	} else if params.Level == "info" {
		level = logrus.InfoLevel
	} else if params.Level == "warn" {
		level = logrus.WarnLevel
	} else if params.Level == "error" {
		level = logrus.ErrorLevel
	}
	l.Logger.SetLevel(level)

	if params.File != "" {
		_, err := os.OpenFile(params.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			logrus.Error(err)
			return
		}

		rotateFileHook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
			Filename:   params.File,
			MaxSize:    params.MaxSize,
			MaxBackups: params.MaxBackups,
			MaxAge:     params.MaxAge, //days
			Level:      level,
			Formatter: &logrus.JSONFormatter{
				TimestampFormat: time.RFC3339,
			},
		})
		if err != nil {
			logrus.Error(err)
			return
		}

		l.Logger.AddHook(rotateFileHook)
	}

	return
}

func (l *Logger) DebugWithFields(fields map[string]interface{}, args ...interface{}) {
	l.Logger.WithFields(fields).Debug(args)
}

func (l *Logger) InfoWithFields(fields map[string]interface{}, args ...interface{}) {
	l.Logger.WithFields(fields).Info(args)
}

func (l *Logger) WarnWithFields(fields map[string]interface{}, args ...interface{}) {
	l.Logger.WithFields(fields).Warn(args)
}

func (l *Logger) ErrorWithFields(fields map[string]interface{}, args ...interface{}) {
	l.Logger.WithFields(fields).Error(args)
}

func (l *Logger) FatalWithFields(fields map[string]interface{}, args ...interface{}) {
	l.Logger.WithFields(fields).Fatal(args)
}

func (l *Logger) PanicWithFields(fields map[string]interface{}, args ...interface{}) {
	l.Logger.WithFields(fields).Panic(args)
}

func Fatal(args ...interface{}) {
	logrus.Fatal(args)
}
