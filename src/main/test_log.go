package main

import (
	"github.com/op/go-logging"
)


func main() {
	var log = logging.MustGetLogger("main")
	var format = logging.MustStringFormatter(
		`%{color}%{time:2006-01-02 15:04:05} %{shortfile} %{shortpkg} %{longfunc} %{level:.4s} ▶%{color:reset} %{message}`,
	)
	logging.SetFormatter(format)
	logging.SetLevel(logging.DEBUG,"main")


	log.Debugf("debug %s", "debug")
	log.Info("info")
	log.Notice("notice")
	log.Warning("warning")
	log.Error("err")
	log.Critical("crit")
}