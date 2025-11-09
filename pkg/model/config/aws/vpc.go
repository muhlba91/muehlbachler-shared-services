package aws

// VPCConfig defines AWS VPC configuration.
type VPCConfig struct {
	// CIDR is the CIDR block for the VPC.
	CIDR string `yaml:"cidr,omitempty"`
	// IPv6Only indicates whether the VPC is IPv6 only.
	IPv6Only bool `yaml:"ipv6Only,omitempty"`
}
