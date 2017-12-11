package aws

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/organizations"
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

	accounts, err := o.GetActiveAccounts()
	if err != nil {
		return err
	}

	for _, account := range accounts {
		// for a newly created account, it can take a while for all the account fields
		// to be populated. So to avoid a panic...
		var name string = "unknown"
		var id string = "unknown"
		if account.Id != nil {
			id = *account.Id
		}
		if account.Name != nil {
			name = *account.Name
		}
		fmt.Printf("%s,%s\n", id, name)
	}
	return nil
}

func (o *Organization) CreateAccount(name string, email string) (*organizations.CreateAccountStatus, error) {
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

	return result.CreateAccountStatus, nil

}

func (o *Organization) WaitForAccountStatus(accountStatus *organizations.CreateAccountStatus) (*organizations.CreateAccountStatus, error) {

	if accountStatus == nil {
		err := fmt.Errorf("error: account status object is nil\n")
		return nil, err
	}

	statusInput := &organizations.DescribeCreateAccountStatusInput{
		CreateAccountRequestId: accountStatus.Id,
	}

	statusOutput, err := o.svc.DescribeCreateAccountStatus(statusInput)
	if err != nil {
		err := fmt.Errorf("error: could not check account status: %s", err.Error())
		return nil, err
	}

	// wait until the request has completed
	for *statusOutput.CreateAccountStatus.State == "IN_PROGRESS" {

		time.Sleep(time.Second * 10)
		statusOutput, err = o.svc.DescribeCreateAccountStatus(statusInput)
		if err != nil {
			err := fmt.Errorf("error: could not check account status: %s", err.Error())
			return nil, err
		}
		fmt.Printf("account status: %s\n", *statusOutput.CreateAccountStatus.State)
	}

	fmt.Printf("account status: %s\n", *statusOutput.CreateAccountStatus.State)

	if *statusOutput.CreateAccountStatus.State == "FAILED" {
		err := fmt.Errorf("error: failed to create account: %s", *statusOutput.CreateAccountStatus.FailureReason)
		return nil, err
	}

	return statusOutput.CreateAccountStatus, err
}
