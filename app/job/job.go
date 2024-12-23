package job

import (
	"github.com/Yoga-Saputra/go-grpc-boilerplate/app/job/createwallet"
	"github.com/Yoga-Saputra/go-grpc-boilerplate/pkg/kafadapter"
	"github.com/hibiken/asynq"
)

type (
	// Meta define meta package that needed by queue job.
	Meta struct {
		Client        *asynq.Client
		KafkaProducer kafadapter.Builder
	}

	// Hold queue and their tasks definition to be registered on service up.
	jobQ struct {
		QueueName string
		Priority  int
		Tasks     map[string]interface{}
	}
)

// Map of function that will be called on Up() method based on their order.
// If have new services, just create new file and their method and register here
var RegiteredTask = []jobQ{
	{
		QueueName: createwallet.QueueName,
		Priority:  1,
		Tasks:     createwallet.CreateWalletTasks,
	},
}

// CreateClient create new asyqn client.
func CreateClient(m Meta) {
	if m.Client == nil {
		panic("asynq client cannot be null")
	}

	// "CreateWallet" queue
	createwallet.CreateClient(m.Client)
}
