package main

import (
	"flag"
	"fmt"
	"github.com/mitchellh/cli"
	"os"
)

type AccountsCommand struct {
	All bool
	Ui  cli.Ui
}

func accountsCmdFactory() (cli.Command, error) {

	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}

	return &AccountsCommand{
		All: false,
		Ui:  ui,
	}, nil
}

func (c *AccountsCommand) Run(args []string) int {

	cmdFlags := flag.NewFlagSet("accounts", flag.ContinueOnError)
	cmdFlags.BoolVar(&c.All, "all", false, "show all accounts including inactive")
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()); os.Exit(1) }
	cmdFlags.Parse(args)

	org, err := newOrganization()
	if err != nil {
		fmt.Printf("error: could not initialize organization: %s\n", err)
		return 1
	}

	if c.All {
		err = org.PrintAccounts()
		if err != nil {
			fmt.Printf("error: could not list accounts: %s\n", err)
			return 1
		}
	} else {
		err = org.PrintActiveAccounts()
		if err != nil {
			fmt.Printf("error: could not list accounts: %s\n", err)
			return 1
		}
	}

	return 0
}

func (c *AccountsCommand) Help() string {
	return `usage: organizer accounts [<args>]

List organization accounts

Options:
	    -all		show all accounts regardless of state. default is to show only active accounts.
	`
}

func (c *AccountsCommand) Synopsis() string {
	return "list all accounts for an organization"
}
