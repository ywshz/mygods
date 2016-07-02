package main

import (
	"fmt"
	"os"
	"github.com/mitchellh/cli"
	"./swiss"
)

const (
	VERSION = "0.1"
)

func main() {
	// Get the command line args. We shortcut "--version" and "-v" to
	// just show the version.
	args := os.Args[1:]
	for _, arg := range args {
		if arg == "-v" || arg == "--version" {
			newArgs := make([]string, len(args)+1)
			newArgs[0] = "version"
			copy(newArgs[1:], args)
			args = newArgs
			break
		}
	}

	c := cli.NewCLI("dkron", VERSION)
	c.Args = args
	c.HelpFunc = cli.BasicHelpFunc("dkron")

	ui := &cli.BasicUi{Writer: os.Stdout}
	c.Commands = map[string]cli.CommandFactory{
		"server": func() (cli.Command, error) {
			return &swiss.ServerCommand{
				Ui:      ui,
				Version : VERSION,
			}, nil
		},
		"worker": func() (cli.Command, error) {
			return &swiss.Worker{
			}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err.Error())
		os.Exit(1)
	}

	os.Exit(exitStatus)
}
