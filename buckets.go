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
type ListBucketsCommand struct {
	AccountId string
	Ui        cli.Ui
}

func listBucketsCmdFactory() (cli.Command, error) {

	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}

	return &ListBucketsCommand{
		AccountId: "",
		Ui:        ui,
	}, nil
}

func (c *ListBucketsCommand) Run(args []string) int {

	cmdFlags := flag.NewFlagSet("list buckets", flag.ContinueOnError)
	cmdFlags.StringVar(&c.AccountId, "accountid", "", "list s3 buckets for a specific account")
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
		err = org.PrintBucketsForAccount(c.AccountId)
		if err != nil {
			fmt.Printf("error: could not list s3 buckets for account: %s\n", err)
			return 1
		}
	} else {
		err = org.PrintBuckets()
		if err != nil {
			fmt.Printf("error: could not list all s3 buckets: %s\n", err)
			return 1
		}
	}

	return 0
}

func (c *ListBucketsCommand) Help() string {
	helpText := `usage: organizer list buckets [<args>]

List all s3 buckets within accounts for an organization

Options:
	    -accountid		list all s3 buckets for a specific account only
	`
	return strings.TrimSpace(helpText)
}

func (c *ListBucketsCommand) Synopsis() string {
	return "list all s3 buckets for all accounts in an organization"
}
