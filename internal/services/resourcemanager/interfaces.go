package resourcemanager

import mqtt "github.com/eclipse/paho.mqtt.golang"

type logger interface {
	Printf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Tracef(format string, args ...interface{})
	Print(args ...interface{})
	Error(args ...interface{})
	Debug(args ...interface{})
	Info(args ...interface{})
	Trace(args ...interface{})
}

type notifier interface {
	Send(msg string)
}

type mqttServer interface {
	Publish(address string, topic string, payload interface{})
	Subscribe(address string, topic string, handler mqtt.MessageHandler)
}
