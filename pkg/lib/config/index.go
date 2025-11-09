package config

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"

	"github.com/muhlba91/muehlbachler-shared-services/pkg/model/config/aws"
)

//nolint:gochecknoglobals // global configuration is acceptable here
var (
	// Environment holds the current deployment environment (e.g., dev, staging, prod).
	Environment string
	// GlobalName is a constant name used across resources.
	GlobalName = "shared-services"
)

// LoadConfig loads the configuration for the given Pulumi context.
// ctx: The Pulumi context.
func LoadConfig(ctx *pulumi.Context) (*aws.Config, error) {
	Environment = ctx.Stack()

	cfg := config.New(ctx, "")

	var awsConfig aws.Config
	cfg.RequireObject("aws", &awsConfig)

	return &awsConfig, nil
}

// CommonLabels returns a map of common labels to be used across resources.
func CommonLabels() map[string]string {
	return map[string]string{
		"environment": Environment,
		"purpose":     GlobalName,
	}
}
