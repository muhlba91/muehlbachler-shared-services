import * as aws from '@pulumi/aws';

/**
 * Defines AWS VPC data.
 */
export type AWSVpcData = {
  readonly vpc: aws.ec2.Vpc;
  readonly subnets: readonly aws.ec2.Subnet[];
};
