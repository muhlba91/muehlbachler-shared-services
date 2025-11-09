package aws

// Config defines AWS-related configuration.
type Config struct {
	// Postgres contains configuration for PostgreSQL RDS instances.
	Postgres *RDSConfig `yaml:"postgres,omitempty"`
}
