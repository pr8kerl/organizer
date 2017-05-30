package main

import (
	"flag"
	"github.com/mitchellh/cli"
	"os"
	"strings"
)

const (
	RunResultHelp = -18511
)

// List Command
type ListCommand struct {
	Ui cli.Ui
}

func listCmdFactory() (cli.Command, error) {

	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}

	return &ListCommand{
		Ui: ui,
	}, nil
}

func (c *ListCommand) Run(args []string) int {

	cmdFlags := flag.NewFlagSet("list", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()); os.Exit(1) }
	cmdFlags.Parse(args)

	return RunResultHelp
}

func (c *ListCommand) Help() string {
	helpText := `usage: organizer list <subcommand> [<args>]

list organization objects

	`

	return strings.TrimSpace(helpText)
}

func (c *ListCommand) Synopsis() string {
	return "list objects for an organization"
}

// Create Command
type CreateCommand struct {
	Ui cli.Ui
}

func createCmdFactory() (cli.Command, error) {

	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}

	return &CreateCommand{
		Ui: ui,
	}, nil
}

func (c *CreateCommand) Run(args []string) int {

	cmdFlags := flag.NewFlagSet("create", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()); os.Exit(1) }
	cmdFlags.Parse(args)

	return RunResultHelp
}

func (c *CreateCommand) Help() string {
	helpText := `usage: organizer create <subcommand> [<args>]

create organization objects

	`

	return strings.TrimSpace(helpText)
}

func (c *CreateCommand) Synopsis() string {
	return "create objects for an organization"
}
