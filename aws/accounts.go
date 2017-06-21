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

	accounts, err := o.GetAccounts()
	if err != nil {
		return err
	}

	for _, account := range accounts {
		if *account.Status == "ACTIVE" {
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

func (o *Organization) PollForAccountStatus(accountStatus *organizations.CreateAccountStatus, progress bool) (*organizations.CreateAccountStatus, error) {

  if accountStatus == nil {
		err := fmt.Errorf("error: account status object is nil\n")
    return nil, err
	}

	statusInput := &organizations.DescribeCreateAccountStatusInput{
		CreateAccountRequestId: accountStatus.Id,
	}
  timeout := make(chan bool, 1)
  result := make(chan *organizations.CreateAccountStatus, 1)
  pollerr := make(chan error, 1)

  // set timeout for 5 mins max
  go func() {
        time.Sleep(time.Minute * 5)
        timeout <- true
  }()

  go func() {
	// wait until the request has completed
	  for *accountStatus.State == "IN_PROGRESS" {

  		time.Sleep(time.Second * 10)
  		statusOutput, err := o.svc.DescribeCreateAccountStatus(statusInput)
  		if err != nil {
  			err := fmt.Errorf("error: could not check account status: %s", err.Error())
  			pollerr <- err
  		}
  		accountStatus = statusOutput.CreateAccountStatus
      if progress {
        fmt.Printf("account status: %s\n", *accountStatus.State)
      }

  	}
    if progress {
      fmt.Printf("account status: %s\n", *accountStatus.State)
    }

  	if *accountStatus.State == "FAILED" {
  		err := fmt.Errorf("error: failed to create account: %s", *accountStatus.FailureReason)
  		pollerr <- err
  	}

    result <- accountStatus

	}()

  // wait for a result
  select {
    case <-result:
  	  return <-result, nil
    case <-timeout:
    	err := fmt.Errorf("error: timeout after 5 mins")
  	  return nil, err
    case <-pollerr:
  	  return nil, <-pollerr
  }

}
