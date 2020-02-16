package main

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

func SetupDefaultLogger() *logrus.Logger{
	return SetupLogger(nil, nil, logrus.TraceLevel)
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
