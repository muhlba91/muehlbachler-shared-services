import { createRDSInstance } from './lib/aws/rds/create';
import { awsConfig, mongodbConfig } from './lib/configuration';
import { createMongoDBCluster } from './lib/mongodb/create';

export = async () => {
  const awsRdsPostgresql = await createRDSInstance(awsConfig.postgres);
  const mongodbCluster = await createMongoDBCluster(mongodbConfig.cluster);

  return {
    aws: {
      postgresql: {
        id: awsRdsPostgresql.rds.id,
        address: awsRdsPostgresql.rds.address,
        port: awsRdsPostgresql.rds.port,
        endpoint: awsRdsPostgresql.rds.endpoint,
        username: awsRdsPostgresql.rds.username,
        password: awsRdsPostgresql.password.password,
      },
    },
    mongodb: {
      atlas: mongodbCluster.cluster.connectionStrings.apply(
        (connectionStrings) => ({
          id: mongodbCluster.cluster.clusterId,
          projectId: mongodbCluster.cluster.projectId,
          name: mongodbConfig.cluster.name,
          clusterName: mongodbCluster.cluster.name,
          endpoint: connectionStrings[0].standard,
        }),
      ),
    },
  };
};
