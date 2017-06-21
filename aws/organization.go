package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/organizations"
	"github.com/aws/aws-sdk-go/service/organizations/organizationsiface"
)

type Organization struct {
	svc      organizationsiface.OrganizationsAPI
	accounts []*organizations.Account
	regions  []string
	region   string
}

func NewOrganization() (*Organization, error) {

	region := "us-east-1"
	config := aws.NewConfig().WithRegion(region)
	sess := session.New(config)
	svc := organizations.New(sess)
	accounts := make([]*organizations.Account, 0, 100)
	regions := make([]string, 0, 20)

	return &Organization{
		svc:      svc,
		accounts: accounts,
		regions:  regions,
		region:   region,
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
