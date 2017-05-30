package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/organizations"
)

func (o *Organization) GetAccounts() ([]*organizations.Account, error) {

	params := &organizations.ListAccountsInput{}
	accounts := make([]*organizations.Account, 0, 100)
	var nextToken *string

	for {
		resp, err := o.svc.ListAccounts(params)
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
	o.accounts = accounts

	return accounts, nil

}

func (o *Organization) GetActiveAccounts() ([]*organizations.Account, error) {

	activeAccounts := make([]*organizations.Account, 0, 100)
	accounts, err := o.GetAccounts()
	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		return nil, err
	}
	for _, account := range accounts {
		if *account.Status == "ACTIVE" {
			activeAccounts = append(accounts, account)
		}
	}
	return activeAccounts, nil

}

func (o *Organization) PrintAccounts() error {

	accounts, err := o.GetAccounts()
	if err != nil {
		return err
	}

	for _, account := range accounts {
		fmt.Printf("%s,%s,%s\n", *account.Id, *account.Name, *account.Status)
	}

	return nil
}

func (o *Organization) PrintActiveAccounts() error {

	accounts, err := o.GetAccounts()
	if err != nil {
		return err
	}

	for _, account := range accounts {
		if *account.Status == "ACTIVE" {
			fmt.Printf("%s,%s,%s\n", *account.Id, *account.Name, *account.Status)
		}
	}
	return nil
}
