import { Config, getStack } from '@pulumi/pulumi';

import { AWSData } from '../model/config/aws';
import { MongoDBData } from '../model/config/mongodb';

export const environment = getStack();

const config = new Config();

export const awsConfig = config.requireObject<AWSData>('aws');
export const mongodbConfig = config.requireObject<MongoDBData>('mongodb');

export const globalName = 'shared-services';

export const commonLabels = {
  environment: environment,
  purpose: globalName,
};
