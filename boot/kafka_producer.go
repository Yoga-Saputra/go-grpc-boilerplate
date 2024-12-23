package boot

import (
	"github.com/Yoga-Saputra/go-grpc-boilerplate/config"
	"github.com/Yoga-Saputra/go-grpc-boilerplate/pkg/kafadapter"
)

// Kafka Adapter pointer value
var KafADPT kafadapter.Builder

// Start kafka producer
func kafkaProducerUp(args *AppArgs) {
	adapter, err := kafadapter.Build(&kafadapter.SegmentioKafka{
		Brokers: config.Of.Kafka.Servers,
	})
	if err != nil {
		panic(err)
	}

	KafADPT = adapter
	printOutUp("New Kafka producer successfully created")
}

func kafkaProducerDown() {
	printOutDown("Closing current kafka producer...")
	KafADPT.CloseProducer()
}
