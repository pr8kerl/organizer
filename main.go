package main

import (
	"fmt"
	"os"

	"github.com/mitchellh/cli"
)

func main() {

	c := cli.NewCLI("organizer", "0.0.1")
	c.Args = os.Args[1:]

	c.Commands = map[string]cli.CommandFactory{
		"list":            listCmdFactory,
		"list accounts":   listAccountsCmdFactory,
		"list buckets":    listBucketsCmdFactory,
		"list cloudfront": listCloudfrontsCmdFactory,
		"list users":      listUsersCmdFactory,
		"create":          createCmdFactory,
		"create account":  createAccountCmdFactory,
		"trails":          trailsCmdFactory,
	}

	exitStatus, err := c.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	os.Exit(exitStatus)
}
