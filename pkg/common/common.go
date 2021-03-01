package common

import (
	log "github.com/sirupsen/logrus"
)

// InitLogger sets global log level and output format
func InitLogger() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(
		&LogFormatter{TimestampFormat: "2006-01-02 15:04:05"})
	log.SetReportCaller(true)
}
