package aws

import (
  "testing"
	"github.com/aws/aws-sdk-go/service/organizations/organizationsiface"
)

var (
  testRegions []string
)

type mockOrganizationsSvc struct {
    organizationsiface.OrganizationsAPI
}

func init() {

	testRegions = []string{
  	"ap-south-1",
  	"eu-west-2",
  	"eu-west-1",
  	"ap-northeast-2",
  	"ap-northeast-1",
  	"sa-east-1",
  	"ca-central-1",
  	"ap-southeast-1",
  	"ap-southeast-2",
  	"eu-central-1",
  	"us-east-1",
  	"us-east-2",
  	"us-west-1",
  	"us-west-2",
	}

}

func (o *Organization) SetSvc(svc *mockOrganizationsSvc)  {
  o.svc = svc
}
func (o *Organization) SetRegions(regions []string)  {
  o.regions = regions
}

func NewMockOrganization() (*Organization, error) {

  org, err := NewOrganization()
	if err != nil {
    return nil, err
	}

  mockSvc := &mockOrganizationsSvc{}
	org.SetSvc(mockSvc)
	org.SetRegions(testRegions)
	return org, nil

}

func TestGetRegions(t *testing.T) {

  org, err := NewMockOrganization()
	if err != nil {
    t.Errorf("could not create mock organization: %s", err)
	}
	regions := org.GetRegions()

  if len(testRegions) != len(regions) {
    t.Errorf("organization GetRegions error")
	}

 for i := range testRegions {
        if testRegions[i] != regions[i] {
            t.Errorf("organization GetRegions return value incorrect")
        }
 }


}
