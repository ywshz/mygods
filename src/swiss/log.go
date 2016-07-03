package swiss

import (
	"github.com/Sirupsen/logrus"
)

var log = logrus.NewEntry(logrus.New())

func InitLogger(logLevel string, node string) {
	formattedLogger := logrus.New()
	formattedLogger.Formatter = &logrus.TextFormatter{FullTimestamp: true}

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logrus.WithError(err).Error("Error parsing log level, using: info")
		level = logrus.InfoLevel
	}

	formattedLogger.Level = level
	log = logrus.NewEntry(formattedLogger).WithField("node", node)
}
