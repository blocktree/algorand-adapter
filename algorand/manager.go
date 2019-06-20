package algorand

import (
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
	// Blockscanner    *ALGOBlockScanner               //区块扫描器
	// client          *Client                         //本地封装的http client
}

func NewWalletManager() *WalletManager {
	wm := WalletManager{}
	wm.Config = NewConfig(Symbol)
	// wm.Blockscanner = NewALGOBlockScanner(&wm)
	wm.Decoder = NewAddressDecoder(&wm)
	// wm.TxDecoder = NewTransactionDecoder(&wm)
	wm.Log = log.NewOWLogger(wm.Symbol())
	// wm.ContractDecoder = NewContractDecoder(&wm)
	return &wm
}

// // GetAccounts 获取账户信息
// // @return 余额，是否存在，错误
// func (wm *WalletManager) GetAccounts(address string) (*AddrBalance, bool, error) {
// 	path := fmt.Sprintf("accounts/%s", address)
// 	r, err := wm.client.Call(path, "GET", nil)
// 	if err != nil {
// 		if err.Code == ErrNotFound {
// 			return nil, false, nil
// 		} else {
// 			return nil, false, err
// 		}
// 	}
// 	account := NewAddrBalance(r)
// 	return account, true, nil
// }

// // BroadcastTransaction 广播交易单
// func (wm *WalletManager) BroadcastTransaction(tx *transaction.Transaction) (string, error) {
// 	_, err := wm.client.Call("transactions", "POST", tx)
// 	if err != nil {
// 		return "", err
// 	}
// 	txid := tx.GetHash()
// 	return txid, nil
// }
