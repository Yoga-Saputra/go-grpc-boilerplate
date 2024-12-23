package config

// Queue configuration key value
// Using redis as queue driver
type queue struct {
	// Queue option
	Option queueOption `json:"option" yaml:"option"`

	// Queue driver connection
	Redis redis `json:"redis" yaml:"redis"`
}

type queueOption struct {
	// Maximum number of concurrent processing of tasks.
	Concurrency int

	// ShutdownTimeout specifies the duration to wait to let workers finish their tasks
	// before forcing them to abort when stopping the server.
	ShutdownTimeout int

	// HealthCheckInterval specifies the interval between healthchecks.
	HealthCheckInterval int

	// for set priority level of queues
	StrictPriority bool
}
