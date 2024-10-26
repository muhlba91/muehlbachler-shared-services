# `muehlbachler`: Shared Services

[![Build status](https://img.shields.io/github/actions/workflow/status/muhlba91/muehlbachler-shared-services/pipeline.yml?style=for-the-badge)](https://github.com/muhlba91/muehlbachler-shared-services/actions/workflows/pipeline.yml)
[![License](https://img.shields.io/github/license/muhlba91/muehlbachler-shared-services?style=for-the-badge)](LICENSE.md)
[![](https://api.scorecard.dev/projects/github.com/muhlba91/muehlbachler-shared-services/badge?style=for-the-badge)](https://scorecard.dev/viewer/?uri=github.com/muhlba91/muehlbachler-shared-services)

This repository contains the Shared Services for `muehlbachler` using [Pulumi](http://pulumi.com).

---

## Requirements

- [NodeJS](https://nodejs.org/en), and [yarn](https://yarnpkg.com)
- [Pulumi](https://www.pulumi.com/docs/install/)

## Creating the Infrastructure

To create the services, a [Pulumi Stack](https://www.pulumi.com/docs/concepts/stack/) with the correct configuration needs to exists.

The stack can be deployed via:

```bash
yarn install
yarn build; pulumi up
```

## Destroying the Infrastructure

The entire infrastructure can be destroyed via:

```bash
yarn install
yarn build; pulumi destroy
```

## Environment Variables

To successfully run, and configure the Pulumi plugins, you need to set a list of environment variables. Alternatively, refer to the used Pulumi provider's configuration documentation.

- `CLOUDSDK_CORE_PROJECT`: the Google Cloud (GCP) project
- `CLOUDSDK_COMPUTE_REGION` the Google Cloud (GCP) region
- `GOOGLE_APPLICATION_CREDENTIALS`: reference to a file containing the Google Cloud (GCP) service account credentials
- `AWS_REGION`: the AWS region
- `AWS_ACCESS_KEY_ID`: the AWS access key identifier
- `AWS_SECRET_ACCESS_KEY`: the AWS secret access key

---

## Configuration

The following section describes the configuration which must be set in the Pulumi Stack.

***Attention:*** do use [Secrets Encryption](https://www.pulumi.com/docs/concepts/secrets/#:~:text=Pulumi%20never%20sends%20authentication%20secrets,“secrets”%20for%20extra%20protection.) provided by Pulumi for secret values!

### AWS

The AWS configuration contains all shared services hosten on AWS.

```yaml
aws:
  postgres: the specification of the RDS PostgreSQL instance
    name: the instance name
    backupRetention: the number of days to keep automated backups
    deletionProtection: protect the instance from deletion; if true it will also keep automated backups on delete
    engine: the engine to use (needs to be 'postgres')
    engineVersion: the PostgreSQL engine version to use (see https://docs.aws.amazon.com/cli/latest/reference/rds/describe-db-engine-versions.html for options)
    dbAdminUser: the database admin user (needs to be 'postgres')
    instanceClass: the instance class identifier (e.g. db.t4g.micro)
    storage: the storage specification (autoexpansion enabled)
      allocated: the default allocated storage in GB (min. 20)
      maximum: the maximum storage the instance can expand to
    vpc: VPC data to host the instance in
      ipv6Only: true if the VPC should be IPv6 only - this provides cost savings due to IPv4 charges; not working for RDS!
      cidr: IPv4 CIDR for the VPC/subnet to create
      ipv6Cidr: IPv6 CIDR for the VPC/subnet to create
```

*Note*: to get the latest PostgreSQL engine version for a given PostgreSQL major version you can run `aws rds describe-db-engine-versions --engine postgres --default-only --db-parameter-group-family 'postgres14'`.

---

## Continuous Integration and Automations

- [GitHub Actions](https://docs.github.com/en/actions) are linting, and verifying the code.
- [Renovate Bot](https://github.com/renovatebot/renovate) is updating NodeJS packages, and GitHub Actions.
