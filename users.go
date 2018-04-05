package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/pr8kerl/organizer/aws"
)

// List Users
type ListUsersCommand struct {
	AccountId string
	Report    bool
	Region    string
	Ui        cli.Ui
}

func listUsersCmdFactory() (cli.Command, error) {

	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}

	return &ListUsersCommand{
		AccountId: "",
		Report:    false,
		Region:    "",
		Ui:        ui,
	}, nil
}

func (c *ListUsersCommand) Run(args []string) int {

	cmdFlags := flag.NewFlagSet("list users", flag.ContinueOnError)
	cmdFlags.StringVar(&c.AccountId, "accountid", "", "list iam users for a specific account")
	cmdFlags.BoolVar(&c.Report, "report", false, "generate and show iam credential reports distributions")
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

		if c.Report {
			err = org.GenerateCredentialReportForAccount(c.AccountId)
			if err != nil {
				fmt.Printf("error: could generate iam credential report for account: %s\n", err)
				return 1
			}
			err = org.PrintCredentialReportForAccount(c.AccountId)
			if err != nil {
				fmt.Printf("error: could get iam credential report for account: %s\n", err)
				return 1
			}
		} else {
			err = org.PrintUsersForAccount(c.AccountId)
			if err != nil {
				fmt.Printf("error: could not list iam users for account: %s\n", err)
				return 1
			}
		}
	} else {
		if c.Report {

			err = org.GenerateCredentialReports()
			if err != nil {
				fmt.Printf("error: could not generate all iam credential reports: %s\n", err)
				return 1
			}
			err = org.PrintCredentialReports()
			if err != nil {
				fmt.Printf("error: could not print all iam credential reports: %s\n", err)
				return 1
			}

		} else {
			err = org.PrintUsers()
			if err != nil {
				fmt.Printf("error: could not list all iam users: %s\n", err)
				return 1
			}
		}
	}

	return 0
}

func (c *ListUsersCommand) Help() string {
	helpText := `usage: organizer list users [<args>]

List all iam users within accounts for an organization

Options:
	    -accountid		list all iam users for a specific account only
	    -report				run a credential report
	`
	return strings.TrimSpace(helpText)
}

func (c *ListUsersCommand) Synopsis() string {
	return "list all iam users for all accounts in an organization"
}
