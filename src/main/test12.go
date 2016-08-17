package main

import (
	"github.com/astaxie/beego/logs"
)

func main() {
	var log = logs.NewLogger(10000)

	log.EnableFuncCallDepth(true)
	//log.SetLogger("file", `{"filename":"jse.log"}`)
	log.SetLogger("console", "")
	log.SetLevel(logs.LevelDebug)

	log.Trace("trace")
	log.Info("info")
	log.Warn("warning")
	log.Debug("debug")
	log.Critical("critical")
}