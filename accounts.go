package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/pr8kerl/organizer/aws"
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
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	org, err := aws.NewOrganization()
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
	AccountName  string
	AccountEmail string
	Ui           cli.Ui
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
	cmdFlags.StringVar(&c.AccountEmail, "email", "", "the account email address to use")
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()); os.Exit(1) }
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	if len(c.AccountName) == 0 {
		c.Ui.Error("error: missing create account --name parameter.")
		cmdFlags.Usage()
		return 1
	}
	if len(c.AccountEmail) == 0 {
		c.Ui.Error("error: missing create account --email parameter.")
		cmdFlags.Usage()
		return 1
	}

	fmt.Printf("create account: %s, %s\n", c.AccountName, c.AccountEmail)

	org, err := aws.NewOrganization()
	if err != nil {
		fmt.Printf("error: could not initialize organization: %s\n", err)
		return 1
	}

	status, err := org.CreateAccount(c.AccountName, c.AccountEmail)
	if err != nil {
		fmt.Printf("error: could not create account: %s\n", err)
		return 1
	}
	account, err := org.WaitForAccountStatus(status)
	if err != nil {
		fmt.Printf("error: could not get account status: %s\n", err)
		return 1
	}

	fmt.Printf("created account %s successfully id: %s\n", c.AccountName, *account.AccountId)
	fmt.Printf("ACCOUNT_NAME=%s\n", *account.AccountName)
	fmt.Printf("ACCOUNT_ID=%s\n", *account.AccountId)

	return 0
}

func (c *CreateAccountCommand) Help() string {
	helpText := `
usage: organizer create account --name <account alias> --email <account email address>

create an organization aws account.

options:

	-name=<account name>	the account name or account alias used to identify the account
	-email=<email addr>	the account email address

	`
	return strings.TrimSpace(helpText)
}

func (c *CreateAccountCommand) Synopsis() string {
	return "create an account for an organization"
}
