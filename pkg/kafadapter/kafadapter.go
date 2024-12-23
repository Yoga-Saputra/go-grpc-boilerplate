// Package Kafka Adpater
package kafadapter

import "errors"

type Messages struct {
	Key   []byte
	Value []byte
}

type (
	Builder interface {
		NewProducer() error
		Publish(topic string, messages []Messages) (err error)
		CloseProducer()
	}

	adapter struct {
		Builder
	}
)

var (
	// Error config brokers must be set
	ErrConfigBrokersRequired = errors.New("config brokers is required")

	// Error producer pointer is nil
	ErrProducerIsNil = errors.New("producer is not defined")
)

// Build kafka adapter based on given builder.
func Build(b Builder) (*adapter, error) {
	if err := b.NewProducer(); err != nil {
		return nil, err
	}
	return &adapter{b}, nil
}
