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
type AliasesPerAccount map[string][]*string

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

func (s AliasesPerAccount) Set(key string, value []*string) {
	s[key] = value
}

func (o *Organization) GetIamSvcForAccount(accountid string) (*iam.IAM, error) {

	role := fmt.Sprintf("arn:aws:iam::%v:role/OrganizationAccountAccessRole", accountid)

	sess := session.Must(session.NewSession())
	stssvc := sts.New(sess)
	sess_name := "organizer-iam-" + accountid

	params := &sts.AssumeRoleInput{
		RoleArn:         aws.String(role),      // Required
		RoleSessionName: aws.String(sess_name), // Required
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
	return iam.New(sess), nil

}

func (o *Organization) GenerateCredentialReportForAccount(accountid string) error {

	svc, err := o.GetIamSvcForAccount(accountid)
	if err != nil {
		return err
	}

	input := &iam.GenerateCredentialReportInput{}
	resp, err := svc.GenerateCredentialReport(input)
	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return err
	}

	// Pretty-print the response data.
	fmt.Println(resp)

	return nil

}

func (o *Organization) GenerateCredentialReports() error {

	accounts, err := o.GetActiveAccounts()
	if err != nil {
		return err
	}

	for _, account := range accounts {
		err := o.GenerateCredentialReportForAccount(*account.Id)
		if err != nil {
			fmt.Printf("warning: could not generate credential report for account %s\n\twarning: %s\n", *account.Name, err)
			continue
		}
		fmt.Printf("info: generated credential report for account %s\n", *account.Id)
	}
	return nil

}

func (o *Organization) GetCredentialReportForAccount(accountid string) (string, error) {

	svc, err := o.GetIamSvcForAccount(accountid)
	if err != nil {
		return "", err
	}

	input := &iam.GetCredentialReportInput{}
	resp, err := svc.GetCredentialReport(input)
	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return "", err
	}

	// Pretty-print the response data.
	// fmt.Println(resp)
	n := len(resp.Content)

	return string(resp.Content[:n]), nil

}

func (o *Organization) PrintCredentialReportForAccount(accountid string) error {

	report, err := o.GetCredentialReportForAccount(accountid)
	if err != nil {
		return err
	}

	if len(report) == 0 {
		fmt.Printf("warning: no credential report found for account %s\n", accountid)
		return nil
	}

	fmt.Printf("\naccount id: %s\n", accountid)
	fmt.Printf("%s\n", report)

	return nil

}

func (o *Organization) PrintCredentialReports() error {

	accounts, err := o.GetActiveAccounts()
	if err != nil {
		return err
	}

	for _, account := range accounts {
		err := o.GenerateCredentialReportForAccount(*account.Id)
		if err != nil {
			fmt.Printf("warning: could not generate credentials report for account %s\n\twarning: %s\n", *account.Id, err)
			continue
		}
	}

	for _, account := range accounts {
		err := o.PrintCredentialReportForAccount(*account.Id)
		if err != nil {
			fmt.Printf("warning: could not print credentials report for account %s\n\twarning: %s\n", *account.Id, err)
			continue
		}
	}
	return nil

}

func (o *Organization) GetUsersForAccount(accountid string) ([]*iam.User, error) {

	svc, err := o.GetIamSvcForAccount(accountid)
	if err != nil {
		return nil, err
	}

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

func (o *Organization) GetAliasesForAccount(accountid string) ([]*string, error) {

	svc, err := o.GetIamSvcForAccount(accountid)
	if err != nil {
		return nil, err
	}

	input := &iam.ListAccountAliasesInput{}
	aresp, err := svc.ListAccountAliases(input)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return nil, err
	}

	// Pretty-print the response data.
	// fmt.Println(resp)

	aliases := make([]*string, 0, 500)

	for _, alias := range aresp.AccountAliases {
		aliases = append(aliases, alias)
	}

	return aliases, nil

}

func (o *Organization) GetAliases() (map[string][]*string, error) {

	accounts, err := o.GetActiveAccounts()
	if err != nil {
		return nil, err
	}

	aliases := make(AliasesPerAccount)

	for _, account := range accounts {
		alii, err := o.GetAliasesForAccount(*account.Id)
		if err != nil {
			fmt.Printf("warning: could not list aliases for account %s\n\twarning: %s\n", *account.Name, err)
			continue
		}
		aliases.Set(*account.Id, alii)
	}
	return aliases, nil

}

func (o *Organization) PrintAliases() error {
	aliases, err := o.GetAliases()
	if err != nil {
		return err
	}

	for accountid, accountaliases := range aliases {
		for _, alias := range accountaliases {
			fmt.Printf("%s,%s\n", accountid, *alias)
		}
	}

	return nil

}

func (o *Organization) PrintAliasesForAccount(accountid string) error {
	aliases, err := o.GetAliasesForAccount(accountid)
	if err != nil {
		return err
	}

	if len(aliases) == 0 {
		fmt.Printf("warning: no aliases found for account %s\n", accountid)
		return nil
	}

	for _, alias := range aliases {
		fmt.Printf("%s,%s\n", accountid, *alias)
	}

	return nil

}
