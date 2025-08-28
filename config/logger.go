package config

import (
	log "github.com/sirupsen/logrus"
)

func SetupLogger() *log.Logger {
	logger := log.New()
	logger.SetReportCaller(true)
	logger.SetFormatter(&log.JSONFormatter{})
	return logger
}