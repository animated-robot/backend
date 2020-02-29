package main

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

func SetupDefaultLogger(level string) *logrus.Logger{
	return SetupLogger(nil, nil, ParseLogLevel(level))
}

func ParseLogLevel(level string) logrus.Level {
	switch level {
	case "TRACE": return logrus.TraceLevel
	case "FATAL": return logrus.FatalLevel
	case "ERROR": return logrus.ErrorLevel
	case "PANIC": return logrus.PanicLevel
	case "INFO" : return logrus.InfoLevel
	default:      return logrus.TraceLevel
	}
}

func SetupLogger(output io.Writer, format logrus.Formatter, level logrus.Level) *logrus.Logger{
	logger := logrus.New()

	if output == nil {
		output = os.Stdout
	}
	logger.SetOutput(output)

	if format == nil {
		format = &logrus.JSONFormatter{}
	}
	logger.SetFormatter(format)

	logger.SetLevel(level)

	return logger
}
