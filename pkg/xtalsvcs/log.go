// Transfer Log service
package xtalsvcs

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Yoga-Saputra/go-grpc-boilerplate/pkg/httpclient"
)

// TransactionLogCentral defines Transaction Log Service
type TransactionLogCentral struct {
	TicketID string `json:"ticketId"`

	BranchCode int    `json:"branchCode"`
	pID        string `json:"pID"`
	Username   string `json:"username"`
	From       string `json:"from"`
	Currency   string `json:"currency"`

	Amount        float64 `json:"amount"`
	BalanceBefore float64 `json:"balanceBefore"`
	BalanceAfter  float64 `json:"balanceAfter"`

	TransactionStatus string    `json:"transactionStatus"`
	TransactionType   string    `json:"transactiontype"`
	Description       string    `json:"description"`
	TransactionDate   time.Time `json:"transactionDate"`
}

// TransferLogResponse defines response from Transaction Log service
type TransferLogResponse struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Error   interface{} `json:"error"`
}

const (
	xferLogSrvKey   string = "log"          // <- Key name of Transfer Log service in config/env file
	xferLogEpInsert string = "/v2/transfer" // <- Endpoint insert
)

// Insert inserting log to Transaction Log service
func Insert(p *TransactionLogCentral) (
	*TransferLogResponse,
	error,
) {
	defer func() {
		p = nil
	}() // <- Freed pointer

	// Is config exists?
	if err := checkConfigExists(xferLogSrvKey); err != nil {
		return nil, err
	}

	// Prepare options
	uri := getConfigValue(xferLogSrvKey, "host") + xferLogEpInsert
	token := getConfigValue(xferLogSrvKey, "token")

	// New http client
	hc := httpclient.New(httpclient.Config{
		Headers: map[string]interface{}{"Content-Type": "application/json; charset=utf-8"},
	})
	if len(token) > 0 {
		hc.Token(token)
	}
	hc.SetTimeout(2)

	// Payloads
	var payloads map[string]interface{}
	b, err := json.Marshal(&p)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, &payloads); err != nil {
		return nil, err
	}

	// Do request
	resp, err := hc.Post(uri, payloads)
	if err != nil {
		log.Printf("[TransactionLogCentral] - Err: %s", err.Error())

		return nil, err
	} else if resp.StatusCode != 200 && resp.StatusCode != 201 {
		err := fmt.Errorf("error with http code %v", resp.StatusCode)
		log.Printf("[TransactionLogCentral] - Err: %s", err.Error())

		return nil, err
	}
	log.Printf("[TransactionLogCentral] - Latency: %s", resp.Latency)

	// Unmarshal data
	var result TransferLogResponse
	if err := json.Unmarshal(resp.BodyBytes, &result); err != nil {
		err = fmt.Errorf("invalid response json: %s", err.Error())
		log.Printf("[TransactionLogCentral] - Err: %s", err.Error())

		return nil, err
	}

	return &result, nil
}
