package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/organizations"
)

type Organization struct {
	Svc      *organizations.Organizations
	Accounts []*organizations.Account
}

func newOrganization() (*Organization, error) {

	config := aws.NewConfig().WithRegion("us-east-1")
	sess := session.New(config)
	svc := organizations.New(sess)
	accounts := make([]*organizations.Account, 0, 100)

	return &Organization{
		Svc:      svc,
		Accounts: accounts,
	}, nil

}

func (o *Organization) GetAccounts() ([]*organizations.Account, error) {

	params := &organizations.ListAccountsInput{}
	accounts := []*organizations.Account{}
	var nextToken *string

	for {
		resp, err := o.Svc.ListAccounts(params)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, resp.Accounts...)
		nextToken = resp.NextToken
		if isNilOrEmpty(nextToken) {
			break
		}
		params.NextToken = nextToken
	}
	o.Accounts = accounts

	return accounts, nil

}

func (o *Organization) GetActiveAccounts() ([]*organizations.Account, error) {

	activeAccounts := []*organizations.Account{}
	accounts, err := o.GetAccounts()
	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
	}
	for _, account := range accounts {
		if *account.Status == "ACTIVE" {
			activeAccounts = append(accounts, account)
		}
	}
	return activeAccounts, nil

}

func (o *Organization) PrintAccounts() error {

	if len(o.Accounts) == 0 {
		_, err := o.GetAccounts()
		if err != nil {
			return err
		}
	}

	for _, account := range o.Accounts {
		fmt.Printf("%s,%s,%s\n", *account.Id, *account.Name, *account.Status)
	}

	return nil
}

func (o *Organization) PrintActiveAccounts() error {

	if len(o.Accounts) == 0 {
		_, err := o.GetAccounts()
		if err != nil {
			return err
		}
	}
	for _, account := range o.Accounts {
		if *account.Status == "ACTIVE" {
			fmt.Printf("%s,%s,%s\n", *account.Id, *account.Name, *account.Status)
		}
	}
	return nil
}
