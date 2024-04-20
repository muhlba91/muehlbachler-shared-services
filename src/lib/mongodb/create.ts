import * as mongodbatlas from '@pulumi/mongodbatlas';
import { Output } from '@pulumi/pulumi';

import { MongoDBClusterConfig } from '../../model/config/mongodb/cluster';
import { MongoDBClusterData } from '../../model/data/mongodb/atlas';
import { commonLabels, environment, globalName } from '../configuration';
import { writeToDoppler } from '../util/doppler/secret';
import { writeToVault } from '../util/vault/secret';

/**
 * Creates the MongoDB cluster.
 *
 * @param {MongoDBClusterConfig} data the cluster data
 * @returns {Promise<MongoDBClusterData>} the cluster
 */
export const createMongoDBCluster = async (
  data: MongoDBClusterConfig,
): Promise<MongoDBClusterData> => {
  const clusterName = 'mdb-' + globalName + '-' + data.name;
  const clusterIdentifier = data.name + '-' + environment;

  new mongodbatlas.ProjectIpAccessList(
    'mdb-project-ip-access-list-' + data.projectId,
    {
      projectId: data.projectId,
      cidrBlock: '0.0.0.0/0',
    },
    {},
  );

  const cluster = new mongodbatlas.AdvancedCluster(
    clusterName,
    {
      name: clusterIdentifier,
      clusterType: 'REPLICASET',
      projectId: data.projectId,
      replicationSpecs: [
        {
          regionConfigs: [
            {
              backingProviderName: data.provider,
              electableSpecs: {
                instanceSize: data.size,
              },
              providerName: 'TENANT',
              regionName: data.region,
              priority: 7,
            },
          ],
        },
      ],
      tags: Object.entries(commonLabels).map(([key, value]) => ({
        key: key,
        value: value,
      })),
    },
    {},
  );

  cluster.connectionStrings.apply((cs) => {
    writeToDoppler(
      `MONGODB_ATLAS_${data.name.toUpperCase()}_URI`,
      Output.create(cs[0].standard),
      globalName,
    );

    writeToVault(
      `mongodb-atlas-${data.name.toLowerCase()}`,
      Output.create(JSON.stringify({ uri: cs[0].standard })),
      globalName,
    );
  });

  return {
    cluster: cluster,
  };
};
