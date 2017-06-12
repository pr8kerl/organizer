package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/organizations"
	"time"
	//	"github.com/kr/pretty"
)

func (o *Organization) GetAccounts() ([]*organizations.Account, error) {

	params := &organizations.ListAccountsInput{}
	accounts := make([]*organizations.Account, 0, 200)
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

	activeAccounts := make([]*organizations.Account, 0, 200)
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
			fmt.Printf("%s,%s\n", *account.Id, *account.Name)
		}
	}
	return nil
}

func (o *Organization) CreateAccount(name string, email string) (*organizations.Account, error) {
	input := &organizations.CreateAccountInput{
		AccountName: aws.String(name),
		Email:       aws.String(email),
		IamUserAccessToBilling: aws.String("ALLOW"),
	}

	result, err := o.svc.CreateAccount(input)
	if err != nil {
		err := fmt.Errorf("error: could not create account: %s", err.Error())
		return nil, err
	}

	// wait until the request has completed
	for *result.CreateAccountStatus.State == "IN_PROGRESS" {
		fmt.Println("account creation in progress")
		time.Sleep(time.Second * 10)
	}

	if *result.CreateAccountStatus.State == "FAILED" {
		err := fmt.Errorf("error: failed to create account: %s", *result.CreateAccountStatus.FailureReason)
		return nil, err
	}
	account_input := &organizations.DescribeAccountInput{
		AccountId: result.CreateAccountStatus.AccountId,
	}

	account_output, err := o.svc.DescribeAccount(account_input)
	if err != nil {
		err := fmt.Errorf("error: failed to describe new account %s: %s", *result.CreateAccountStatus.AccountId, *result.CreateAccountStatus.FailureReason)
		return nil, err
	}
	return account_output.Account, nil

}
