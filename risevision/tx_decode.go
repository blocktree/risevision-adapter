package risevision

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/blocktree/go-owcrypt"
	"github.com/blocktree/risevision-adapter/crypto"
	"github.com/blocktree/risevision-adapter/risevision_txsigner"
	"github.com/blocktree/risevision-adapter/transactions"
	"github.com/blocktree/openwallet/common"
	"github.com/blocktree/openwallet/log"
	"github.com/blocktree/openwallet/openwallet"
	"github.com/shopspring/decimal"
	"math/big"
	"sort"
	"time"
)

type TransactionDecoder struct {
	openwallet.TransactionDecoderBase
	wm *WalletManager //钱包管理者
}

//NewTransactionDecoder 交易单解析器
func NewTransactionDecoder(wm *WalletManager) *TransactionDecoder {
	decoder := TransactionDecoder{}
	decoder.wm = wm
	return &decoder
}


//CreateRawTransaction 创建交易单
func (decoder *TransactionDecoder) CreateRawTransaction(wrapper openwallet.WalletDAI, rawTx *openwallet.RawTransaction) error {

	var (
		decimals        = decoder.wm.Decimal()
		accountID       = rawTx.Account.AccountID
		fixFees         = big.NewInt(0)
		findAddrBalance *AddrBalance
		feeInfo         *txFeeInfo
	)

	//获取wallet
	addresses, err := wrapper.GetAddressList(0, -1, "AccountID", accountID) //wrapper.GetWallet().GetAddressesByAccount(rawTx.Account.AccountID)
	if err != nil {
		return err
	}

	if len(addresses) == 0 {
		return fmt.Errorf("[%s] have not addresses", accountID)
	}

	searchAddrs := make([]string, 0)
	for _, address := range addresses {
		searchAddrs = append(searchAddrs, address.Address)
	}

	addrBalanceArray, err := decoder.wm.Blockscanner.GetBalanceByAddress(searchAddrs...)
	if err != nil {
		return err
	}

	var amountStr string
	for _, v := range rawTx.To {
		amountStr = v
		break
	}

	//地址余额从大到小排序
	sort.Slice(addrBalanceArray, func(i int, j int) bool {
		a_amount, _ := decimal.NewFromString(addrBalanceArray[i].Balance)
		b_amount, _ := decimal.NewFromString(addrBalanceArray[j].Balance)
		if a_amount.LessThan(b_amount) {
			return true
		} else {
			return false
		}
	})

	amount := common.StringNumToBigIntWithExp(amountStr, decimals)

	if len(rawTx.FeeRate) > 0 {
		fixFees = common.StringNumToBigIntWithExp(rawTx.FeeRate, decimals)
	} else {
		fixFees = common.StringNumToBigIntWithExp(decoder.wm.Config.FixFees, decimals)
	}

	//计算手续费
	feeInfo = &txFeeInfo{
		Fee:      fixFees,
		GasPrice: fixFees,
		GasUsed:  big.NewInt(1),
	}

	for _, addrBalance := range addrBalanceArray {

		addrBalance_BI := common.StringNumToBigIntWithExp(addrBalance.Balance, decimals)

		//总消耗数量 = 转账数量 + 手续费
		totalAmount := new(big.Int)
		totalAmount.Add(amount, feeInfo.Fee)

		//余额不足查找下一个地址
		if addrBalance_BI.Cmp(totalAmount) < 0 {
			continue
		}

		//只要找到一个合适使用的地址余额就停止遍历
		findAddrBalance = &AddrBalance{Address: addrBalance.Address, Balance: addrBalance_BI}
		break
	}

	if findAddrBalance == nil {
		return fmt.Errorf("all address's balance of account is not enough")
	}

	//最后创建交易单
	err = decoder.createRawTransaction(
		wrapper,
		rawTx,
		findAddrBalance,
		feeInfo,
		"")
	if err != nil {
		return err
	}

	return nil

}

//createRawTransaction
func (decoder *TransactionDecoder) createRawTransaction(
	wrapper openwallet.WalletDAI,
	rawTx *openwallet.RawTransaction,
	addrBalance *AddrBalance,
	feeInfo *txFeeInfo,
	callData string) error {

	var (
		accountTotalSent = decimal.Zero
		txFrom           = make([]string, 0)
		txTo             = make([]string, 0)
		keySignList      = make([]*openwallet.KeySignature, 0)
		amountStr        string
		destination      string
	)

	decimals := int32(0)
	if rawTx.Coin.IsContract {
		decimals = int32(rawTx.Coin.Contract.Decimals)
	} else {
		decimals = decoder.wm.Decimal()
	}
	//isContract := rawTx.Coin.IsContract
	//contractAddress := rawTx.Coin.Contract.Address
	//tokenCoin := rawTx.Coin.Contract.Token
	//tokenDecimals := int(rawTx.Coin.Contract.Decimals)
	//coinDecimals := this.wm.Decimal()

	for k, v := range rawTx.To {
		destination = k
		amountStr = v
		break
	}

	//计算账户的实际转账amount
	accountTotalSentAddresses, findErr := wrapper.GetAddressList(0, -1, "AccountID", rawTx.Account.AccountID, "Address", destination)
	if findErr != nil || len(accountTotalSentAddresses) == 0 {
		amountDec, _ := decimal.NewFromString(amountStr)
		accountTotalSent = accountTotalSent.Add(amountDec)
	}

	txFrom = []string{fmt.Sprintf("%s:%s", addrBalance.Address, amountStr)}
	txTo = []string{fmt.Sprintf("%s:%s", destination, amountStr)}

	addr, err := wrapper.GetAddress(addrBalance.Address)
	if err != nil {
		return err
	}

	//decoder.wm.Log.Debugf("nonce: %d", nonce)
	//decoder.wm.Log.Debugf("pending: %d", pending)

	amount := common.StringNumToBigIntWithExp(amountStr, decimals)

	timestamp := transactions.GetCurrentTimeWithOffset(0)

	decoder.wm.Log.Debugf("timestamp: %d", timestamp)

	senderPk, err := hex.DecodeString(addr.PublicKey)
	if err != nil {
		return err
	}
	// create the SpendTransaction
	transaction := &transactions.Transaction{
		Type:            0,
		Amount:          amount.Uint64(),
		RecipientID:     destination,
		Timestamp:       timestamp,
		SenderPublicKey: senderPk,
	}





	txRaw, err := transaction.MarshalJSON()
	if err != nil {
		return err
	}

	rawTx.RawHex = string(txRaw)


	hash, err := transaction.Hash()
	if err != nil {
		return err
	}

	txHex := string(hash)
	log.Warn(crypto.GetBigNumberStringFromBytes(crypto.GetFirstEightBytesReversed(hash)))
	if rawTx.Signatures == nil {
		rawTx.Signatures = make(map[string][]*openwallet.KeySignature)
	}

	msg := append([]byte(decoder.wm.Config.NetworkID), txHex...)

	signature := openwallet.KeySignature{
		EccType: decoder.wm.Config.CurveType,
		Address: addr,
		Message: hex.EncodeToString(msg),
	}
	keySignList = append(keySignList, &signature)

	feesAmount := common.BigIntToDecimals(feeInfo.Fee, decimals)
	gasPrice := common.BigIntToDecimals(feeInfo.GasPrice, decimals)
	accountTotalSent = accountTotalSent.Add(feesAmount)
	accountTotalSent = decimal.Zero.Sub(accountTotalSent)

	//rawTx.RawHex = rawHex
	rawTx.Signatures[rawTx.Account.AccountID] = keySignList
	rawTx.FeeRate = gasPrice.String()
	rawTx.Fees = feesAmount.String()
	rawTx.IsBuilt = true
	rawTx.TxAmount = accountTotalSent.StringFixed(decimals)
	rawTx.TxFrom = txFrom
	rawTx.TxTo = txTo


	return nil
}

//SignRawTransaction 签名交易单
func (decoder *TransactionDecoder) SignRawTransaction(wrapper openwallet.WalletDAI, rawTx *openwallet.RawTransaction) error {

	if rawTx.Signatures == nil || len(rawTx.Signatures) == 0 {
		//this.wm.Log.Std.Error("len of signatures error. ")
		return fmt.Errorf("transaction signature is empty")
	}

	key, err := wrapper.HDKey()
	if err != nil {
		return err
	}

	keySignatures := rawTx.Signatures[rawTx.Account.AccountID]
	if keySignatures != nil {
		for _, keySignature := range keySignatures {

			childKey, err := key.DerivedKeyWithPath(keySignature.Address.HDPath, keySignature.EccType)

			keyBytes, err := childKey.GetPrivateKeyBytes()

			if err != nil {
				return err
			}

			msg ,_ :=  hex.DecodeString(keySignature.Message)
			//if err != nil {
			//	return fmt.Errorf("decoder transaction hash failed, unexpected err: %v", err)
			//}

			sig, err := risevision_txsigner.Default.SignTransactionHash(msg, keyBytes, owcrypt.ECC_CURVE_ED25519)

			//sig2,_ := risevision_txsigner.Default.SignTransactionHash(msg, keyBytes, keySignature.EccType)


			if err != nil {

				return fmt.Errorf("sign transaction hash failed, unexpected err: %v", err)
			}

			keySignature.Signature = hex.EncodeToString(sig)
		}
	}

	decoder.wm.Log.Info("transaction hash sign success")

	rawTx.Signatures[rawTx.Account.AccountID] = keySignatures

	return nil
}

//VerifyRawTransaction 验证交易单，验证交易单并返回加入签名后的交易单
func (decoder *TransactionDecoder) VerifyRawTransaction(wrapper openwallet.WalletDAI, rawTx *openwallet.RawTransaction) error {

	if rawTx.Signatures == nil || len(rawTx.Signatures) == 0 {
		//this.wm.Log.Std.Error("len of signatures error. ")
		return fmt.Errorf("serializableTransaction signature is empty")
	}
	//
	var serializableTransaction *transactions.SerializableTransaction
	transJson := []byte(rawTx.RawHex)
	err := json.Unmarshal(transJson, &serializableTransaction)
	if err != nil {
		return fmt.Errorf("serializableTransaction decode failed, unexpected error: %v", err)
	}

	//支持多重签名
	for accountID, keySignatures := range rawTx.Signatures {
		decoder.wm.Log.Debug("accountID Signatures:", accountID)
		for _, keySignature := range keySignatures {

			msg ,_ :=  hex.DecodeString(keySignature.Message)
			signature, _ := hex.DecodeString(keySignature.Signature)
			publicKey, _ := hex.DecodeString(keySignature.Address.PublicKey)

			//decoder.wm.Log.Debug("txHex:", hex.EncodeToString(txHex))
			//decoder.wm.Log.Debug("Signature:", keySignature.Signature)

			//验证签名
			ret := owcrypt.Verify(publicKey, nil, 0, msg, uint16(len(msg)), signature,owcrypt.ECC_CURVE_ED25519)
			if ret != owcrypt.SUCCESS {
				return fmt.Errorf("serializableTransaction verify failed")
			}

			serializableTransaction.Signature = keySignature.Signature

			txRaw, err := json.Marshal(serializableTransaction)
			if err != nil {
				return fmt.Errorf("serializableTransaction verify Marshal failed,err: %v", err)
			}

			rawTx.IsCompleted = true
			rawTx.RawHex = string(txRaw)
			break

		}
	}

	return nil
}

//SendRawTransaction 广播交易单
func (decoder *TransactionDecoder) SubmitRawTransaction(wrapper openwallet.WalletDAI, rawTx *openwallet.RawTransaction) (*openwallet.Transaction, error) {

	var serializableTransaction *transactions.SerializableTransaction
	transJson := []byte(rawTx.RawHex)
	err := json.Unmarshal(transJson, &serializableTransaction)
	if err != nil {
		return nil, fmt.Errorf("serializableTransaction decode failed, unexpected error: %v", err)
	}



	senderPk, err := hex.DecodeString(serializableTransaction.SenderPublicKey)
	if err != nil {
		return nil,err
	}
	sig, err := hex.DecodeString(serializableTransaction.Signature)
	if err != nil {
		return nil,err
	}
	// create the SpendTransaction
	transaction := &transactions.Transaction{
		Type:            0,
		Amount:         serializableTransaction.Amount,
		RecipientID:     serializableTransaction.RecipientID,
		Timestamp:       serializableTransaction.Timestamp,
		SenderPublicKey: senderPk,
		Signature:sig,
	}


	resp, err := decoder.wm.Api.SendTransaction(decoder.wm.Context, transaction)
	if err != nil {
		return nil, err
	}

	decoder.wm.Log.Info("Transaction [%s] submitted to the network successfully.", resp.Success)




	rawTx.TxID,err = transaction.ID()
	if err != nil {
		return nil, err
	}

	rawTx.IsSubmit = true

	decimals := decoder.wm.Decimal()

	//记录一个交易单
	tx := &openwallet.Transaction{
		From:       rawTx.TxFrom,
		To:         rawTx.TxTo,
		Amount:     rawTx.TxAmount,
		Coin:       rawTx.Coin,
		TxID:       rawTx.TxID,
		Decimal:    decimals,
		AccountID:  rawTx.Account.AccountID,
		Fees:       rawTx.Fees,
		SubmitTime: time.Now().Unix(),
	}

	tx.WxID = openwallet.GenTransactionWxID(tx)

	return tx, nil
}

//GetRawTransactionFeeRate 获取交易单的费率
func (decoder *TransactionDecoder) GetRawTransactionFeeRate() (feeRate string, unit string, err error) {
	return "0.1", "RISE", nil
}

//CreateSummaryRawTransaction 创建汇总交易
func (decoder *TransactionDecoder) CreateSummaryRawTransaction(wrapper openwallet.WalletDAI, sumRawTx *openwallet.SummaryRawTransaction) ([]*openwallet.RawTransaction, error) {

	var (
		decimals        = decoder.wm.Decimal()
		rawTxArray      = make([]*openwallet.RawTransaction, 0)
		accountID       = sumRawTx.Account.AccountID
		minTransfer     = common.StringNumToBigIntWithExp(sumRawTx.MinTransfer, decimals)
		retainedBalance = common.StringNumToBigIntWithExp(sumRawTx.RetainedBalance, decimals)
		fixFees         = big.NewInt(0)
		feeInfo         *txFeeInfo
	)

	if minTransfer.Cmp(retainedBalance) < 0 {
		return nil, fmt.Errorf("mini transfer amount must be greater than address retained balance")
	}

	//获取wallet
	addresses, err := wrapper.GetAddressList(sumRawTx.AddressStartIndex, sumRawTx.AddressLimit,
		"AccountID", sumRawTx.Account.AccountID)
	if err != nil {
		return nil, err
	}

	if len(addresses) == 0 {
		return nil, fmt.Errorf("[%s] have not addresses", accountID)
	}

	searchAddrs := make([]string, 0)
	for _, address := range addresses {
		searchAddrs = append(searchAddrs, address.Address)
	}

	addrBalanceArray, err := decoder.wm.Blockscanner.GetBalanceByAddress(searchAddrs...)
	if err != nil {
		return nil, err
	}

	if len(sumRawTx.FeeRate) > 0 {
		fixFees = common.StringNumToBigIntWithExp(sumRawTx.FeeRate, decimals)
	} else {
		fixFees = common.StringNumToBigIntWithExp(decoder.wm.Config.FixFees, decimals)
	}

	//计算手续费
	//计算手续费
	feeInfo = &txFeeInfo{
		Fee: fixFees,
		GasPrice: fixFees,
		GasUsed: big.NewInt(1),
	}

	for _, addrBalance := range addrBalanceArray {

		//检查余额是否超过最低转账
		addrBalance_BI := common.StringNumToBigIntWithExp(addrBalance.Balance, decimals)

		if addrBalance_BI.Cmp(minTransfer) < 0 || addrBalance_BI.Cmp(big.NewInt(0)) <= 0 {
			continue
		}
		//计算汇总数量 = 余额 - 保留余额
		sumAmount_BI := new(big.Int)
		sumAmount_BI.Sub(addrBalance_BI, retainedBalance)

		//减去手续费
		sumAmount_BI.Sub(sumAmount_BI, feeInfo.Fee)
		if sumAmount_BI.Cmp(big.NewInt(0)) <= 0 {
			continue
		}

		sumAmount := common.BigIntToDecimals(sumAmount_BI, decimals)
		feesAmount := common.BigIntToDecimals(feeInfo.Fee, decimals)

		decoder.wm.Log.Debugf("balance: %v", addrBalance.Balance)
		decoder.wm.Log.Debugf("fees: %v", feesAmount)
		decoder.wm.Log.Debugf("sumAmount: %v", sumAmount)

		//创建一笔交易单
		rawTx := &openwallet.RawTransaction{
			Coin:    sumRawTx.Coin,
			Account: sumRawTx.Account,
			To: map[string]string{
				sumRawTx.SummaryAddress: sumAmount.StringFixed(decoder.wm.Decimal()),
			},
			Required: 1,
		}

		createErr := decoder.createRawTransaction(
			wrapper,
			rawTx,
			&AddrBalance{Address: addrBalance.Address, Balance: addrBalance_BI},
			feeInfo,
			"")
		if createErr != nil {
			return nil, createErr
		}

		//创建成功，添加到队列
		rawTxArray = append(rawTxArray, rawTx)

	}

	return rawTxArray, nil

}


//CreateSummaryRawTransactionWithError 创建汇总交易，返回能原始交易单数组（包含带错误的原始交易单）
func (decoder *TransactionDecoder) CreateSummaryRawTransactionWithError(wrapper openwallet.WalletDAI, sumRawTx *openwallet.SummaryRawTransaction) ([]*openwallet.RawTransactionWithError, error) {
	raTxWithErr := make([]*openwallet.RawTransactionWithError, 0)
	rawTxs, err := decoder.CreateSummaryRawTransaction(wrapper, sumRawTx)
	if err != nil {
		return nil, err
	}
	for _, tx := range rawTxs {
		raTxWithErr = append(raTxWithErr, &openwallet.RawTransactionWithError{
			RawTx: tx,
			Error: nil,
		})
	}
	return raTxWithErr, nil
}