package createwallet

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Yoga-Saputra/go-grpc-boilerplate/app/usecase/v1/wallet"
	"github.com/hibiken/asynq"
)

type (
	// QPayload define payload for create wallet queue.
	QPayload struct {
		BranchID uint16
		MemberID uint64
		PID      string
		Currency string
		Username string
	}
)

// Package constant
const (
	QueueName = "CREATEWALLET"
	TaskName  = QueueName + ":Run"

	maxRetry = 3
)

// Local variable
var (
	CreateWalletTasks = map[string]interface{}{
		TaskName: Handler,
	}

	client *asynq.Client
)

// CreateClient create new asynq client.
func CreateClient(c *asynq.Client) {
	client = c
}

// Enqueue job to the queue system.
func Enqueue(p *QPayload) (*asynq.TaskInfo, error) {
	// Create task
	b, err := json.Marshal(p)
	if err != nil {
		p.logger(fmt.Sprintf("(Init) Failed marshaling the payload: %s", err.Error()))
		return nil, err
	}

	task := asynq.NewTask(TaskName, b)

	// Enqueue job
	info, err := client.Enqueue(
		task,
		asynq.TaskID(p.taskID()),
		asynq.Queue(QueueName),
		asynq.MaxRetry(maxRetry),
	)
	if err != nil {
		p.logger(fmt.Sprintf("(Init) Failed to enqueue task: %s", err.Error()))
		return nil, err
	}

	p.logger(fmt.Sprintf("(Init) Successfully enqueue task: %s", info.ID))
	return info, nil
}

func Handler(c context.Context, t *asynq.Task) error {
	// Prepare the task payload
	var p *QPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %s -> %s", err.Error(), asynq.SkipRetry)
	}

	// create wallet
	if err := wallet.CreateWallet(p.BranchID, p.MemberID, p.PID, p.Currency, p.Username); err != nil {
		p.logger(fmt.Sprintf("(CreateWallet) with p_id %v:  Error: %s", p.PID, err.Error()))
		return err
	}

	p.logger("(CreateWallet) Wallet successfully inserted")
	return nil
}

func (p *QPayload) taskID() string {
	return fmt.Sprintf("%s-%s", p.Currency, p.PID)
}

func (p *QPayload) logger(m string) {
	prefix := fmt.Sprintf("CreateWallet-%v", p.PID)

	log.Printf("[%s] %s", prefix, m)
}
