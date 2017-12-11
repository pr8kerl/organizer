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
type ListCloudfrontsCommand struct {
	All bool
	Ui  cli.Ui
}

func listCloudfrontsCmdFactory() (cli.Command, error) {

	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}

	return &ListCloudfrontsCommand{
		Ui: ui,
	}, nil
}

func (c *ListCloudfrontsCommand) Run(args []string) int {

	cmdFlags := flag.NewFlagSet("list cloudfront", flag.ContinueOnError)
	cmdFlags.BoolVar(&c.All, "all", false, "show all cloudfront distributions")
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()); os.Exit(1) }
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	org, err := aws.NewOrganization()
	if err != nil {
		fmt.Printf("error: could not initialize organization: %s\n", err)
		return 1
	}

	err = org.PrintCloudfronts()
	if err != nil {
		fmt.Printf("error: could not list cloudfront distributions: %s\n", err)
		return 1
	}

	return 0
}

func (c *ListCloudfrontsCommand) Help() string {
	helpText := `usage: organizer list cloudfront [<args>]

List cloudfront distributions per account

	`
	return strings.TrimSpace(helpText)
}

func (c *ListCloudfrontsCommand) Synopsis() string {
	return "list all cloudfront distributions for a organizational accounts"
}
