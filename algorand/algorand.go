package algorand

import (
	"github.com/algorand/go-algorand-sdk/client/algod"
	"github.com/astaxie/beego/config"
	"github.com/blocktree/openwallet/log"
	"github.com/blocktree/openwallet/openwallet"
)

//CurveType 曲线类型
func (wm *WalletManager) CurveType() uint32 {
	return wm.Config.CurveType
}

//FullName 币种全名
func (wm *WalletManager) FullName() string {
	return "algorand"
}

//Symbol 币种标识
func (wm *WalletManager) Symbol() string {
	return wm.Config.Symbol
}

//Decimal 小数位精度
func (wm *WalletManager) Decimal() int32 {
	return 6
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

	return nil
}

//LoadAssetsConfig 加载外部配置
func (wm *WalletManager) LoadAssetsConfig(c config.Configer) error {

	wm.Config.ServerAPI = c.String("ServerAPI")
	wm.Config.ServerToken = c.String("ServerToken")
	client, _ := algod.MakeClient(wm.Config.ServerAPI, wm.Config.ServerToken)
	wm.client = &client
	// wm.Config.NetworkID = c.String("networkID")
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
