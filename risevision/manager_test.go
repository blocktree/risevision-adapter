package risevision

import (
	"github.com/astaxie/beego/config"
	"github.com/blocktree/openwallet/log"
	"path/filepath"
	"testing"
)

func testNewWalletManager() *WalletManager {
	wm := NewWalletManager()

	//读取配置
	absFile := filepath.Join("conf", "RISE.ini")
	//log.Debug("absFile:", absFile)
	c, err := config.NewConfig("ini", absFile)
	if err != nil {
		return nil
	}
	wm.LoadAssetsConfig(c)
	return wm
}

func init() {

}

func TestWalletManager_GetAccount(t *testing.T) {
	wm := testNewWalletManager()
	r, err := wm.GetAccount("11000023207492553399L")
	if err != nil {
		log.Errorf("unexpected error: %v", err)
		return
	}
	log.Infof("%+v", r)
}

func TestWalletManager_GetAccountPendingTxCount(t *testing.T) {
	wm := testNewWalletManager()
	r, err := wm.GetAccountPendingTxCount("ak_rozWtRmHh91aEu1Qo46wGSHtJfaGtbgRgPEezZvTmHtRu1fqe")
	if err != nil {
		log.Errorf("unexpected error: %v", err)
		return
	}
	log.Infof("%+v", r)
}
