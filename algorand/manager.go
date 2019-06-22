package algorand

import (
	"github.com/algorand/go-algorand-sdk/client/algod"
	"github.com/blocktree/openwallet/log"
	"github.com/blocktree/openwallet/openwallet"
)

type WalletManager struct {
	openwallet.AssetsAdapterBase

	Config          *WalletConfig                   // 节点配置
	Decoder         openwallet.AddressDecoder       //地址编码器
	TxDecoder       openwallet.TransactionDecoder   //交易单编码器
	Log             *log.OWLogger                   //日志工具
	ContractDecoder openwallet.SmartContractDecoder //智能合约解析器
	Blockscanner    *AlgoBlockScanner               //区块扫描器
	client          *algod.Client                   //algod client
}

func NewWalletManager() *WalletManager {
	wm := WalletManager{}
	wm.Config = NewConfig(Symbol)
	wm.Blockscanner = NewAlgoBlockScanner(&wm)
	wm.Decoder = NewAddressDecoder(&wm)
	wm.TxDecoder = NewTransactionDecoder(&wm)
	wm.Log = log.NewOWLogger(wm.Symbol())
	return &wm
}
