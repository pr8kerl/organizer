package main

import (
	"flag"
	"fmt"
	"github.com/mitchellh/cli"
	"os"
	"strings"
)

// List Account
type ListAccountsCommand struct {
	All bool
	Ui  cli.Ui
}

func listAccountsCmdFactory() (cli.Command, error) {

	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}

	return &ListAccountsCommand{
		All: false,
		Ui:  ui,
	}, nil
}

func (c *ListAccountsCommand) Run(args []string) int {

	cmdFlags := flag.NewFlagSet("list accounts", flag.ContinueOnError)
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

func (c *ListAccountsCommand) Help() string {
	helpText := `usage: organizer list accounts [<args>]

List organization accounts

Options:
	    -all		show all accounts regardless of state. default is to show only active accounts.
	`
	return strings.TrimSpace(helpText)
}

func (c *ListAccountsCommand) Synopsis() string {
	return "list all accounts for an organization"
}

// Create Account
type CreateAccountCommand struct {
  AccountName string
  AccountEmail string
	Ui cli.Ui
}

func createAccountCmdFactory() (cli.Command, error) {

	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}

	return &CreateAccountCommand{
		Ui: ui,
	}, nil
}

func (c *CreateAccountCommand) Run(args []string) int {

	cmdFlags := flag.NewFlagSet("create account", flag.ContinueOnError)
	cmdFlags.StringVar(&c.AccountName, "name", "", "the account name to use")
	cmdFlags.StringVar(&c.AccuntEmail, "email", "", "the account email address to use")
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()); os.Exit(1) }
	cmdFlags.Parse(args)

	org, err := newOrganization()
	if err != nil {
		fmt.Printf("error: could not initialize organization: %s\n", err)
		return 1
	}

/*
	err = org.CreateAccount(accountname)
	if err != nil {
		fmt.Printf("error: could not create account: %s\n", err)
		return 1
	}
*/

	return 0
}

func (c *CreateAccountCommand) Help() string {
	helpText := `usage: organizer create accounts [<args>]

Create an organization account

	`
	return strings.TrimSpace(helpText)
}

func (c *CreateAccountCommand) Synopsis() string {
	return "create an account for an organization"
}
