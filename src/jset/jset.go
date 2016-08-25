package main

import (
	"jse"
	"fmt"
	"github.com/jessevdk/go-flags"
	"encoding/json"
)

type Options struct {
	Script string `short:"s" long:"script" description:"script"`
	Params string `short:"p" long:"params" description:"params"`
}

func main() {
	opts := &Options{}
	flags.Parse(opts)
	if opts.Script == "" || opts.Params == "" {
		fmt.Printf(`Usage: jse -s console.log("hello"+name) -p {"name":"jset"}\n`)
		return
	}

	fmt.Printf(opts.Script,opts.Params)
	engine := jse.NewJsEngine()
	var paramsMap map[string]interface{}
	json.Unmarshal([]byte(opts.Params ), &paramsMap)
	engine.Run(opts.Script, paramsMap)
}
