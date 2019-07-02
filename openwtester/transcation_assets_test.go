/*
 * Copyright 2018 The openwallet Authors
 * This file is part of the openwallet library.
 *
 * The openwallet library is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The openwallet library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.
 */

package openwtester

import (
	"github.com/blocktree/openwallet/openw"
	"testing"

	"github.com/blocktree/openwallet/log"
	"github.com/blocktree/openwallet/openwallet"
)

func testGetAssetsAccountBalance(tm *openw.WalletManager, walletID, accountID string) {
	balance, err := tm.GetAssetsAccountBalance(testApp, walletID, accountID)
	if err != nil {
		log.Error("GetAssetsAccountBalance failed, unexpected error:", err)
		return
	}
	log.Info("balance:", balance)
}

func testGetAssetsAccountTokenBalance(tm *openw.WalletManager, walletID, accountID string, contract openwallet.SmartContract) {
	balance, err := tm.GetAssetsAccountTokenBalance(testApp, walletID, accountID, contract)
	if err != nil {
		log.Error("GetAssetsAccountTokenBalance failed, unexpected error:", err)
		return
	}
	log.Info("token balance:", balance.Balance)
}

func testCreateTransactionStep(tm *openw.WalletManager, walletID, accountID, to, amount, feeRate string, contract *openwallet.SmartContract) (*openwallet.RawTransaction, error) {

	//err := tm.RefreshAssetsAccountBalance(testApp, accountID)
	//if err != nil {
	//	log.Error("RefreshAssetsAccountBalance failed, unexpected error:", err)
	//	return nil, err
	//}

	rawTx, err := tm.CreateTransaction(testApp, walletID, accountID, amount, to, feeRate, "", contract)

	if err != nil {
		log.Error("CreateTransaction failed, unexpected error:", err)
		return nil, err
	}

	return rawTx, nil
}

func testCreateSummaryTransactionStep(
	tm *openw.WalletManager,
	walletID, accountID, summaryAddress, minTransfer, retainedBalance, feeRate string,
	start, limit int,
	contract *openwallet.SmartContract) ([]*openwallet.RawTransactionWithError, error) {

	rawTxArray, err := tm.CreateSummaryRawTransactionWithError(testApp, walletID, accountID, summaryAddress, minTransfer,
		retainedBalance, feeRate, start, limit, contract, nil)

	if err != nil {
		log.Error("CreateSummaryTransaction failed, unexpected error:", err)
		return nil, err
	}

	return rawTxArray, nil
}

func testSignTransactionStep(tm *openw.WalletManager, rawTx *openwallet.RawTransaction) (*openwallet.RawTransaction, error) {

	_, err := tm.SignTransaction(testApp, rawTx.Account.WalletID, rawTx.Account.AccountID, "12345678", rawTx)
	if err != nil {
		log.Error("SignTransaction failed, unexpected error:", err)
		return nil, err
	}

	log.Infof("rawTx: %+v", rawTx)
	return rawTx, nil
}

func testVerifyTransactionStep(tm *openw.WalletManager, rawTx *openwallet.RawTransaction) (*openwallet.RawTransaction, error) {

	//log.Info("rawTx.Signatures:", rawTx.Signatures)

	_, err := tm.VerifyTransaction(testApp, rawTx.Account.WalletID, rawTx.Account.AccountID, rawTx)
	if err != nil {
		log.Error("VerifyTransaction failed, unexpected error:", err)
		return nil, err
	}

	log.Infof("rawTx: %+v", rawTx)
	return rawTx, nil
}

func testSubmitTransactionStep(tm *openw.WalletManager, rawTx *openwallet.RawTransaction) (*openwallet.RawTransaction, error) {

	tx, err := tm.SubmitTransaction(testApp, rawTx.Account.WalletID, rawTx.Account.AccountID, rawTx)
	if err != nil {
		log.Error("SubmitTransaction failed, unexpected error:", err)
		return nil, err
	}

	log.Std.Info("tx: %+v", tx)
	log.Info("wxID:", tx.WxID)
	log.Info("txID:", rawTx.TxID)

	return rawTx, nil
}

func TestTransfer_RISE(t *testing.T) {
	tm := testInitWalletManager()
	walletID := "WCoZhSoyQqmm6EJz4DeKQQeQLU7RZb6u1J"
	accountID := "49m6MqzUhvjvrbMm64pqzrNrSKe4uAx7myUPTs8S5MoZ"
	to := "4748974534685299179R"

	//accountID := "3biDqABdY66PtH5R8xk2upxz5xcLLpu7pDvpq3Bb9Yec"
	//to := "7061235840388139206L"

	testGetAssetsAccountBalance(tm, walletID, accountID)

	for i:= 0 ;i<1 ;i++{
		rawTx, err := testCreateTransactionStep(tm, walletID, accountID, to, "1", "", nil)
		if err != nil {
			return
		}

		log.Std.Info("rawTx: %+v", rawTx)

		_, err = testSignTransactionStep(tm, rawTx)
		if err != nil {
			return
		}

		_, err = testVerifyTransactionStep(tm, rawTx)
		if err != nil {
			return
		}

		_, err = testSubmitTransactionStep(tm, rawTx)
		if err != nil {
			return
		}
	}


}

func TestSummary_RISE(t *testing.T) {
	tm := testInitWalletManager()
	walletID := "WCoZhSoyQqmm6EJz4DeKQQeQLU7RZb6u1J"
	accountID := "49m6MqzUhvjvrbMm64pqzrNrSKe4uAx7myUPTs8S5MoZ"
	summaryAddress := "12832451703076591480R"

	testGetAssetsAccountBalance(tm, walletID, accountID)

	for i:=0;i<3;i++{
		rawTxArray, err := testCreateSummaryTransactionStep(tm, walletID, accountID,
			summaryAddress, "", "", "0.1",
			0, 100, nil)
		if err != nil {
			log.Errorf("CreateSummaryTransaction failed, unexpected error: %v", err)
			return
		}

		//执行汇总交易
		for _, rawTxWithErr := range rawTxArray {

			if rawTxWithErr.Error != nil {
				log.Error(rawTxWithErr.Error.Error())
				continue
			}

			_, err = testSignTransactionStep(tm, rawTxWithErr.RawTx)
			if err != nil {
				return
			}

			_, err = testVerifyTransactionStep(tm, rawTxWithErr.RawTx)
			if err != nil {
				return
			}

			_, err = testSubmitTransactionStep(tm, rawTxWithErr.RawTx)
			if err != nil {
				return
			}
		}
	}

}
