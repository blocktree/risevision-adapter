package risevision

import (
	"fmt"
	"github.com/blocktree/risevision-adapter/api"
	"github.com/blocktree/openwallet/common"
	"github.com/blocktree/openwallet/crypto"
	"github.com/blocktree/openwallet/openwallet"
	"math/big"
)

type AddrBalance struct {
	Address      string
	Balance      *big.Int
	TokenBalance *big.Int
}

type txFeeInfo struct {
	GasUsed  *big.Int
	GasPrice *big.Int
	Fee      *big.Int
}

func (f *txFeeInfo) CalcFee() error {
	fee := new(big.Int)
	fee.Mul(f.GasUsed, f.GasPrice)
	f.Fee = fee
	return nil
}

type Block struct {
	Hash              string
	Confirmations     uint64
	Merkleroot        string
	Previousblockhash string
	Height            uint64 `storm:"id"`
	Version           uint64
	Time              uint64
	Fork              bool
}

func NewBlock(generation *api.Block) *Block {
	obj := &Block{}
	obj.Height = uint64(generation.Height)
	obj.Hash = string(generation.ID)
	obj.Previousblockhash = string(generation.PreviousBlockID)
	obj.Time = uint64(generation.Timestamp)

	return obj
}

//BlockHeader 区块链头
func (b *Block) BlockHeader(symbol string) *openwallet.BlockHeader {

	obj := openwallet.BlockHeader{}
	//解析json
	obj.Hash = b.Hash
	//obj.Confirmations = b.Confirmations
	//obj.Merkleroot = b.TransactionMerkleRoot
	obj.Previousblockhash = b.Previousblockhash
	obj.Height = b.Height
	obj.Version = uint64(b.Version)
	obj.Time = b.Time
	obj.Symbol = symbol

	return &obj
}



//UnscanRecords 扫描失败的区块及交易
type UnscanRecord struct {
	ID          string `storm:"id"` // primary key
	BlockHeight uint64
	BlockID     string
	Reason      string
}

func NewUnscanRecord(height uint64, blockID, reason string) *UnscanRecord {
	obj := UnscanRecord{}
	obj.BlockHeight = height
	obj.BlockID = blockID
	obj.Reason = reason
	obj.ID = common.Bytes2Hex(crypto.SHA256([]byte(fmt.Sprintf("%d_%s", height, blockID))))
	return &obj
}
