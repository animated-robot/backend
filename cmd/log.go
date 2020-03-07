package main

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"time"
)

func SetupDefaultLogger(level string, logPath string) *logrus.Logger{
	logFile := getLogFilePath(logPath)
	return SetupLogger(nil, nil, ParseLogLevel(level), logFile)
}

func getLogFilePath(logPath string) string {
	if logPath == "" {
		cw, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		logPath = cw
	}

	timestamp := time.Now().Format(time.RFC3339)
	logName := "log_" + timestamp + ".txt"
	return path.Join(logPath, logName)
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

func SetupLogger(output io.Writer, format logrus.Formatter, level logrus.Level, logFile string) *logrus.Logger{
	logger := logrus.New()

	if output == nil {
		logFile, err := os.OpenFile(logFile, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		mw := io.MultiWriter(os.Stdout, logFile)
		output = mw
	}
	logger.SetOutput(output)

	if format == nil {
		format = &logrus.JSONFormatter{}
	}
	logger.SetFormatter(format)
	logger.SetLevel(level)

	logger.WithFields(logrus.Fields{
		"file": logFile,
	}).Info("SetupLogger: Logger Created")

	return logger
}
