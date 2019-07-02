package risevision

import (
	"context"
	"fmt"
	apiclient "github.com/blocktree/risevision-adapter/api"
	"github.com/blocktree/openwallet/log"
	"github.com/blocktree/openwallet/openwallet"
	"math/big"
)

type WalletManager struct {
	openwallet.AssetsAdapterBase

	Api             *apiclient.Client               // 节点客户端
	Config          *WalletConfig                   // 节点配置
	Decoder         openwallet.AddressDecoder       //地址编码器
	TxDecoder       openwallet.TransactionDecoder   //交易单编码器
	Log             *log.OWLogger                   //日志工具
	ContractDecoder openwallet.SmartContractDecoder //智能合约解析器
	Blockscanner    *RISEBlockScanner                //区块扫描器
	client          *Client                         //本地封装的http client
	Context         context.Context
}

func NewWalletManager() *WalletManager {
	wm := WalletManager{}
	wm.Config = NewConfig(Symbol)
	wm.Blockscanner = NewRISEBlockScanner(&wm)
	wm.Decoder = NewAddressDecoder(&wm)
	wm.TxDecoder = NewTransactionDecoder(&wm)
	wm.Log = log.NewOWLogger(wm.Symbol())

	wm.Context = context.TODO()
	return &wm
}

type RISEAccount struct {
	Address            string
	Balance            *big.Int
	UnconfirmedBalance *big.Int
}

//GetAccount
func (wm *WalletManager) GetAccount(address string) (*RISEAccount, error) {

	if wm.Api == nil {
		return nil, fmt.Errorf("risevision API is not inited")
	}

	accountReq := &apiclient.AccountRequest{
		Address: address,
		ListOptions: apiclient.ListOptions{
			Limit: 1,
		},
	}
	q, err := wm.Api.GetAccounts(wm.Context, accountReq)
	if err != nil {
		return nil, err
	}
	riseAccount := &RISEAccount{
		Address: address,
	}
	if q.Account != nil  {
		account := q.Account
		riseAccount.Balance = IntToBalance(account.Balance)
		riseAccount.UnconfirmedBalance = IntToBalance(account.UnconfirmedBalance)

	} else {
		riseAccount.Balance = IntToBalance(0)
		riseAccount.UnconfirmedBalance = IntToBalance(0)
	}
	return riseAccount, nil
}

//GetAccountPendingTxCount
func (wm *WalletManager) GetAccountPendingTxCount(address string) (uint64, error) {

	if wm.client == nil {
		return 0, fmt.Errorf("risevision API is not inited")
	}
	//GetPendingAccountTransactionsByPubkey有bug
	//p := external.NewGetPendingAccountTransactionsByPubkeyParams().WithPubkey(address)
	//result, err := wm.Api.Node.External.GetPendingAccountTransactionsByPubkey(p)
	//if err != nil {
	//	return 0, err
	//}

	path := fmt.Sprintf("/accounts/%s/transactions/pending", address)
	result, err := wm.client.Call(path, "GET", nil)
	if err != nil {
		return 0, err
	}

	txs := result.Get("transactions")
	if !txs.IsArray() {
		return 0, nil
	}

	return uint64(len(txs.Array())), nil
}

// BroadcastTransaction recalculates the transaction hash and sends the transaction to the node.
func (wm *WalletManager) BroadcastTransaction(txHex string) (string, error) {
	//txBytes, err := hex.DecodeString(txHex)
	//if err != nil {
	//	return "", fmt.Errorf("transaction decode failed, unexpected error: %v", err)
	//}
	//signedEncodedTx := aeternity.Encode(aeternity.PrefixTransaction, txBytes)
	// calculate the hash of the decoded txRLP
	//rlpTxHashRaw := owcrypt.Hash(txBytes, 32, owcrypt.HASH_ALG_BLAKE2B)
	//// base58/64 encode the hash with the th_ prefix
	//signedEncodedTxHash := aeternity.Encode(aeternity.PrefixTransactionHash, rlpTxHashRaw)

	// send it to the network
	//return postTransaction(wm.Api.Node, signedEncodedTx)
	return "", nil
}

func IntToBalance(balance int64) *big.Int {
	return big.NewInt(balance)
}
