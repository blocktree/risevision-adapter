package api

import (
	"context"
	"github.com/pkg/errors"
	"strconv"
)

type (
	// BlockRequest is the request body to request blocks
	BlockRequest struct {
		// BlockID of the block
		BlockID string
		// Height if the block
		Height int64
		// GeneratorPublicKey is the public key of the delegate that forged the block
		GeneratorPublicKey string

		ListOptions
	}

	// BlockResponse is the API response for block requests
	BlockResponse struct {
		// Blocks are the results
		Blocks []*Block `json:"blocks"`

		Success bool `json:"success"`
	}

	// Block is a Lisk block
	Block struct {
		// ID of the block
		ID string `json:"id"`
		// Version of the block
		Version int `json:"version"`
		// Height of the block
		Height int `json:"height"`
		// Timestamp of the block
		Timestamp int `json:"timestamp"`
		// GeneratorPublicKey of the block
		GeneratorPublicKey string `json:"generatorPublicKey"`
		// OptionsLength of the block
		OptionsLength int `json:"optionsLength"`
		// OptionsHash of the block
		OptionsHash string `json:"optionsHash"`
		// BlockSignature of the block
		BlockSignature string `json:"blockSignature"`
		// Confirmations of the block
		Confirmations int `json:"confirmations"`
		// PreviousBlockID of the block
		PreviousBlockID string `json:"previousBlock"`
		// NumberOfTransactions of the block
		NumberOfTransactions int `json:"numberOfTransactions"`
		// TotalAmount of the block
		TotalAmount int64 `json:"totalAmount"`
		// TotalFee of the block
		TotalFee int64 `json:"totalFee"`
		// Reward of the block
		Reward int64 `json:"reward"`
		// TotalForged of the block
		TotalForged string `json:"totalForged"`
	}
)

// GetBlocks searches for blocks on the blockchain.
// Search parameters can be specified in options.
// Limit is set to 100 by default
func (c *Client) GetBlocks(ctx context.Context, options *BlockRequest) (*BlockResponse, error) {
	req := c.restClient.R().SetContext(ctx)

	if options != nil {
		if options.BlockID != "" {
			req.SetQueryParam("blockId", options.BlockID)
		}
		if options.Height != 0 {
			req.SetQueryParam("height", strconv.FormatInt(options.Height, 10))
		}

	}

	req.SetResult(&BlockResponse{})
	req.SetError(Error{})

	res, err := req.Get("api/blocks")
	if err != nil {
		return nil, err
	}

	if res.IsError(){
		return res.Result().(*BlockResponse), res.Error().(error)
	}

	if !res.Result().(*BlockResponse).Success{
		return nil, errors.New("GetBlocks not success")
	}

	return res.Result().(*BlockResponse),nil
}
