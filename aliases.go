package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/pr8kerl/organizer/aws"
)

// List Aliases
type ListAliasesCommand struct {
	AccountId string
	Report    bool
	Region    string
	Ui        cli.Ui
}

func listAliasesCmdFactory() (cli.Command, error) {

	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}

	return &ListAliasesCommand{
		AccountId: "",
		Report:    false,
		Region:    "",
		Ui:        ui,
	}, nil
}

func (c *ListAliasesCommand) Run(args []string) int {

	cmdFlags := flag.NewFlagSet("list aliases", flag.ContinueOnError)
	cmdFlags.StringVar(&c.AccountId, "accountid", "", "list account aliases for a specific account")
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()); os.Exit(1) }
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	org, err := aws.NewOrganization()
	if err != nil {
		fmt.Printf("error: could not initialize organization: %s\n", err)
		return 1
	}

	if len(c.AccountId) > 0 {

		err = org.PrintAliasesForAccount(c.AccountId)
		if err != nil {
			fmt.Printf("error: could not list aliases for account: %s\n", err)
			return 1
		}

	} else {

		err = org.PrintAliases()
		if err != nil {
			fmt.Printf("error: could not list all account aliases: %s\n", err)
			return 1
		}
	}

	return 0
}

func (c *ListAliasesCommand) Help() string {
	helpText := `usage: organizer list aliases [<args>]

List account aliases for accounts within an organization

Options:
	    -accountid		list all iam users for a specific account only
	`
	return strings.TrimSpace(helpText)
}

func (c *ListAliasesCommand) Synopsis() string {
	return "list account aliases for all accounts within an organization"
}
