package boot

import (
	"context"
	"fmt"
	"time"

	"github.com/Yoga-Saputra/go-grpc-boilerplate/app/job"
	"github.com/Yoga-Saputra/go-grpc-boilerplate/config"
	"github.com/hibiken/asynq"
)

// Queue driver pointer value fo redis adapter
var Queue *asynq.Server

// Start queue server
func queueUp(args *AppArgs) {
	var conn asynq.RedisConnOpt

	if config.Of.Queue.Redis.CSE {
		conn = asynq.RedisClusterClientOpt{
			Addrs:    []string{fmt.Sprintf("%s:%v", config.Of.Queue.Redis.Host, config.Of.Queue.Redis.Port)},
			Password: config.Of.Queue.Redis.Password,
		}
	} else {
		conn = asynq.RedisClientOpt{
			Addr:     fmt.Sprintf("%s:%v", config.Of.Queue.Redis.Host, config.Of.Queue.Redis.Port),
			Password: config.Of.Queue.Redis.Password,
			DB:       config.Of.Queue.Redis.Database,
			PoolSize: config.Of.Queue.Redis.PoolSize,
		}
	}

	// Create queue server
	createServer(conn)

	// Create queue client
	createClient(conn)
}

// Stop queue server
func queueDown() {
	printOutFinishTask("Wait until all queue tasks is finished...")
	printOutFinishTask("and closing current Queue Redis connection...")

	// Shutdown the server
	Queue.Shutdown()
}

// Create queue server
func createServer(conn asynq.RedisConnOpt) {
	// Define asynq config
	aqcfg := asynq.Config{
		Concurrency:         config.Of.Queue.Option.Concurrency,
		ShutdownTimeout:     time.Duration(config.Of.Queue.Option.ShutdownTimeout) * time.Second,
		HealthCheckInterval: time.Duration(config.Of.Queue.Option.HealthCheckInterval) * time.Second,
		RetryDelayFunc: func(n int, e error, t *asynq.Task) time.Duration {
			return 2 * time.Second
		},
		StrictPriority: config.Of.Queue.Option.StrictPriority,
	}

	// Set the queues
	queues := make(map[string]int, len(job.RegiteredTask))
	for _, taq := range job.RegiteredTask {
		queues[taq.QueueName] = taq.Priority
	}
	if len(queues) > 0 {
		aqcfg.Queues = queues
	}

	// Create new queue server
	qSrv := asynq.NewServer(conn, aqcfg)
	Queue = qSrv

	// Create new ServerMux
	mux := asynq.NewServeMux()

	// Assign registered task into mux
	for _, taq := range job.RegiteredTask {
		for tn, hd := range taq.Tasks {
			handler, ok := hd.(func(context.Context, *asynq.Task) error)
			if !ok {
				panic(fmt.Sprintf("mismatch queue task handler: %v", hd))
			}

			mux.HandleFunc(tn, handler)
		}
	}

	// Start Queue Server using goroutine
	go func() {
		if err := qSrv.Start(mux); err != nil {
			panic(err)
		}
	}()

	printOutUp("New Queue Server successfully open")
}

// Create queue client
func createClient(conn asynq.RedisConnOpt) {
	// Asynq client
	client := asynq.NewClient(conn)

	job.CreateClient(job.Meta{
		Client:        client,
		KafkaProducer: KafADPT,
	})
	printOutUp("New Queue Client successfully open")
}
