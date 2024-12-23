package config

// Database configuration key value
type database struct {
	// Host name where the Database is hosted
	Host string `json:"host" yaml:"host"`

	// Port number of Database connection
	Port int `json:"port" yaml:"port"`

	// User name of Database connection
	User string `json:"user" yaml:"user"`

	// Password of Database conenction
	Password string `json:"password" yaml:"password"`

	// Name of Database that want to connect
	Name string `json:"name" yaml:"name"`

	// Dialect is varian or type of database query language
	Dialect string `json:"dialect" yaml:"dialect"`

	// Identifier represent custom identifier of this connection
	Identifier string `json:"identifier" yaml:"identifier"`

	// Database resolver
	// Can be called as other database sources
	Resolver resolver `json:"resolver" yaml:"resolver"`

	// Another DB tool
	// Can be called to using specific DB connection to handle some tools
	Tools []database `json:"tools" yaml:"tools"`
}

// Database resolver, both another sources and replicas
type resolver struct {
	// Array of other database sources
	Sources []database `json:"sources" yaml:"sources"`

	// Array of other database that used as replications
	Replicas []database `json:"replicas" yaml:"replicas"`
}
