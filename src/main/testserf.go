package main

import (
	"fmt"
	"github.com/mitchellh/cli"
	"os"
)

func main() {
	args := os.Args[1:]
	c := cli.NewCLI("dkron", "1")
	c.Args = args
	c.HelpFunc = cli.BasicHelpFunc("dkron")

	ui := &cli.BasicUi{Writer: os.Stdout}
	c.Commands = map[string]cli.CommandFactory{
		"agent": func() (cli.Command, error) {
			return &AgentCommand{
				Ui:               ui,
				Version:          "1",
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