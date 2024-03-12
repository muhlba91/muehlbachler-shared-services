import * as aws from '@pulumi/aws';

import { RandomPasswordData } from '../../random';

/**
 * Defines AWS RDS data.
 */
export type AWSRdsData = {
  readonly rds: aws.rds.Instance;
  readonly password: RandomPasswordData;
};
