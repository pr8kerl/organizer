package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/sts"
)

type CloudfrontsPerAccount map[string][]*cloudfront.DistributionSummary

func (s CloudfrontsPerAccount) Add(account string, value *cloudfront.DistributionSummary) {
	_, ok := s[account]
	if !ok {
		s[account] = make([]*cloudfront.DistributionSummary, 0, 100)
	}
	s[account] = append(s[account], value)
}

func (s CloudfrontsPerAccount) Get(key string) ([]*cloudfront.DistributionSummary, bool) {
	slice, ok := s[key]
	if !ok || len(slice) == 0 {
		return nil, false
	}
	return s[key], true
}

func (s CloudfrontsPerAccount) Set(key string, value []*cloudfront.DistributionSummary) {
	s[key] = value
}

func (o *Organization) GetCloudfrontsForAccount(accountid string) ([]*cloudfront.DistributionSummary, error) {

	role := fmt.Sprintf("arn:aws:iam::%v:role/OrganizationAccountAccessRole", accountid)

	sess := session.Must(session.NewSession())
	stssvc := sts.New(sess)

	params := &sts.AssumeRoleInput{
		RoleArn:         aws.String(role),                   // Required
		RoleSessionName: aws.String("organizer-cloudfront"), // Required
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
	svc := cloudfront.New(sess)

	cparams := &cloudfront.ListDistributionsInput{
		MaxItems: aws.Int64(100),
	}
	cresp, err := svc.ListDistributions(cparams)
	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return nil, err
	}

	// Pretty-print the response data.
	// fmt.Println(resp)

	distros := make([]*cloudfront.DistributionSummary, 0, 500)

	for _, distro := range cresp.DistributionList.Items {
		distros = append(distros, distro)
	}

	return distros, nil

}

func (o *Organization) GetCloudfronts() (map[string][]*cloudfront.DistributionSummary, error) {

	accounts, err := o.GetActiveAccounts()
	if err != nil {
		return nil, err
	}

	distros := make(CloudfrontsPerAccount)

	for _, account := range accounts {
		distributions, err := o.GetCloudfrontsForAccount(*account.Id)
		if err != nil {
			fmt.Printf("warning: could not list users for account %s\n\twarning: %s\n", *account.Name, err)
			continue
		}
		distros.Set(*account.Id, distributions)
	}
	return distros, nil

}

func (o *Organization) PrintCloudfrontsForAccount(accountid string) error {
	distros, err := o.GetCloudfrontsForAccount(accountid)
	if err != nil {
		return err
	}

	if len(distros) == 0 {
		fmt.Printf("warning: no cloudfront distributions found in account %s\n", accountid)
		return nil
	}

	for _, distro := range distros {
		fmt.Printf("%s,%s\n", accountid, *distro.DomainName)
	}

	return nil

}

func (o *Organization) PrintCloudfronts() error {
	distros, err := o.GetCloudfronts()
	if err != nil {
		return err
	}

	for accountid, accountdistros := range distros {
		for _, distro := range accountdistros {
			fmt.Printf("%s,%s\n", accountid, *distro.DomainName)
		}
	}
	return nil

}
