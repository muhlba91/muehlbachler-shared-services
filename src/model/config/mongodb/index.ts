import { MongoDBClusterConfig } from './cluster';

/**
 * Defines MongoDB data.
 */
export type MongoDBData = {
  readonly cluster: MongoDBClusterConfig;
};
