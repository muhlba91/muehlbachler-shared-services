import * as mongodbatlas from '@pulumi/mongodbatlas';

/**
 * Defines MongoDB Cluster data.
 */
export type MongoDBClusterData = {
  readonly cluster: mongodbatlas.AdvancedCluster;
};
