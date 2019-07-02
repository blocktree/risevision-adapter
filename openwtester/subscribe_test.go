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
	"github.com/astaxie/beego/config"
	"github.com/blocktree/openwallet/log"
	"github.com/blocktree/openwallet/openw"
	"github.com/blocktree/openwallet/openwallet"
	"path/filepath"
	"testing"
)

////////////////////////// 测试单个扫描器 //////////////////////////

type subscriberSingle struct {
}

//BlockScanNotify 新区块扫描完成通知
func (sub *subscriberSingle) BlockScanNotify(header *openwallet.BlockHeader) error {
	log.Notice("header:", header)
	return nil
}

//BlockTxExtractDataNotify 区块提取结果通知
func (sub *subscriberSingle) BlockExtractDataNotify(sourceKey string, data *openwallet.TxExtractData) error {
	log.Notice("account:", sourceKey)

	for i, input := range data.TxInputs {
		log.Std.Notice("data.TxInputs[%d]: %+v", i, input)
	}

	for i, output := range data.TxOutputs {
		log.Std.Notice("data.TxOutputs[%d]: %+v", i, output)
	}

	log.Std.Notice("data.Transaction: %+v", data.Transaction)

	return nil
}

func TestSubscribeAddress(t *testing.T) {

	var (
		endRunning = make(chan bool, 1)
		symbol     = "RISE"
		addrs      = map[string]string{
			"17894957342263131751R": "sender",
		}
	)

	//GetSourceKeyByAddress 获取地址对应的数据源标识
	scanTargetFunc := func(target openwallet.ScanTarget) (string, bool) {
		//如果余额模型是地址，查找地址表
		if target.BalanceModelType == openwallet.BalanceModelTypeAddress {
			key, ok := addrs[target.Address]
			if !ok {
				return "", false
			}
			return key, true
		} else {
			//如果余额模型是账户，用别名操作账户的别名
			key, ok := addrs[target.Alias]
			if !ok {
				return "", false
			}
			return key, true
		}

	}

	assetsMgr, err := openw.GetAssetsAdapter(symbol)
	if err != nil {
		log.Error(symbol, "is not support")
		return
	}

	//读取配置
	absFile := filepath.Join(configFilePath, symbol+".ini")

	c, err := config.NewConfig("ini", absFile)
	if err != nil {
		return
	}
	assetsMgr.LoadAssetsConfig(c)

	assetsLogger := assetsMgr.GetAssetsLogger()
	if assetsLogger != nil {
		assetsLogger.SetLogFuncCall(true)
	}

	//log.Debug("already got scanner:", assetsMgr)
	scanner := assetsMgr.GetBlockScanner()
	scanner.SetRescanBlockHeight(2111247)

	if scanner == nil {
		log.Error(symbol, "is not support block scan")
		return
	}

	scanner.SetBlockScanTargetFunc(scanTargetFunc)

	sub := subscriberSingle{}
	scanner.AddObserver(&sub)

	scanner.Run()

	<-endRunning
}

func TestBlockScanner_ExtractTransactionData(t *testing.T) {

	var (
		symbol = "RISE"
		txid   = "6891565370527412645"
		addrs  = map[string]string{
			"9789709449210829854L": "sender",
		}
	)

	//GetSourceKeyByAddress 获取地址对应的数据源标识
	scanTargetFunc := func(target openwallet.ScanTarget) (string, bool) {
		key, ok := addrs[target.Address]
		if !ok {
			return "", false
		}
		return key, true
	}

	assetsMgr, err := openw.GetAssetsAdapter(symbol)
	if err != nil {
		log.Error(symbol, "is not support")
		return
	}

	//读取配置
	absFile := filepath.Join(configFilePath, symbol+".ini")

	c, err := config.NewConfig("ini", absFile)
	if err != nil {
		return
	}
	assetsMgr.LoadAssetsConfig(c)

	assetsLogger := assetsMgr.GetAssetsLogger()
	if assetsLogger != nil {
		assetsLogger.SetLogFuncCall(true)
	}

	//log.Debug("already got scanner:", assetsMgr)
	scanner := assetsMgr.GetBlockScanner()
	//scanner.SetRescanBlockHeight(6518561)

	if scanner == nil {
		log.Error(symbol, "is not support block scan")
		return
	}
	result, err := scanner.ExtractTransactionData(txid, scanTargetFunc)
	if err != nil {
		t.Errorf("ExtractTransactionData unexpected error %v", err)
		return
	}

	for sourceKey, keyData := range result {
		log.Notice("account:", sourceKey)
		for _, data := range keyData {

			for i, input := range data.TxInputs {
				log.Std.Notice("data.TxInputs[%d]: %+v", i, input)
			}

			for i, output := range data.TxOutputs {
				log.Std.Notice("data.TxOutputs[%d]: %+v", i, output)
			}

			log.Std.Notice("data.Transaction: %+v", data.Transaction)
		}
	}

}
