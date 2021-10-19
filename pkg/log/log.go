package log

import (
	"net"
	"os"

	logrustash "github.com/bshuster-repo/logrus-logstash-hook"
	"github.com/openzipkin/zipkin-go"
	log "github.com/sirupsen/logrus"
)

func SetLoglevel(level string) {
	switch level {
	case "info":
		log.SetLevel(log.InfoLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "panic":
		log.SetLevel(log.PanicLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	default:
		log.SetLevel(log.ErrorLevel)
	}
}
func LogStashRegister() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.ErrorLevel)
	connLogstash, err := net.Dial("tcp", os.Getenv("LOGSTASH_IP")+`:`+os.Getenv("LOGSTASH_PORT"))
	Fatal(err)
	hook, err := logrustash.NewHookWithConn(connLogstash, os.Getenv("APPLICATION_NAME"))
	Fatal(err)
	log.AddHook(hook)
}
func Fatal(err error) {
	if err != nil {
		log.WithFields(log.Fields{
			"Application name": os.Getenv("APPLICATION_NAME"),
		}).Fatal(err)
	}
}
func Info(message string) {
	log.WithFields(log.Fields{
		"Application name": os.Getenv("APPLICATION_NAME"),
	}).Info(message)
}
func Warn(message string) {
	log.WithFields(log.Fields{
		"Application name": os.Getenv("APPLICATION_NAME"),
	}).Warn(message)
}
func Error(err error) {
	if err != nil {
		log.WithFields(log.Fields{
			"Application name": os.Getenv("APPLICATION_NAME"),
		}).Error(err)
	}
}
func InfoWithTraceID(message string, Span zipkin.Span) {
	var logFeilds = log.Fields{
		"application name": os.Getenv("APPLICATION_NAME"),
	}
	if Span.Context().ParentID == nil {
		logFeilds["TraceID"] = Span.Context().TraceID.String()
	} else {
		logFeilds["TraceID"] = Span.Context().ParentID.String()
	}

	log.WithFields(logFeilds).Info(message)
}
func WarnWithTraceID(message string, Span zipkin.Span) {
	var logFeilds = log.Fields{
		"application name": os.Getenv("APPLICATION_NAME"),
	}
	if Span.Context().ParentID == nil {
		logFeilds["TraceID"] = Span.Context().TraceID.String()
	} else {
		logFeilds["TraceID"] = Span.Context().ParentID.String()
	}
	log.WithFields(logFeilds).Warn(message)
}
func FatalWithTraceID(Error string, Span zipkin.Span) {
	var logFeilds = log.Fields{
		"application name": os.Getenv("APPLICATION_NAME"),
	}
	if Span.Context().ParentID == nil {
		logFeilds["TraceID"] = Span.Context().TraceID.String()
	} else {
		logFeilds["TraceID"] = Span.Context().ParentID.String()
	}
	if Error != "" {
		log.WithFields(logFeilds).Fatal(Error)
	}
}
func ErrorWithTraceID(error error, span zipkin.Span) {
	var logFeilds = log.Fields{
		"application name": os.Getenv("APPLICATION_NAME"),
	}
	if span != nil {
		if span.Context().ParentID == nil {
			logFeilds["TraceID"] = span.Context().TraceID.String()
		} else {
			logFeilds["TraceID"] = span.Context().ParentID.String()
		}
	}
	log.WithFields(logFeilds).Error(error)
}
func PanicWithTraceID(message string, Span zipkin.Span) {
	var logFeilds = log.Fields{
		"application name": os.Getenv("APPLICATION_NAME"),
	}
	if Span.Context().ParentID == nil {
		logFeilds["TraceID"] = Span.Context().TraceID.String()
	} else {
		logFeilds["TraceID"] = Span.Context().ParentID.String()
	}
	log.WithFields(logFeilds).Panic(message)
}
