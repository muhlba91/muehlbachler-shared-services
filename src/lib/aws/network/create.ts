import * as aws from '@pulumi/aws';

import { AWSVpcData } from '../../../model/data/aws/vpc';
import { commonLabels } from '../../configuration';

/**
 * Creates the VPC.
 *
 * @param {string} name the VPC name
 * @param {boolean} ipv6Only defines if this VPC is IPv6 only
 * @param {string} cidr the VPC CIDR block
 * @returns {Promise<AWSVpcData>} the security group
 */
export const createAWSVpc = async (
  name: string,
  ipv6Only: boolean,
  cidr: string,
): Promise<AWSVpcData> => {
  const vpc = new aws.ec2.Vpc(
    'vpc-' + name,
    {
      cidrBlock: cidr,
      assignGeneratedIpv6CidrBlock: true,
      enableDnsHostnames: true,
      enableDnsSupport: true,
      tags: {
        Name: name,
        ...commonLabels,
      },
    },
    {},
  );

  const igw = new aws.ec2.InternetGateway(
    'igw-' + name,
    {
      vpcId: vpc.id,
      tags: {
        Name: name,
        ...commonLabels,
      },
    },
    {},
  );

  const routeTable = new aws.ec2.RouteTable(
    'route-table-' + name,
    {
      vpcId: vpc.id,
      routes: [
        {
          cidrBlock: '0.0.0.0/0',
          gatewayId: igw.id,
        },
        {
          ipv6CidrBlock: '::/0',
          gatewayId: igw.id,
        },
      ],
      tags: {
        Name: name,
        ...commonLabels,
      },
    },
    {},
  );
  new aws.ec2.MainRouteTableAssociation(
    'main-route-table-association-' + name,
    {
      vpcId: vpc.id,
      routeTableId: routeTable.id,
    },
    {},
  );

  const subnets = (await aws.getAvailabilityZones({ state: 'available' })).names
    .sort()
    .map(
      (az, idx) =>
        new aws.ec2.Subnet(
          'subnet-' + name + '-' + az,
          {
            vpcId: vpc.id,
            availabilityZone: az,
            cidrBlock: ipv6Only
              ? undefined
              : vpc.cidrBlock.apply((cidr) =>
                  cidr.replace('.0.0/16', `.${idx}.0/24`),
                ),
            ipv6CidrBlock: vpc.ipv6CidrBlock.apply((cidr) =>
              cidr.replace('00::/56', `0${idx}::/64`),
            ),
            assignIpv6AddressOnCreation: true,
            mapPublicIpOnLaunch: true,
            ipv6Native: ipv6Only,
            enableDns64: false,
            enableResourceNameDnsAaaaRecordOnLaunch: true,
            enableResourceNameDnsARecordOnLaunch: !ipv6Only,
            tags: {
              Name: name + '-' + az,
              ...commonLabels,
            },
          },
          {},
        ),
    );

  return {
    vpc: vpc,
    subnets: subnets,
  };
};
