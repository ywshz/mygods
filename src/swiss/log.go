package swiss

import (
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("swiss")

func init() {
	var format = logging.MustStringFormatter(
		`%{time:2006/01/02 15:04:05.999} %{shortfile} %{shortpkg} %{longfunc} %{level:.4s} ▶%{color:reset} %{message}`,
	)
	logging.SetFormatter(format)
	logging.SetLevel(logging.DEBUG, "swiss")

	//log.Debugf("debug %s", "debug")
	//log.Info("info")
	//log.Notice("notice")
	//log.Warning("warning")
	//log.Error("err")
	//log.Critical("crit")
}


//import (
//	"github.com/Sirupsen/logrus"
//)
//
//var log = logrus.NewEntry(logrus.New())
//
//func init(){
//	InitLogger("debug","")
//}
//
//func InitLogger(logLevel string, node string) {
//	formattedLogger := logrus.New()
//	formattedLogger.Formatter = &logrus.TextFormatter{FullTimestamp: true}
//
//	level, err := logrus.ParseLevel(logLevel)
//	if err != nil {
//		logrus.WithError(err).Error("Error parsing log level, using: info")
//		level = logrus.DebugLevel
//	}
//
//	formattedLogger.Level = level
//	log = logrus.NewEntry(formattedLogger)
//}
