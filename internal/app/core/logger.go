package core

import (
	"fmt"
	graylog "github.com/gemnasium/logrus-graylog-hook"
	"github.com/sirupsen/logrus"
	"os"
)

type Logger struct {
	*logrus.Logger
}

var logger *Logger

func GetLogger() *Logger {
	if logger == nil {
		logger = initLogger()
	}

	return logger
}

func initLogger() *Logger {
	config := GetConfig()
	logLevel, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		logLevel = logrus.DebugLevel
	}

	externalLogger := logrus.New()
	externalLogger.Out = os.Stdout
	externalLogger.Level = logLevel
	externalLogger.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "02-01 15:04:05.000",
		PadLevelText:    true,
	})

	if len(config.GraylogHost) > 0 && len(config.GraylogPort) > 0 {
		graylogCredentials := fmt.Sprintf("%s:%s", config.GraylogHost, config.GraylogPort)
		hook := graylog.NewGraylogHook(graylogCredentials, map[string]interface{}{"service": "sentinel_proxy"})
		externalLogger.AddHook(hook)
	}

	return &Logger{externalLogger}
}
