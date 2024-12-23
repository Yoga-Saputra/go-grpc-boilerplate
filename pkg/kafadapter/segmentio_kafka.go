// Package Kafka Adpater
package kafadapter

import (
	"context"
	"errors"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/compress"
)

type (
	SegmentioKafka struct {
		Brokers []string
	}
)

var segmentioProducer *kafka.Writer

// NewProducer create new kafka producer.
func (k *SegmentioKafka) NewProducer() error {
	w := &kafka.Writer{
		Addr:        kafka.TCP(k.Brokers...),
		Compression: compress.Gzip,

		AllowAutoTopicCreation: true,
	}

	segmentioProducer = w
	return nil
}

// Publish message to given kafka topic.
func (k *SegmentioKafka) Publish(topic string, messages []Messages) (err error) {
	// Validate pointer of producer
	if segmentioProducer == nil {
		return ErrProducerIsNil
	}

	// Convert message from adapter to segmentio kafka message.
	var kafkaMsgs []kafka.Message
	for _, m := range messages {
		kafkaMsgs = append(kafkaMsgs, kafka.Message{Topic: topic, Key: m.Key, Value: m.Value})
	}

	// Do attempt to publish the message
	retries := 3
	for i := 0; i < retries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Do publish message
		err = segmentioProducer.WriteMessages(ctx, kafkaMsgs...)

		// Break iteration if produce message return without any error
		if err == nil {
			break
		}

		// Do retries until n... if produce message got "LeaderNotAvailable" or "DeadlineExceeded" error
		// otherwise return error without retry
		if errors.Is(err, kafka.LeaderNotAvailable) || errors.Is(err, context.DeadlineExceeded) {
			time.Sleep(250 * time.Millisecond)
			continue
		} else {
			return err
		}
	}

	return nil
}

// CloseProducer closing kafka connection.
func (k *SegmentioKafka) CloseProducer() {
	if segmentioProducer != nil {
		segmentioProducer.Close()
	}
}
