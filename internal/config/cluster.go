package config

// DefaultPFAddress specifies the default PortForward host address.
const DefaultPFAddress = "localhost"

// Cluster tracks W2 cluster configuration.
type Cluster struct {
	Namespace          *Namespace    `yaml:"namespace"`
	View               *View         `yaml:"view"`
	FeatureGates       *FeatureGates `yaml:"featureGates"`
	PortForwardAddress string        `yaml:"portForwardAddress"`
}

// NewCluster creates a new cluster configuration.
func NewCluster() *Cluster {
	return &Cluster{
		Namespace:          NewNamespace(),
		View:               NewView(),
		PortForwardAddress: DefaultPFAddress,
		FeatureGates:       NewFeatureGates(),
	}
}
