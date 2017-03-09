package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/organizations"
	"github.com/mitchellh/cli"
	"os"
)

/*
 {
  Arn: "arn:aws:organizations::376681487066:account/o-j6lf48521y/833136982555",
  Id: "833136982555",
  JoinedMethod: "INVITED",
  JoinedTimestamp: 2016-12-04 21:49:05 +0000 UTC,
  Name: "EX Task Service Production",
  Status: "ACTIVE"
}
*/

type AccountsCommand struct {
	Region string
	All    bool
	Ui     cli.Ui
}

func accountsCmdFactory() (cli.Command, error) {

	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}

	return &AccountsCommand{
		Region: "us-east-1",
		All:    false,
		Ui:     ui,
	}, nil
}

func (c *AccountsCommand) Run(args []string) int {

	cmdFlags := flag.NewFlagSet("accounts", flag.ContinueOnError)
	cmdFlags.BoolVar(&c.All, "all", false, "show all accounts including inactive")
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()); os.Exit(1) }
	cmdFlags.Parse(args)

	config := aws.NewConfig().WithRegion(c.Region)
	sess := session.New(config)
	svc := organizations.New(sess)

	params := &organizations.ListAccountsInput{}
	accounts := []*organizations.Account{}
	var nextToken *string

	for {
		resp, err := svc.ListAccounts(params)
		if err != nil {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
			return 1
		}
		accounts = append(accounts, resp.Accounts...)
		nextToken = resp.NextToken
		if isNilOrEmpty(nextToken) {
			break
		}
		params.NextToken = nextToken
	}

	c.printAccountInfo(accounts)
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

func (c *AccountsCommand) printAccountInfo(accounts []*organizations.Account) {
	for _, account := range accounts {
		fmt.Printf("%s,%s,%s\n", *account.Id, *account.Name, *account.Status)
	}
}
