package risevision

import (
	"github.com/astaxie/beego/config"
	"github.com/blocktree/openwallet/log"
	"github.com/blocktree/openwallet/openwallet"
	"github.com/blocktree/risevision-adapter/api"
	"strconv"
	"strings"
)

//CurveType 曲线类型
func (wm *WalletManager) CurveType() uint32 {
	return wm.Config.CurveType
}

//FullName 币种全名
func (wm *WalletManager) FullName() string {
	return "RISE"
}

//Symbol 币种标识
func (wm *WalletManager) Symbol() string {
	return wm.Config.Symbol
}

//Decimal 小数位精度
func (wm *WalletManager) Decimal() int32 {
	return 8
}

//BalanceModelType 余额模型类别
func (wm *WalletManager) BalanceModelType() openwallet.BalanceModelType {
	return openwallet.BalanceModelTypeAddress
}

//GetAddressDecode 地址解析器
func (wm *WalletManager) GetAddressDecode() openwallet.AddressDecoder {
	return wm.Decoder
}

//GetTransactionDecoder 交易单解析器
func (wm *WalletManager) GetTransactionDecoder() openwallet.TransactionDecoder {
	return wm.TxDecoder
}

//GetBlockScanner 获取区块链
func (wm *WalletManager) GetBlockScanner() openwallet.BlockScanner {

	return wm.Blockscanner
}

//LoadAssetsConfig 加载外部配置
func (wm *WalletManager) LoadAssetsConfig(c config.Configer) error {

	wm.Config.ServerAPI = c.String("serverAPI")
	wm.Config.FixFees = "0.1"
	wm.Config.NetworkID = c.String("networkID")
	//wm.client = NewClient(wm.Config.ServerAPI, false)
	serverAPIStrs := strings.Split(wm.Config.ServerAPI, ":")
	if len(serverAPIStrs) > 1 {

		p, err := strconv.Atoi(serverAPIStrs[1])
		if err != nil {
			log.Error("string转int失败")
		} else {

			host := api.Host{
				Hostname: serverAPIStrs[0],
				Port:     p,
			}
			config2 := &api.Config{
				Host:       host,
				RandomHost: false,
			}
			wm.Api = api.NewClientWithCustomConfig(config2)
		}

	} else {
		host := api.Host{
			Hostname: serverAPIStrs[0],
			Secure:   true,
		}
		config2 := &api.Config{
			Host:       host,
			RandomHost: false,
		}
		wm.Api = api.NewClientWithCustomConfig(config2)
	}
	wm.Config.DataDir = c.String("dataDir")

	//数据文件夹
	wm.Config.makeDataDir()

	return nil
}

//InitAssetsConfig 初始化默认配置
func (wm *WalletManager) InitAssetsConfig() (config.Configer, error) {
	return config.NewConfigData("ini", []byte(wm.Config.DefaultConfig))
}

//GetAssetsLogger 获取资产账户日志工具
func (wm *WalletManager) GetAssetsLogger() *log.OWLogger {
	return wm.Log
}

//GetSmartContractDecoder 获取智能合约解析器
func (wm *WalletManager) GetSmartContractDecoder() openwallet.SmartContractDecoder {
	return wm.ContractDecoder
}
