package main

import (
	"flag"
	"fmt"
	"github.com/mitchellh/cli"
	"os"
)

type TrailsCommand struct {
	AccountId string
	Ui        cli.Ui
}

func trailsCmdFactory() (cli.Command, error) {

	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}

	return &TrailsCommand{
		AccountId: "",
		Ui:        ui,
	}, nil
}

func (c *TrailsCommand) Run(args []string) int {

	cmdFlags := flag.NewFlagSet("trails", flag.ContinueOnError)
	cmdFlags.StringVar(&c.AccountId, "accountid", "", "work on cloudtrails in this account")
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()); os.Exit(1) }
	cmdFlags.Parse(args)

	org, err := newOrganization()
	if err != nil {
		fmt.Printf("error: could not initialize organization: %s\n", err)
		return 1
	}

	if len(c.AccountId) == 0 {
		fmt.Printf("error: trails subcommand requires an accountid options\n")
		return 1

	}

	trails, err := org.GetTrailArnsForAccount(c.AccountId)
	if err != nil {
		fmt.Printf("error: could not list trails for account %s: %s\n", c.AccountId, err)
		return 1
	}
	for _, trail := range trails {
		fmt.Printf("%s,%s\n", c.AccountId, trail)
	}

	return 0
}

func (c *TrailsCommand) Help() string {
	return `usage: organizer trails [<args>]

List cloudtrails within an account

Options:
	    -accountid		specify the account id to work against
	`
}

func (c *TrailsCommand) Synopsis() string {
	return "list all trails for an account"
}
