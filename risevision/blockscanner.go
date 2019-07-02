package risevision

import (
	"fmt"
	"github.com/blocktree/risevision-adapter/api"
	"github.com/blocktree/openwallet/common"
	"github.com/blocktree/openwallet/log"
	"github.com/blocktree/openwallet/openwallet"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"math/big"
	"time"
)

const (
	blockchainBucket = "blockchain" // blockchain dataset
	//periodOfTask      = 5 * time.Second // task interval
	maxExtractingSize = 0 // thread count
	successTxType     = 0
)

//AEBlockScanner RISE block scanner
type RISEBlockScanner struct {
	*openwallet.BlockScannerBase

	CurrentBlockHeight   uint64         //当前区块高度
	extractingCH         chan struct{}  //扫描工作令牌
	wm                   *WalletManager //钱包管理者
	RescanLastBlockCount uint64         //重扫上N个区块数量
}

//ExtractResult extract result
type ExtractTxResult struct {
	extractData map[string]*openwallet.TxExtractData
	TxID        string
}

//ExtractResult extract result
type ExtractResult struct {
	extractData []*ExtractTxResult
	BlockHash   string
	BlockHeight uint64
	BlockTime   int64
	Success     bool
}

//SaveResult result
type SaveResult struct {
	TxID        string
	BlockHeight uint64
	Success     bool
}

// NewEOSBlockScanner create a block scanner
func NewRISEBlockScanner(wm *WalletManager) *RISEBlockScanner {
	bs := RISEBlockScanner{
		BlockScannerBase: openwallet.NewBlockScannerBase(),
	}

	bs.extractingCH = make(chan struct{}, 3)
	bs.wm = wm

	bs.RescanLastBlockCount = maxExtractingSize

	// set task
	bs.SetTask(bs.ScanBlockTask)

	return &bs
}

//GetCurrentBlock 获取当前最新区块
func (bs *RISEBlockScanner) GetCurrentBlock() (*api.BlockResponse, error) {
	request := &api.BlockRequest{
		ListOptions: api.ListOptions{
			Limit: 1,
		},
	}
	resp, err := bs.wm.Api.GetBlocks(bs.wm.Context, request)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

//GetBlockHeight 获取区块链高度
func (bs *RISEBlockScanner) GetGlobalMaxBlockHeight() uint64 {

	resp, err := bs.GetCurrentBlock()
	if err != nil || resp == nil {
		log.Errorf("resp is nil,err is : ", err)
		return 0
	}
	blocks := resp.Blocks
	if len(blocks) <= 0 {
		log.Errorf("blocks length is 0")
		return 0
	}
	block := blocks[0]
	return uint64(block.Height)
}

func (bs *RISEBlockScanner) GetBlockByHeight(height uint64) (*Block, error) {
	br := &api.BlockRequest{
		Height: int64(height),
	}
	keyBlock, err := bs.wm.Api.GetBlocks(bs.wm.Context, br)
	if err != nil {
		return nil, err
	}
	blocks := keyBlock.Blocks
	if blocks == nil || len(blocks) == 0 {
		return nil, errors.New("blocks is nil or blocks length is 0")
	}

	blockTemp := blocks[0]
	block := NewBlock(blockTemp)
	return block, nil
}

//GetScannedBlockHeader 获取当前扫描的区块头
func (bs *RISEBlockScanner) GetScannedBlockHeader() (*openwallet.BlockHeader, error) {

	var (
		blockHeader *openwallet.BlockHeader
		blockHeight uint64 = 0
		hash        string
		err         error
	)

	blockHeight, hash = bs.wm.GetLocalNewBlock()

	//如果本地没有记录，查询接口的高度
	if blockHeight == 0 {
		blockHeader, err = bs.GetCurrentBlockHeader()
		if err != nil {

			return nil, err
		}
		blockHeight = blockHeader.Height
		//就上一个区块链为当前区块
		blockHeight = blockHeight - 1

		block, err := bs.GetBlockByHeight(blockHeight)
		if err != nil {
			return nil, err
		}
		hash = block.Hash
	}

	return &openwallet.BlockHeader{Height: blockHeight, Hash: hash}, nil
}

//GetCurrentBlockHeader 获取当前区块高度
func (bs *RISEBlockScanner) GetCurrentBlockHeader() (*openwallet.BlockHeader, error) {

	br := &api.BlockRequest{
		ListOptions: api.ListOptions{
			Limit: 0,
		},
	}
	keyBlock, err := bs.wm.Api.GetBlocks(bs.wm.Context, br)
	if err != nil {
		return nil, err
	}
	blocks := keyBlock.Blocks
	if blocks == nil || len(blocks) == 0 {
		return nil, errors.New("blocks is nil")
	}
	block := blocks[0]

	return &openwallet.BlockHeader{Height: uint64(block.Height), Hash: string(block.ID)}, nil
}

//GetBalanceByAddress 查询地址余额
func (bs *RISEBlockScanner) GetBalanceByAddress(address ...string) ([]*openwallet.Balance, error) {

	addrBalanceArr := make([]*openwallet.Balance, 0)
	for _, a := range address {
		acc, err := bs.wm.GetAccount(a)
		if err == nil {

			b := common.BigIntToDecimals(acc.Balance, bs.wm.Decimal())
			ub := common.BigIntToDecimals(acc.UnconfirmedBalance, bs.wm.Decimal())
			obj := &openwallet.Balance{
				Symbol:           bs.wm.Symbol(),
				Address:          a,
				Balance:          ub.String(),
				UnconfirmBalance: ub.String(),
				ConfirmBalance:   b.String(),
			}

			//log.Warn("address:",a,",ubalance:",ub.String(),"confirmBlance:",b.String())
			addrBalanceArr = append(addrBalanceArr, obj)
			//return nil, err
		}

	}

	return addrBalanceArr, nil
}

//GetScannedBlockHeight 获取已扫区块高度
func (bs *RISEBlockScanner) GetScannedBlockHeight() uint64 {
	localHeight, _ := bs.wm.GetLocalNewBlock()
	return localHeight
}

//GetTransactionsByBlockHash
func (bs *RISEBlockScanner) GetTransactionsByBlock(block *Block) ([]*api.Transaction, error) {

	transReq := &api.TransactionRequest{
		BlockID: block.Hash,
	}

	trans, err := bs.wm.Api.GetTransactions(bs.wm.Context, transReq)

	if err != nil {
		return nil, err
	}

	if trans.Transactions == nil {

		return nil, errors.New("trans.Transactions is nil")
	}

	return trans.Transactions, nil
}

//GetBlockHeight 获取区块链高度
func (bs *RISEBlockScanner) GetBlockHeight() (uint64, error) {
	//暂时只有一个接口获取钱包高度
	return bs.GetGlobalMaxBlockHeight(), nil
}

//newBlockNotify 获得新区块后，通知给观测者
func (bs *RISEBlockScanner) newBlockNotify(block *Block, isFork bool) {
	header := block.BlockHeader(bs.wm.Symbol())
	header.Fork = isFork
	bs.NewBlockNotify(header)
}

//extractRuntime 提取运行时
func (bs *RISEBlockScanner) extractRuntime(producer chan ExtractResult, worker chan ExtractResult, quit chan struct{}) {

	var (
		values = make([]ExtractResult, 0)
	)

	for {

		var activeWorker chan<- ExtractResult
		var activeValue ExtractResult

		//当数据队列有数据时，释放顶部，传输给消费者
		if len(values) > 0 {
			activeWorker = worker
			activeValue = values[0]

		}

		select {

		//生成者不断生成数据，插入到数据队列尾部
		case pa := <-producer:
			values = append(values, pa)
		case <-quit:
			//退出
			return
		case activeWorker <- activeValue:
			//wm.Log.Std.Info("Get %d", len(activeValue))
			values = values[1:]
		}
	}

}

//ExtractTransaction 提取交易单
func (bs *RISEBlockScanner) ExtractTransaction(block *Block, scanTargetFunc openwallet.BlockScanTargetFunc) (ExtractResult, error) {

	result := ExtractResult{
		Success:     true,
		BlockHeight: block.Height,
		extractData: make([]*ExtractTxResult, 0),
	}
	transactionList := make([]*api.Transaction, 0)

		req := &api.TransactionRequest{
			BlockID: block.Hash,
		}

		txInfos, err := bs.wm.Api.GetTransactions(bs.wm.Context, req)
		if err != nil {
			log.Errorf("cant find the transaction by height %d ,err : %s", block.Height, err.Error())
		}

		if txInfos != nil {
			transactionList = append(transactionList, txInfos.Transactions...)
		}


	if len(transactionList) != 0 {
		for _, v := range transactionList {
			resultTx, err := bs.changeTrans(v, scanTargetFunc)
			if err != nil {
				bs.wm.Log.Std.Error("trans ID: %d, save unscan record failed. unexpected error: %v", v.ID, err.Error())
				continue
			}
			result.extractData = append(result.extractData, &resultTx)
		}

	}

	return result, nil

}

//ExtractTransaction 提取交易单
func (bs *RISEBlockScanner) ExtractTransactionSingleTx(block *Block, tx *api.Transaction, scanTargetFunc openwallet.BlockScanTargetFunc) (ExtractTxResult, error) {
	return bs.changeTrans(tx, scanTargetFunc)

}

func (bs *RISEBlockScanner) changeTrans(trans *api.Transaction, scanTargetFunc openwallet.BlockScanTargetFunc) (ExtractTxResult, error) {
	v := trans
	resultTx := ExtractTxResult{
		TxID:        trans.ID,
		extractData: make(map[string]*openwallet.TxExtractData),
	}
	var (
		txID     = string(trans.ID)
		createAt = time.Now().Unix()
		txType   = v.Type
		decimals = bs.wm.Decimal()
	)

	switch txType {
	case successTxType:
		amountDec := decimal.New(v.Amount,0)

		feeDec := decimal.New(v.Fee,0)

		amount := common.BigIntToDecimals(big.NewInt(amountDec.IntPart()), decimals).String()
		fees := common.BigIntToDecimals(big.NewInt(feeDec.IntPart()), decimals).String()
		from := string(v.SenderID)
		to := string(v.RecipientID)

		sourceKey, ok := scanTargetFunc(
			openwallet.ScanTarget{
				Address:          from,
				BalanceModelType: openwallet.BalanceModelTypeAddress,
			})
		if ok {
			input := openwallet.TxInput{}
			input.TxID = txID
			input.Address = string(v.SenderID)
			input.Amount = amount
			input.Coin = openwallet.Coin{
				Symbol:     bs.wm.Symbol(),
				IsContract: false,
			}
			input.Index = 0
			input.Sid = openwallet.GenTxInputSID(txID, bs.wm.Symbol(), "", uint64(0))
			//input.CreateAt = createAt
			input.BlockHeight = uint64(v.Height)
			//input.BlockHash = string(trx.BlockHash)
			input.BlockHash = trans.BlockID //TODO: 先记录keyblock的hash方便上层计算确认次数，以后做扩展
			ed := resultTx.extractData[sourceKey]
			if ed == nil {
				ed = openwallet.NewBlockExtractData()
				resultTx.extractData[sourceKey] = ed
			}

			ed.TxInputs = append(ed.TxInputs, &input)

			//手续费也作为一个输出
			tmp := *&input
			feeCharge := &tmp
			feeCharge.Amount = fees
			ed.TxInputs = append(ed.TxInputs, feeCharge)
		}

		sourceKey2, ok2 := scanTargetFunc(
			openwallet.ScanTarget{
				Address:          string(v.RecipientID),
				BalanceModelType: openwallet.BalanceModelTypeAddress,
			})
		if ok2 {
			output := openwallet.TxOutPut{}
			output.TxID = txID
			output.Address = to
			output.Amount = amount
			output.Coin = openwallet.Coin{
				Symbol:     bs.wm.Symbol(),
				IsContract: false,
			}
			output.Index = 0
			output.Sid = openwallet.GenTxOutPutSID(txID, bs.wm.Symbol(), "", 0)
			output.CreateAt = createAt
			output.BlockHeight = uint64(v.Height)
			//output.BlockHash = string(trx.BlockHash)
			output.BlockHash = trans.BlockID //TODO: 先记录keyblock的hash方便上层计算确认次数，以后做扩展
			ed := resultTx.extractData[sourceKey2]
			if ed == nil {
				ed = openwallet.NewBlockExtractData()
				resultTx.extractData[sourceKey2] = ed
			}

			ed.TxOutputs = append(ed.TxOutputs, &output)
		}

		for _, extractData := range resultTx.extractData {
			status := "1"
			reason := ""
			tx := &openwallet.Transaction{
				From:   []string{from + ":" + amount},
				To:     []string{to + ":" + amount},
				Amount: amount,
				Fees:   fees,
				Coin: openwallet.Coin{
					Symbol:     bs.wm.Symbol(),
					IsContract: false,
				},
				//BlockHash:   string(trx.BlockHash),
				BlockHash:   trans.BlockID, //TODO: 先记录keyblock的hash方便上层计算确认次数，以后做扩展
				BlockHeight: uint64(v.Height),
				TxID:        txID,
				Decimal:     decimals,
				Status:      status,
				Reason:      reason,
				Confirm:     int64(trans.Confirmations),
				//SubmitTime:  int64(block.Time),
				ConfirmTime: int64(trans.Timestamp),
			}
			wxID := openwallet.GenTransactionWxID(tx)
			tx.WxID = wxID
			extractData.Transaction = tx

		}
	default:
		return resultTx, nil
	}

	return resultTx, nil
}

//newExtractDataNotify 发送通知
func (bs *RISEBlockScanner) newExtractDataNotify(height uint64, extractTxResult []*ExtractTxResult) error {

	for o, _ := range bs.Observers {
		for _, txResult := range extractTxResult {
			for key, data := range txResult.extractData {
				err := o.BlockExtractDataNotify(key, data)
				if err != nil {
					bs.wm.Log.Error("BlockExtractDataNotify unexpected error:", err)
					//记录未扫区块
					unscanRecord := NewUnscanRecord(height, "", "ExtractData Notify failed.")
					err = bs.SaveUnscanRecord(unscanRecord)
					if err != nil {
						bs.wm.Log.Std.Error("block height: %d, save unscan record failed. unexpected error: %v", height, err.Error())
					}

				}
			}
		}
	}

	return nil
}

//BatchExtractTransaction 批量提取交易单
//bitcoin 1M的区块链可以容纳3000笔交易，批量多线程处理，速度更快
func (bs *RISEBlockScanner) BatchExtractTransaction(block *Block) error {

	var (
		quit   = make(chan struct{})
		failed = 0
	)

	//生产通道
	producer := make(chan ExtractResult)
	defer close(producer)

	//消费通道
	worker := make(chan ExtractResult)
	defer close(worker)

	//保存工作
	saveWork := func(height uint64, result chan ExtractResult) {
		//回收创建的地址
		for gets := range result {

			if gets.Success {

				notifyErr := bs.newExtractDataNotify(height, gets.extractData)
				//saveErr := bs.SaveRechargeToWalletDB(height, gets.Recharges)
				if notifyErr != nil {
					failed++ //标记保存失败数
					bs.wm.Log.Std.Info("newExtractDataNotify unexpected error: %v", notifyErr)
				}

			} else {
				//记录未扫区块
				unscanRecord := NewUnscanRecord(height, "", "")
				bs.SaveUnscanRecord(unscanRecord)
				bs.wm.Log.Std.Info("block height: %d extract failed.", height)
				failed++ //标记保存失败数
			}
			//累计完成的线程数
			close(quit) //关闭通道，等于给通道传入nil
		}
	}

	//提取工作
	extractWork := func(eBlock *Block, eProducer chan ExtractResult) {
		bs.extractingCH <- struct{}{}
		//shouldDone++
		go func(mBlock *Block, end chan struct{}, mProducer chan<- ExtractResult) {
			result, err := bs.ExtractTransaction(mBlock, bs.ScanTargetFunc)
			if err != nil {
				log.Error("extractWork err :", err)
			}
			//导出提出的交易
			mProducer <- result
			//释放
			<-end

		}(eBlock, bs.extractingCH, eProducer)
	}

	/*	开启导出的线程	*/

	//独立线程运行消费
	go saveWork(block.Height, worker)

	//独立线程运行生产
	go extractWork(block, producer)

	//以下使用生产消费模式
	bs.extractRuntime(producer, worker, quit)

	if failed > 0 {
		return fmt.Errorf("block scanner saveWork failed")
	} else {
		return nil
	}
}

//rescanFailedRecord 重扫失败记录
func (bs *RISEBlockScanner) RescanFailedRecord() {

	list, err := bs.wm.GetUnscanRecords()
	if err != nil {
		bs.wm.Log.Std.Info("block scanner can not get rescan data; unexpected error: %v", err)
	}

	if list == nil || len(list) == 0 {
		bs.wm.Log.Std.Info("block scanner can not get rescan data; list is nil")
		return
	}

	for _, l := range list {

		if l.BlockHeight == 0 {
			continue
		}

		bs.wm.Log.Std.Info("block scanner rescanning height: %d ...", l.BlockHeight)

		block, err := bs.GetBlockByHeight(l.BlockHeight)
		if err != nil {
			bs.wm.Log.Std.Info("block scanner can not get new block data; unexpected error: %v", err)
			continue
		}

		err = bs.BatchExtractTransaction(block)
		if err != nil {
			bs.wm.Log.Std.Info("block scanner can not extractRechargeRecords; unexpected error: %v", err)
			continue
		}

		//删除未扫记录
		bs.wm.DeleteUnscanRecord(l.BlockHeight)
	}

	//删除未没有找到交易记录的重扫记录
	bs.wm.DeleteUnscanRecordNotFindTX()
}

//RISEBlockScanner 扫描任务
func (bs *RISEBlockScanner) ScanBlockTask() {

	//获取本地区块高度
	blockHeader, err := bs.GetScannedBlockHeader()
	if err != nil {
		bs.wm.Log.Std.Info("block scanner can not get new block height; unexpected error: %v", err)
		return
	}

	currentHeight := blockHeader.Height
	currentHash := blockHeader.Hash

	for {

		if !bs.Scanning {
			//区块扫描器已暂停，马上结束本次任务
			return
		}

		//获取最大高度
		maxHeight, err := bs.GetBlockHeight()
		if err != nil {
			//下一个高度找不到会报异常
			bs.wm.Log.Std.Info("block scanner can not get rpc-server block height; unexpected error: %v", err)
			break
		}

		//是否已到最新高度
		if currentHeight >= maxHeight {
			bs.wm.Log.Std.Info("block scanner has scanned full chain data. Current height: %d", maxHeight)
			break
		}

		//继续扫描下一个区块
		currentHeight = currentHeight + 1

		bs.wm.Log.Std.Info("block scanner scanning height: %d ...", currentHeight)

		//获取最大高度
		block, err := bs.GetBlockByHeight(currentHeight)
		if err != nil {
			//记录未扫区块
			unscanRecord := NewUnscanRecord(currentHeight, "", err.Error())
			bs.SaveUnscanRecord(unscanRecord)
			bs.wm.Log.Std.Info("block height: %d extract failed.", currentHeight)
			return
		}

		hash := block.Hash

		isFork := false

		//判断hash是否上一区块的hash
		if currentHash != block.Previousblockhash {

			bs.wm.Log.Std.Info("block has been fork on height: %d.", currentHeight)
			bs.wm.Log.Std.Info("block height: %d local hash = %s ", currentHeight-1, currentHash)
			bs.wm.Log.Std.Info("block height: %d mainnet hash = %s ", currentHeight-1, block.Previousblockhash)

			bs.wm.Log.Std.Info("delete recharge records on block height: %d.", currentHeight-1)

			//查询本地分叉的区块
			forkBlock, _ := bs.wm.GetLocalBlock(currentHeight - 1)

			//删除上一区块链的所有充值记录
			//bs.DeleteRechargesByHeight(currentHeight - 1)
			//删除上一区块链的未扫记录
			bs.wm.DeleteUnscanRecord(currentHeight - 1)
			currentHeight = currentHeight - 2 //倒退2个区块重新扫描
			if currentHeight <= 0 {
				currentHeight = 1
			}

			localBlock, err := bs.wm.GetLocalBlock(currentHeight)
			if err != nil {
				bs.wm.Log.Std.Warning("block scanner can not get local block; unexpected error: %v", err)

				//查找core钱包的RPC
				bs.wm.Log.Info("block scanner prev block height:", currentHeight)

				localBlock, err = bs.GetBlockByHeight(currentHeight)
				if err != nil {
					bs.wm.Log.Std.Error("block scanner can not get prev block; unexpected error: %v", err)
					break
				}

			}

			//重置当前区块的hash
			currentHash = localBlock.Hash

			bs.wm.Log.Std.Info("rescan block on height: %d, hash: %s .", currentHeight, currentHash)

			//重新记录一个新扫描起点
			bs.wm.SaveLocalNewBlock(localBlock.Height, localBlock.Hash)

			isFork = true

			if forkBlock != nil {

				//通知分叉区块给观测者，异步处理
				bs.newBlockNotify(forkBlock, isFork)
			}

		} else {

			err = bs.BatchExtractTransaction(block)
			if err != nil {
				bs.wm.Log.Std.Info("block scanner can not extractRechargeRecords; unexpected error: %v", err)
			}

			//重置当前区块的hash
			currentHash = hash

			//保存本地新高度
			bs.wm.SaveLocalNewBlock(currentHeight, currentHash)
			bs.wm.SaveLocalBlock(block)

			isFork = false

			//通知新区块给观测者，异步处理
			bs.newBlockNotify(block, isFork)
		}

	}

	//重扫前N个块，为保证记录找到
	for i := currentHeight - bs.RescanLastBlockCount; i < currentHeight; i++ {
		bs.scanBlock(i + 1)
	}

	//重扫失败区块
	bs.RescanFailedRecord()

}

//ScanBlock 扫描指定高度区块
func (bs *RISEBlockScanner) ScanBlock(height uint64) error {

	block, err := bs.scanBlock(height)
	if err != nil {
		return err
	}

	//通知新区块给观测者，异步处理
	bs.newBlockNotify(block, false)

	return nil
}

//ScanBlock 扫描指定高度区块
func (bs *RISEBlockScanner) scanBlock(height uint64)(*Block, error) {

	block, err := bs.GetBlockByHeight(height)
	if err != nil {
		return nil,err
	}
	err = bs.BatchExtractTransaction(block)
	if err != nil {
		bs.wm.Log.Std.Info("block scanner can not extractRechargeRecords; unexpected error: %v", err)
		return nil,err
	}
	return block,nil
}


//SetRescanBlockHeight 重置区块链扫描高度
func (bs *RISEBlockScanner) SetRescanBlockHeight(height uint64) error {
	height = height - 1
	if height < 0 {
		return fmt.Errorf("block height to rescan must greater than 0.")
	}
	block, err := bs.GetBlockByHeight(height)
	if err != nil {
		return err
	}

	bs.wm.SaveLocalNewBlock(height, block.Hash)

	return nil
}


//ExtractTransactionData
func (bs *RISEBlockScanner) ExtractTransactionData(txid string, scanAddressFunc openwallet.BlockScanTargetFunc) (map[string][]*openwallet.TxExtractData, error) {
	transReq := &api.TransactionRequest{
		ID: txid,
	}

	trans, err := bs.wm.Api.GetTransactions(bs.wm.Context, transReq)

	if err != nil {
		return nil, err
	}

	if trans.Transactions == nil || len(trans.Transactions) == 0 {

		return nil, errors.New("trans.Transactions is nil")
	}
	tx := trans.Transactions[0]
	block, err := bs.GetBlockByHeight(uint64(tx.Height))
	if err != nil {
		return nil, err
	}
	result, err := bs.ExtractTransactionSingleTx(block, tx, scanAddressFunc)
	if err != nil {
		return nil, err
	}
	extData := make(map[string][]*openwallet.TxExtractData)
	for key, data := range result.extractData {
		txs := extData[key]
		if txs == nil {
			txs = make([]*openwallet.TxExtractData, 0)
		}
		txs = append(txs, data)
		extData[key] = txs
	}
	return extData, nil
}
