package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sts"
)

type BucketsPerAccount map[string][]string

func (s BucketsPerAccount) Add(key, value string) {
	_, ok := s[key]
	if !ok {
		s[key] = make([]string, 0, 20)
	}
	s[key] = append(s[key], value)
}

func (s BucketsPerAccount) Get(key string) (string, bool) {
	slice, ok := s[key]
	if !ok || len(slice) == 0 {
		return "", false
	}
	return s[key][0], true
}

func (s BucketsPerAccount) Set(key string, value []string) {
	s[key] = value
}

func (o *Organization) GetBucketsForAccount(accountid string) ([]string, error) {

	role := fmt.Sprintf("arn:aws:iam::%v:role/OrganizationAccountAccessRole", accountid)

	sess := session.Must(session.NewSession())
	stssvc := sts.New(sess)

	params := &sts.AssumeRoleInput{
		RoleArn:         aws.String(role),                // Required
		RoleSessionName: aws.String("organizer-buckets"), // Required
		DurationSeconds: aws.Int64(900),
	}
	resp, err := stssvc.AssumeRole(params)
	if err != nil {
		fmt.Printf("error: could not assume role %s\n", err.Error())
		return nil, err
	}

	config := aws.NewConfig().WithCredentials(
		credentials.NewStaticCredentials(
			*resp.Credentials.AccessKeyId,
			*resp.Credentials.SecretAccessKey,
			*resp.Credentials.SessionToken,
		),
	)

	sess = session.New(config)
	svc := s3.New(sess)

	input := &s3.ListBucketsInput{}
	bresp, err := svc.ListBuckets(input)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return nil, err
	}

	// Pretty-print the response data.
	// fmt.Println(resp)

	buckets := make([]string, 0, 200)

	for _, bucket := range bresp.Buckets {
		buckets = append(buckets, *bucket.Name)
	}

	return buckets, nil

}

func (o *Organization) GetBuckets() (map[string][]string, error) {

	accounts, err := o.GetActiveAccounts()
	if err != nil {
		return nil, err
	}

	buckets := make(BucketsPerAccount)

	for _, account := range accounts {
		bucks, err := o.GetBucketsForAccount(*account.Id)
		if err != nil {
			fmt.Printf("warning: could not list buckets for account %s\n\twarning: %s\n", *account.Name, err)
			continue
		}
		buckets.Set(*account.Id, bucks)
	}
	return buckets, nil

}

func (o *Organization) PrintBucketsForAccount(accountid string) error {
	buckets, err := o.GetBucketsForAccount(accountid)
	if err != nil {
		return err
	}

	if len(buckets) == 0 {
		fmt.Printf("warning: no buckets found in account %s\n", accountid)
		return nil
	}

	for _, bucket := range buckets {
		fmt.Printf("%s,%s\n", accountid, bucket)
	}

	return nil

}

func (o *Organization) PrintBuckets() error {
	buckets, err := o.GetBuckets()
	if err != nil {
		return err
	}

	for accountid, accountbuckets := range buckets {
		for _, bucket := range accountbuckets {
			fmt.Printf("%s,%s\n", accountid, bucket)
		}
	}
	return nil

}
