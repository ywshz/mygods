package main

import (
	"jse"
	"fmt"
	"github.com/jessevdk/go-flags"
)

type Options struct {
	Port string `short:"p" long:"port" description:"Http Server port"`
}

func main() {
	opts := &Options{}
	flags.Parse(opts)
	if opts.Port == "" {
		fmt.Printf("Usage: jse -p 8888\n")
		return
	}
	fmt.Printf("Port set to %v\n", opts.Port)

	s := jse.NewServer()

	s.Start(opts.Port)
}
