package vpc

import (
	"fmt"
	"sort"
	"strings"

	"github.com/muhlba91/pulumi-shared-library/pkg/util/metadata"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"github.com/muhlba91/muehlbachler-shared-services/pkg/lib/config"
	vpcModel "github.com/muhlba91/muehlbachler-shared-services/pkg/model/aws/vpc"
	awsConfig "github.com/muhlba91/muehlbachler-shared-services/pkg/model/config/aws"
)

// Create creates a VPC with an IGW, main route table, and one public subnet per AZ.
// ctx: Pulumi context.
// name: Name of the VPC.
// vpcConfig: VPC configuration.
//
//nolint:funlen
func Create(ctx *pulumi.Context, name string, vpcConfig *awsConfig.VPCConfig) (*vpcModel.Data, error) {
	tags := config.CommonLabels()
	tags["Name"] = name
	pTags := metadata.LabelsToStringMap(tags)

	vpc, vErr := ec2.NewVpc(ctx, fmt.Sprintf("vpc-%s", name), &ec2.VpcArgs{
		CidrBlock:                    pulumi.StringPtr(vpcConfig.CIDR),
		AssignGeneratedIpv6CidrBlock: pulumi.Bool(true),
		EnableDnsHostnames:           pulumi.Bool(true),
		EnableDnsSupport:             pulumi.Bool(true),
		Tags:                         pTags,
	})
	if vErr != nil {
		return nil, vErr
	}

	igw, igwErr := ec2.NewInternetGateway(ctx, fmt.Sprintf("igw-%s", name), &ec2.InternetGatewayArgs{
		VpcId: vpc.ID(),
		Tags:  pTags,
	})
	if igwErr != nil {
		return nil, igwErr
	}

	routeTable, rtErr := ec2.NewRouteTable(ctx, fmt.Sprintf("route-table-%s", name), &ec2.RouteTableArgs{
		VpcId: vpc.ID(),
		Routes: ec2.RouteTableRouteArray{
			&ec2.RouteTableRouteArgs{
				CidrBlock: pulumi.String("0.0.0.0/0"),
				GatewayId: igw.ID(),
			},
			&ec2.RouteTableRouteArgs{
				Ipv6CidrBlock: pulumi.String("::/0"),
				GatewayId:     igw.ID(),
			},
		},
		Tags: pTags,
	})
	if rtErr != nil {
		return nil, rtErr
	}

	_, rtaErr := ec2.NewMainRouteTableAssociation(
		ctx,
		fmt.Sprintf("main-route-table-association-%s", name),
		&ec2.MainRouteTableAssociationArgs{
			VpcId:        vpc.ID(),
			RouteTableId: routeTable.ID(),
		},
	)
	if rtaErr != nil {
		return nil, rtaErr
	}

	azs, azErr := aws.GetAvailabilityZones(ctx, &aws.GetAvailabilityZonesArgs{
		State: pulumi.StringRef("available"),
	}, nil)
	if azErr != nil {
		return nil, azErr
	}

	names := append([]string{}, azs.Names...)
	sort.Strings(names)

	subnets := make([]*ec2.Subnet, 0, len(names))
	for idx, az := range names {
		var ipv4Cidr pulumi.StringPtrInput
		if !vpcConfig.IPv6Only {
			ipv4Cidr, _ = vpc.CidrBlock.ApplyT(func(cidr string) *string {
				s := strings.Replace(cidr, ".0.0/16", fmt.Sprintf(".%d.0/24", idx), 1)
				return &s
			}).(pulumi.StringPtrOutput)
		}

		ipv6Cidr, _ := vpc.Ipv6CidrBlock.ApplyT(func(block string) *string {
			s := strings.Replace(block, "00::/56", fmt.Sprintf("0%d::/64", idx), 1)
			return &s
		}).(pulumi.StringPtrOutput)

		tags["Name"] = fmt.Sprintf("%s-%s", name, az)

		sn, snErr := ec2.NewSubnet(ctx, fmt.Sprintf("subnet-%s-%s", name, az), &ec2.SubnetArgs{
			VpcId:                                   vpc.ID(),
			AvailabilityZone:                        pulumi.String(az),
			CidrBlock:                               ipv4Cidr,
			Ipv6CidrBlock:                           ipv6Cidr,
			AssignIpv6AddressOnCreation:             pulumi.Bool(true),
			MapPublicIpOnLaunch:                     pulumi.Bool(true),
			Ipv6Native:                              pulumi.Bool(vpcConfig.IPv6Only),
			EnableDns64:                             pulumi.Bool(false),
			EnableResourceNameDnsAaaaRecordOnLaunch: pulumi.Bool(true),
			EnableResourceNameDnsARecordOnLaunch:    pulumi.Bool(!vpcConfig.IPv6Only),
			Tags:                                    metadata.LabelsToStringMap(tags),
		})
		if snErr != nil {
			return nil, snErr
		}
		subnets = append(subnets, sn)
	}

	return &vpcModel.Data{
		Vpc:     vpc,
		Subnets: subnets,
	}, nil
}
