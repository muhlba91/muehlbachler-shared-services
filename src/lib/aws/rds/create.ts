import * as aws from '@pulumi/aws';

import { AWSRdsConfig } from '../../../model/config/aws/rds';
import { AWSRdsData } from '../../../model/data/aws/rds';
import { commonLabels, environment, globalName } from '../../configuration';
import { writeToDoppler } from '../../util/doppler/secret';
import { createRandomPassword } from '../../util/random';
import { createAWSVpc } from '../network/create';

/**
 * Creates the RDS instance.
 *
 * @param {AWSRdsConfig} data the database data
 * @returns {Promise<AWSRdsData>} the database instance
 */
export const createRDSInstance = async (
  data: AWSRdsConfig,
): Promise<AWSRdsData> => {
  const instanceName = 'rds-' + globalName + '-' + data.name;
  const instanceIdentifier = data.name + '-' + environment;
  const instanceLabels = {
    Name: data.name,
    ...commonLabels,
  };

  const kmsKeyId = (
    await aws.kms.getAlias({
      name: 'alias/aws/rds',
    })
  ).arn;

  const vpc = await createAWSVpc(
    instanceName,
    data.vpc.ipv6Only,
    data.vpc.cidr,
  );

  const securityGroup = new aws.ec2.SecurityGroup(
    instanceName + '-security-group',
    {
      description: instanceName + ': PostgreSQL',
      vpcId: vpc.vpc.id,
      ingress: [
        {
          fromPort: 5432,
          toPort: 5432,
          protocol: 'tcp',
          cidrBlocks: ['0.0.0.0/0'],
          ipv6CidrBlocks: ['::/0'],
        },
      ],
      egress: [
        {
          fromPort: 0,
          toPort: 0,
          protocol: '-1',
          cidrBlocks: ['0.0.0.0/0'],
          ipv6CidrBlocks: ['::/0'],
        },
      ],
      tags: instanceLabels,
    },
    {},
  );

  const dbSubnet = new aws.rds.SubnetGroup(
    instanceName + '-subnet-group',
    {
      name: instanceIdentifier,
      description: instanceIdentifier,
      subnetIds: vpc.subnets.map((subnet) => subnet.id),
      tags: instanceLabels,
    },
    {},
  );
  const password = createRandomPassword(instanceName, {
    length: 32,
    special: false,
  });

  const instance = new aws.rds.Instance(
    instanceName,
    {
      identifier: instanceIdentifier,
      multiAz: false,
      instanceClass: data.instanceClass,
      storageType: 'gp3',
      allocatedStorage: data.storage.allocated,
      maxAllocatedStorage: data.storage.maximum,
      storageEncrypted: true,
      kmsKeyId: kmsKeyId,
      engine: data.engine,
      engineVersion: data.engineVersion,
      username: data.dbAdminUser,
      password: password.password,
      backupRetentionPeriod: data.backupRetention,
      deleteAutomatedBackups: true,
      backupWindow: '20:00-23:30',
      deletionProtection: data.deletionProtection,
      skipFinalSnapshot: !data.deletionProtection,
      finalSnapshotIdentifier: instanceIdentifier,
      copyTagsToSnapshot: true,
      allowMajorVersionUpgrade: false,
      autoMinorVersionUpgrade: true,
      maintenanceWindow: 'Sun:00:00-Sun:06:00',
      iamDatabaseAuthenticationEnabled: false,
      performanceInsightsEnabled: true,
      performanceInsightsRetentionPeriod: 7,
      performanceInsightsKmsKeyId: kmsKeyId,
      enabledCloudwatchLogsExports: undefined,
      monitoringInterval: 0,
      vpcSecurityGroupIds: [securityGroup.id],
      dbSubnetGroupName: dbSubnet.name,
      networkType: 'IPV4', // DUAL is not allowed for public instances
      publiclyAccessible: true,
      tags: instanceLabels,
    },
    {
      ignoreChanges: ['kmsKeyId', 'performanceInsightsKmsKeyId'],
    },
  );

  writeToDoppler(
    `AWS_RDS_${data.name.toUpperCase()}_ENDPOINT`,
    instance.endpoint,
    globalName,
  );

  writeToDoppler(
    `AWS_RDS_${data.name.toUpperCase()}_PORT`,
    instance.port.apply((port) => port.toString()),
    globalName,
  );

  writeToDoppler(
    `AWS_RDS_${data.name.toUpperCase()}_ADMIN_USERNAME`,
    instance.username,
    globalName,
  );

  writeToDoppler(
    `AWS_RDS_${data.name.toUpperCase()}_ADMIN_PASSWORD`,
    password.password,
    globalName,
  );

  return {
    rds: instance,
    password: password,
  };
};
