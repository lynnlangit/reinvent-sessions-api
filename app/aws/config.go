// Package aws helps developers to use AWS cloud API more friendly
package aws

import "github.com/aws/aws-sdk-go/aws"

func config() *aws.Config {
	awscfg := &aws.Config{}
	return awscfg
}
