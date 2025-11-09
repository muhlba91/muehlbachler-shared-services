package rds

import (
	"github.com/muhlba91/pulumi-shared-library/pkg/model/random"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/rds"
)

// Data defines AWS RDS data.
type Data struct {
	// Rds is the RDS instance.
	Rds *rds.Instance
	// Password is the generated password.
	Password *random.PasswordData
}
