package logger

import (
	"fmt"
	"memberserver/pkg/slack"

	log "github.com/sirupsen/logrus"
)

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
