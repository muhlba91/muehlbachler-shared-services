package vpc

import "github.com/pulumi/pulumi-aws/sdk/v7/go/aws/ec2"

// Data defines AWS VPC data.
type Data struct {
	// VPC is the AWS VPC resource.
	Vpc *ec2.Vpc
	// Subnets are the AWS Subnet resources associated with the VPC.
	Subnets []*ec2.Subnet
}
