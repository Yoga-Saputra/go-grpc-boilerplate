package config

// Kafka configuration key value
type kafka struct {
	// Kafka server list
	Servers []string `json:"servers" yaml:"servers"`
}
