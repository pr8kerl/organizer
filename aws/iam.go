package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
)

type UsersPerAccount map[string][]*iam.User

func (s UsersPerAccount) Add(account string, value *iam.User) {
	_, ok := s[account]
	if !ok {
		s[account] = make([]*iam.User, 0, 100)
	}
	s[account] = append(s[account], value)
}

func (s UsersPerAccount) Get(key string) ([]*iam.User, bool) {
	slice, ok := s[key]
	if !ok || len(slice) == 0 {
		return nil, false
	}
	return s[key], true
}

func (s UsersPerAccount) Set(key string, value []*iam.User) {
	s[key] = value
}

func (o *Organization) GetUsersForAccount(accountid string) ([]*iam.User, error) {

	role := fmt.Sprintf("arn:aws:iam::%v:role/OrganizationAccountAccessRole", accountid)

	sess := session.Must(session.NewSession())
	stssvc := sts.New(sess)

	params := &sts.AssumeRoleInput{
		RoleArn:         aws.String(role),            // Required
		RoleSessionName: aws.String("organizer-iam"), // Required
		DurationSeconds: aws.Int64(900),
	}
	resp, err := stssvc.AssumeRole(params)
	if err != nil {
		return nil, err
	}

	config := aws.NewConfig().WithCredentials(
		credentials.NewStaticCredentials(
			*resp.Credentials.AccessKeyId,
			*resp.Credentials.SecretAccessKey,
			*resp.Credentials.SessionToken,
		),
	).WithRegion(o.region)

	sess = session.New(config)
	svc := iam.New(sess)

	input := &iam.ListUsersInput{}
	uresp, err := svc.ListUsers(input)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return nil, err
	}

	// Pretty-print the response data.
	// fmt.Println(resp)

	users := make([]*iam.User, 0, 500)

	for _, user := range uresp.Users {
		users = append(users, user)
	}

	return users, nil

}

func (o *Organization) GetUsers() (map[string][]*iam.User, error) {

	accounts, err := o.GetActiveAccounts()
	if err != nil {
		return nil, err
	}

	users := make(UsersPerAccount)

	for _, account := range accounts {
		lusers, err := o.GetUsersForAccount(*account.Id)
		if err != nil {
			fmt.Printf("warning: could not list users for account %s\n\twarning: %s\n", *account.Name, err)
			continue
		}
		users.Set(*account.Id, lusers)
	}
	return users, nil

}

func (o *Organization) PrintUsersForAccount(accountid string) error {
	users, err := o.GetUsersForAccount(accountid)
	if err != nil {
		return err
	}

	if len(users) == 0 {
		fmt.Printf("warning: no users found in account %s\n", accountid)
		return nil
	}

	for _, user := range users {
		fmt.Printf("%s,%s,%s,%s\n", accountid, *user.UserName, user.CreateDate, user.PasswordLastUsed)
	}

	return nil

}

func (o *Organization) PrintUsers() error {
	users, err := o.GetUsers()
	if err != nil {
		return err
	}

	for accountid, accountusers := range users {
		for _, user := range accountusers {
			fmt.Printf("%s,%s,%s,%s\n", accountid, *user.UserName, user.CreateDate, user.PasswordLastUsed)
		}
	}
	return nil

}
