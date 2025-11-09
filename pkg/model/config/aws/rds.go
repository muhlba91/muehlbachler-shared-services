package aws

// RDSConfig defines AWS RDS configuration.
type RDSConfig struct {
	// Name is the name of the RDS instance.
	Name string `yaml:"name,omitempty"`
	// VPC is the VPC configuration.
	VPC *VPCConfig `yaml:"vpc,omitempty"`
	// InstanceClass is the instance class.
	InstanceClass string `yaml:"instanceClass,omitempty"`
	// Storage is the storage configuration.
	Storage *RDSStorageConfig `yaml:"storage,omitempty"`
	// Engine is the database engine.
	Engine string `yaml:"engine,omitempty"`
	// EngineVersion is the database engine version.
	EngineVersion string `yaml:"engineVersion,omitempty"`
	// DBAdminUser is the database admin user.
	DBAdminUser string `yaml:"dbAdminUser,omitempty"`
	// BackupRetention is the backup retention period in days.
	BackupRetention int `yaml:"backupRetention,omitempty"`
	// DeletionProtection indicates whether deletion protection is enabled.
	DeletionProtection bool `yaml:"deletionProtection,omitempty"`
}

// RDSStorageConfig defines AWS RDS storage configuration.
type RDSStorageConfig struct {
	// Allocated is the allocated storage in GB.
	Allocated int `yaml:"allocated,omitempty"`
	// Maximum is the maximum storage in GB.
	Maximum int `yaml:"maximum,omitempty"`
}
