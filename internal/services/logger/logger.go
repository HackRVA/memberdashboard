package logger

import (
	"fmt"
	"memberserver/pkg/slack"

	log "github.com/sirupsen/logrus"
)

type Logger struct {
	level log.Level
}

func New() *Logger {
	return &Logger{}
}

func (l *Logger) SetLevel(level log.Level) {
	l.level = level
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	Errorf(format, args...)
}
func (l *Logger) Debugf(format string, args ...interface{}) {
	Debugf(format, args...)
}
func (l *Logger) Infof(format string, args ...interface{}) {
	Infof(format, args...)
}
func (l *Logger) Fatalf(format string, args ...interface{}) {
	Fatalf(format, args...)
}
func (l *Logger) Tracef(format string, args ...interface{}) {
	Tracef(format, args...)
}

func init() {
	log.SetLevel(log.DebugLevel)
}

func Errorf(format string, args ...interface{}) {
	go slack.Send("[member-server-error]: " + fmt.Sprintf(format, args...))
	log.Errorf(format, args...)
}

var Info = log.Info
var Debug = log.Debug

var Error = log.Error
var Infof = log.Infof

var Debugf = log.Debugf

var Fatal = log.Fatal
var Fatalf = log.Fatalf

var Trace = log.Trace
var Tracef = log.Tracef
