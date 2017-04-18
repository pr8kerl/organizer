package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudtrail"
	"github.com/aws/aws-sdk-go/service/organizations"
	"github.com/aws/aws-sdk-go/service/sts"
)

type Organization struct {
	svc      *organizations.Organizations
	accounts []*organizations.Account
	regions  []string
}

func newOrganization() (*Organization, error) {

	config := aws.NewConfig().WithRegion("us-east-1")
	sess := session.New(config)
	svc := organizations.New(sess)
	accounts := make([]*organizations.Account, 0, 100)
	regions := make([]string, 0, 20)

	return &Organization{
		svc:      svc,
		accounts: accounts,
		regions:  regions,
	}, nil

}

func (o *Organization) GetRegions() []string {

	if len(o.regions) == 0 {
		p := endpoints.AwsPartition()
		for id := range p.Regions() {
			o.regions = append(o.regions, id)
		}
	}

	return o.regions

}

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

func (o *Organization) GetTrailArnsForAccount(accountid string) ([]string, error) {

	regions := o.GetRegions()
	role := fmt.Sprintf("arn:aws:iam::%v:role/OrganizationAccountAccessRole", accountid)

	sess := session.Must(session.NewSession())
	stssvc := sts.New(sess)

	params := &sts.AssumeRoleInput{
		RoleArn:         aws.String(role),                     // Required
		RoleSessionName: aws.String("organizer-cloudtrailer"), // Required
		DurationSeconds: aws.Int64(900),
	}
	resp, err := stssvc.AssumeRole(params)
	if err != nil {
		fmt.Printf("error: could not assume role %s\n", err.Error())
		return nil, err
	}

	trailmap := make(map[string]bool, 100)

	for _, region := range regions {

		config := aws.NewConfig().WithCredentials(
			credentials.NewStaticCredentials(
				*resp.Credentials.AccessKeyId,
				*resp.Credentials.SecretAccessKey,
				*resp.Credentials.SessionToken,
			),
		).WithRegion(region)

		sess = session.New(config)
		svc := cloudtrail.New(sess)

		params := &cloudtrail.DescribeTrailsInput{
			IncludeShadowTrails: aws.Bool(true),
			TrailNameList:       []*string{},
		}
		resp, err := svc.DescribeTrails(params)

		if err != nil {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
			return nil, err
		}

		// Pretty-print the response data.
		// fmt.Println(resp)
		if len(resp.TrailList) == 0 {
			fmt.Printf("warning: no trails defined in account %s in region %s\n", accountid, region)
			continue
		} else {
			for _, trail := range resp.TrailList {
				trailmap[*trail.TrailARN] = true
			}
		}

	}
	trails := make([]string, len(trailmap))
	i := 0
	for trail, _ := range trailmap {
		trails[i] = trail
		i++
	}

	return trails, nil

}

func (o *Organization) PurgeTrailsForAccount(accountid string) ([]string, error) {

	regions := o.GetRegions()
	role := fmt.Sprintf("arn:aws:iam::%v:role/OrganizationAccountAccessRole", accountid)

	sess := session.Must(session.NewSession())
	stssvc := sts.New(sess)

	params := &sts.AssumeRoleInput{
		RoleArn:         aws.String(role),                     // Required
		RoleSessionName: aws.String("organizer-cloudtrailer"), // Required
		DurationSeconds: aws.Int64(900),
	}
	resp, err := stssvc.AssumeRole(params)
	if err != nil {
		fmt.Printf("error: could not assume role %s\n", err.Error())
		return nil, err
	}

	trailmap := make(map[string]bool, 100)

	for _, region := range regions {

		config := aws.NewConfig().WithCredentials(
			credentials.NewStaticCredentials(
				*resp.Credentials.AccessKeyId,
				*resp.Credentials.SecretAccessKey,
				*resp.Credentials.SessionToken,
			),
		).WithRegion(region)

		sess = session.New(config)
		svc := cloudtrail.New(sess)

		params := &cloudtrail.DescribeTrailsInput{
			IncludeShadowTrails: aws.Bool(true),
			TrailNameList:       []*string{},
		}
		resp, err := svc.DescribeTrails(params)

		if err != nil {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
			return nil, err
		}

		// Pretty-print the response data.
		// fmt.Println(resp)
		if len(resp.TrailList) == 0 {
			fmt.Printf("warning: no trails defined in account %s in region %s\n", accountid, region)
			continue
		} else {
			for _, trail := range resp.TrailList {
				trailmap[*trail.TrailARN] = true
				params := &cloudtrail.DeleteTrailInput{
					Name: trail.TrailARN,
				}
				_, err := svc.DeleteTrail(params)
				if err != nil {
					// Print the error, cast err to awserr.Error to get the Code and
					// Message from an error.
					fmt.Println(err.Error())
					return nil, err
				}
			}
		}

	}
	trails := make([]string, len(trailmap))
	i := 0
	for trail, _ := range trailmap {
		trails[i] = trail
		i++
	}

	return trails, nil

}
