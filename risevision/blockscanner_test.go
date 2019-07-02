package risevision

import (
	"github.com/blocktree/openwallet/log"
	"github.com/blocktree/openwallet/openwallet"
	"testing"
)


func TestRISEBlockScanner_GetGlobalMaxBlockHeight(t *testing.T) {
	wm := testNewWalletManager()
	scanner := wm.GetBlockScanner()
	height :=scanner.GetGlobalMaxBlockHeight()
	log.Infof("height: %v", height)
}




func TestAEBlockScanner_GetBlockByHeight(t *testing.T) {
	wm := testNewWalletManager()
	block, err := wm.Blockscanner.GetBlockByHeight(67213)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	log.Infof("block: %v", block)
}


func TestAEBlockScanner_ExtractTransactionData(t *testing.T) {

	//GetSourceKeyByAddress 获取地址对应的数据源标识
	scanTargetFunc := func(target openwallet.ScanTarget) (string, bool) {
		//if target.Address == "	ak_qcqXt6ySgRPvBkNwEpNMvaKWzrhPZsoBHLvgg68qg9vRht62y" {
		//	return "sender", true
		//} else if target.Address == "ak_mPXUBSsSCJgfu3yz2i2AiVTtLA2TzMyMJL5e6X7shM9Qa246t" {
		//	return "recipient", true
		//}
		return "", true
	}

	wm := testNewWalletManager()
	txs, err := wm.Blockscanner.ExtractTransactionData("10532021534936299833", scanTargetFunc)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	log.Infof("txs: %v", txs)
}
