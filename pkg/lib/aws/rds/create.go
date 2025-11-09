package rds

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/muhlba91/pulumi-shared-library/pkg/lib/random"
	"github.com/muhlba91/pulumi-shared-library/pkg/lib/vault/secret"
	"github.com/muhlba91/pulumi-shared-library/pkg/util/metadata"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/ec2"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/kms"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/rds"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"github.com/muhlba91/muehlbachler-shared-services/pkg/lib/aws/vpc"
	"github.com/muhlba91/muehlbachler-shared-services/pkg/lib/config"
	rdsModel "github.com/muhlba91/muehlbachler-shared-services/pkg/model/aws/rds"
	awsConfig "github.com/muhlba91/muehlbachler-shared-services/pkg/model/config/aws"
)

// Create creates an AWS RDS instance with the specified configuration.
// ctx: Pulumi context
// rdsConfig: RDS configuration
//
//nolint:mnd,funlen // magic number is acceptable here for configuration defaults
func Create(
	ctx *pulumi.Context,
	rdsConfig *awsConfig.RDSConfig,
	vpcConfig *awsConfig.VPCConfig,
) (*rdsModel.Data, error) {
	instanceName := fmt.Sprintf("rds-%s-%s", config.GlobalName, rdsConfig.Name)
	instanceIdentifier := fmt.Sprintf("%s-%s", rdsConfig.Name, config.Environment)

	tags := config.CommonLabels()
	tags["Name"] = rdsConfig.Name
	pTags := metadata.LabelsToStringMap(tags)

	kAlias, kErr := kms.LookupAlias(ctx, &kms.LookupAliasArgs{
		Name: "alias/aws/rds",
	}, nil)
	if kErr != nil {
		return nil, kErr
	}
	kmsKeyID := kAlias.Arn

	vpc, vErr := vpc.Create(ctx, instanceName, vpcConfig)
	if vErr != nil {
		return nil, vErr
	}

	sg, sgErr := ec2.NewSecurityGroup(ctx, fmt.Sprintf("%s-security-group", instanceName), &ec2.SecurityGroupArgs{
		Description: pulumi.String(fmt.Sprintf("%s: PostgreSQL", instanceName)),
		VpcId:       vpc.Vpc.ID(),
		Ingress: ec2.SecurityGroupIngressArray{
			ec2.SecurityGroupIngressArgs{
				FromPort:       pulumi.Int(5432),
				ToPort:         pulumi.Int(5432),
				Protocol:       pulumi.String("tcp"),
				CidrBlocks:     pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				Ipv6CidrBlocks: pulumi.StringArray{pulumi.String("::/0")},
			},
		},
		Egress: ec2.SecurityGroupEgressArray{
			ec2.SecurityGroupEgressArgs{
				FromPort:       pulumi.Int(0),
				ToPort:         pulumi.Int(0),
				Protocol:       pulumi.String("-1"),
				CidrBlocks:     pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				Ipv6CidrBlocks: pulumi.StringArray{pulumi.String("::/0")},
			},
		},
		Tags: pTags,
	})
	if sgErr != nil {
		return nil, sgErr
	}

	var subnetIDs pulumi.StringArray
	for _, sn := range vpc.Subnets {
		subnetIDs = append(subnetIDs, sn.ID())
	}

	dbSubnet, sErr := rds.NewSubnetGroup(ctx, fmt.Sprintf("%s-subnet-group", instanceName), &rds.SubnetGroupArgs{
		Name:        pulumi.String(instanceIdentifier),
		Description: pulumi.String(instanceIdentifier),
		SubnetIds:   subnetIDs,
		Tags:        pTags,
	})
	if sErr != nil {
		return nil, sErr
	}

	password, pErr := random.CreatePassword(ctx, fmt.Sprintf("%s-password", instanceName), &random.PasswordOptions{
		Length:  32,
		Special: false,
	})
	if pErr != nil {
		return nil, pErr
	}

	instance, rErr := rds.NewInstance(ctx, instanceName, &rds.InstanceArgs{
		Identifier:                         pulumi.String(instanceIdentifier),
		MultiAz:                            pulumi.Bool(false),
		InstanceClass:                      pulumi.String(rdsConfig.InstanceClass),
		StorageType:                        pulumi.String("gp3"),
		AllocatedStorage:                   pulumi.Int(rdsConfig.Storage.Allocated),
		MaxAllocatedStorage:                pulumi.Int(rdsConfig.Storage.Maximum),
		StorageEncrypted:                   pulumi.Bool(true),
		KmsKeyId:                           pulumi.String(kmsKeyID),
		Engine:                             pulumi.String(rdsConfig.Engine),
		EngineVersion:                      pulumi.String(rdsConfig.EngineVersion),
		Username:                           pulumi.String(rdsConfig.DBAdminUser),
		Password:                           password.Password,
		BackupRetentionPeriod:              pulumi.Int(rdsConfig.BackupRetention),
		DeleteAutomatedBackups:             pulumi.Bool(true),
		BackupWindow:                       pulumi.String("20:00-23:30"),
		DeletionProtection:                 pulumi.Bool(rdsConfig.DeletionProtection),
		SkipFinalSnapshot:                  pulumi.Bool(!rdsConfig.DeletionProtection),
		FinalSnapshotIdentifier:            pulumi.String(instanceIdentifier),
		CopyTagsToSnapshot:                 pulumi.Bool(true),
		AllowMajorVersionUpgrade:           pulumi.Bool(false),
		AutoMinorVersionUpgrade:            pulumi.Bool(true),
		MaintenanceWindow:                  pulumi.String("Sun:00:00-Sun:06:00"),
		IamDatabaseAuthenticationEnabled:   pulumi.Bool(false),
		PerformanceInsightsEnabled:         pulumi.Bool(true),
		PerformanceInsightsRetentionPeriod: pulumi.Int(7),
		PerformanceInsightsKmsKeyId:        pulumi.String(kmsKeyID),
		MonitoringInterval:                 pulumi.Int(0),
		VpcSecurityGroupIds:                pulumi.StringArray{sg.ID()},
		DbSubnetGroupName:                  dbSubnet.Name,
		NetworkType:                        pulumi.String("IPV4"), // DUAL is not allowed for public instances
		PubliclyAccessible:                 pulumi.Bool(true),
		Tags:                               pTags,
	}, pulumi.IgnoreChanges([]string{"kmsKeyId", "performanceInsightsKmsKeyId"}))
	if rErr != nil {
		return nil, rErr
	}

	secretVal, _ := pulumi.All(instance.Endpoint, instance.Port, instance.Username, password.Password).ApplyT(
		func(args []interface{}) (string, error) {
			endpoint := args[0].(string)
			port := args[1].(int)
			username := args[2].(string)
			pass := args[3].(string)

			m := map[string]string{
				"endpoint": endpoint,
				"port":     strconv.Itoa(port),
				"username": username,
				"password": pass,
			}

			b, jErr := json.Marshal(m)
			if jErr != nil {
				return "", jErr
			}
			return string(b), nil
		},
	).(pulumi.StringOutput)

	if _, vsErr := secret.Write(ctx, &secret.WriteArgs{
		Key:   fmt.Sprintf("aws-rds-%s", strings.ToLower(rdsConfig.Name)),
		Value: secretVal,
		Path:  config.GlobalName,
	}); vsErr != nil {
		return nil, vsErr
	}

	return &rdsModel.Data{
		Rds:      instance,
		Password: password,
	}, nil
}
