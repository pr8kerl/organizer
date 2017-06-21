package aws

import (
	"os"
)

func getRegion() string {

	// set the default to ap-southeast-2 if no other is set
	regenv := os.Getenv("AWS_REGION")
	if regenv != "" {
		return regenv
	}
	regenv = os.Getenv("AWS_DEFAULT_REGION")
	if regenv != "" {
		return regenv
	}
	return "us-east-1"

}

func isNilOrEmpty(s *string) bool {
	return s == nil || *s == ""
}
