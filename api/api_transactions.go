package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/blocktree/risevision-adapter/transactions"
)

type (
	// TransactionRequest is the request body for a transaction request
	TransactionRequest struct {
		// ID of the transaction
		ID string
		// RecipientID of the transaction
		RecipientID string
		// RecipientPublicKey of the transaction
		RecipientPublicKey string
		// SenderID of the transaction
		SenderID string
		// SenderPublicKey of the transaction
		SenderPublicKey string
		// BlockID of the transaction
		BlockID string
		// Type of the transaction
		Type *int
		// Height of the transaction
		Height int64
		// MinAmount of the transaction
		MinAmount *int64
		// MaxAmount of the transaction
		MaxAmount *int64

		// FromTimestamp only returns transactions after this time
		FromTimestamp int64
		// ToTimestamp only returns transactions before this time
		ToTimestamp int64
		ListOptions
	}

	// TransactionsResponse is the API response for transaction requests
	TransactionsResponse struct {
		// Transactions are the results
		Transactions []*Transaction `json:"transactions"`
		Success      bool           `json:"success"`
	}

	// TransactionSendResponse is the API response for transaction send requests
	TransactionSendResponse struct {
		Success bool `json:"success"`
	}
)

// GetTransactions searches for transactions on the blockchain.
// Search parameters can be specified in options.
// Limit is set to 100 by default
func (c *Client) GetTransactions(ctx context.Context, options *TransactionRequest) (*TransactionsResponse, error) {
	req := c.restClient.R().SetContext(ctx)

	if options != nil {
		if options.ID != "" {
			req.SetQueryParam("id", options.ID)
		}
		if options.RecipientID != "" {
			req.SetQueryParam("recipientId", options.RecipientID)
		}
		if options.RecipientPublicKey != "" {
			req.SetQueryParam("recipientPublicKey", options.RecipientPublicKey)
		}
		if options.SenderID != "" {
			req.SetQueryParam("senderId", options.SenderID)
		}
		if options.SenderPublicKey != "" {
			req.SetQueryParam("senderPublicKey", options.SenderPublicKey)
		}
		if options.BlockID != "" {
			req.SetQueryParam("blockId", options.BlockID)
		}

		if options.Sort != "" {
			req.SetQueryParam("sort", string(options.Sort))
		}
	}

	req.SetResult(&TransactionsResponse{})
	req.SetError(Error{})

	res, err := req.Get("api/transactions")
	if err != nil {
		return nil, err
	}
	if res.IsError() {
		return res.Result().(*TransactionsResponse), res.Error().(error)
	}

	if !res.Result().(*TransactionsResponse).Success {
		return nil, errors.New("GetTransactions not success")
	}

	return res.Result().(*TransactionsResponse), nil
}

// SendTransaction submits the transaction to the network.
func (c *Client) SendTransaction(ctx context.Context, transaction *transactions.Transaction) (*TransactionSendResponse, error) {
	req := c.restClient.R().SetContext(ctx)

	//posts := make([]*transactions.Transaction, 0)
	//posts = append(posts, transaction)

	//resultByte,_ := transaction.Serialize()
	result,err := json.Marshal(transaction)
	//fmt.Println(hex.EncodeToString(resultByte))


	transactionResult := `{"transaction":` + string(result) + `}`
	fmt.Println(transactionResult)
	req.SetBody(transactionResult)



	req.SetResult(&TransactionSendResponse{})
	req.SetError(Error{})
	req.SetHeader("Content-Type","application/json")
	res, err := req.Put("api/transactions")
	if err != nil {
		return nil, err
	}

	fmt.Println(string(res.Body()))
	if res.IsError() {
		return res.Result().(*TransactionSendResponse), res.Error().(error)
	}

	if !res.Result().(*TransactionSendResponse).Success {
		return nil, errors.New("SendTransaction not success")
	}
	return res.Result().(*TransactionSendResponse), nil
}

// SendTransaction submits the transaction to the network.
func (c *Client) SendSerializableTransaction(ctx context.Context, transaction *transactions.SerializableTransaction) (*TransactionSendResponse, error) {
	req := c.restClient.R().SetContext(ctx)

	req.SetBody(transaction)

	req.SetResult(&TransactionSendResponse{})
	req.SetError(Error{})

	res, err := req.Post("api/transactions")
	if err != nil {
		return nil, err
	}

	return res.Result().(*TransactionSendResponse), res.Error().(error)
}
