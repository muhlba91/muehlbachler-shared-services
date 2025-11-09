package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"github.com/muhlba91/muehlbachler-shared-services/pkg/lib/aws/rds"
	"github.com/muhlba91/muehlbachler-shared-services/pkg/lib/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		awsConfig, err := config.LoadConfig(ctx)
		if err != nil {
			return err
		}

		awsRdsPostgresql, err := rds.Create(ctx, awsConfig.Postgres, awsConfig.Postgres.VPC)
		if err != nil {
			return err
		}

		ctx.Export("aws", pulumi.Map{
			"postgresql": pulumi.Map{
				"id":       awsRdsPostgresql.Rds.ID(),
				"address":  awsRdsPostgresql.Rds.Address,
				"port":     awsRdsPostgresql.Rds.Port,
				"endpoint": awsRdsPostgresql.Rds.Endpoint,
				"username": awsRdsPostgresql.Rds.Username,
				"password": awsRdsPostgresql.Password.Password,
			},
		})

		return nil
	})
}
