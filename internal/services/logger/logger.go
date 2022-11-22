package logger

import (
	"fmt"

	"github.com/HackRVA/memberserver/pkg/slack"

	log "github.com/sirupsen/logrus"
)

func New() *Logger {
	return &Logger{}
}

func init() {
	log.SetLevel(log.DebugLevel)
}

func Errorf(format string, args ...interface{}) {
	go slack.Send("[member-server-error]: " + fmt.Sprintf(format, args...))
	log.Errorf(format, args...)
}

var (
	Info  = log.Info
	Debug = log.Debug

	Error = log.Error
	Infof = log.Infof

	Debugf = log.Debugf

	Fatal  = log.Fatal
	Fatalf = log.Fatalf

	Trace  = log.Trace
	Tracef = log.Tracef

	Print   = log.Print
	Printf  = log.Printf
	Println = log.Println
)

type Logger struct {
	level log.Level
}

func (l *Logger) SetLevel(level log.Level) {
	l.level = level
}

func (l *Logger) Println(args ...interface{}) {
	Println(args...)
}
func (l *Logger) Printf(format string, args ...interface{}) {
	Printf(format, args...)
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
func (l *Logger) Print(args ...interface{}) {
	Print(args...)
}
func (l *Logger) Error(args ...interface{}) {
	Error(args...)
}
func (l *Logger) Debug(args ...interface{}) {
	Debug(args...)
}
func (l *Logger) Info(args ...interface{}) {
	Info(args...)
}
func (l *Logger) Fatal(args ...interface{}) {
	Fatal(args...)
}
func (l *Logger) Trace(args ...interface{}) {
	Trace(args...)
}
