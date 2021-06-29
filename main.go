package main

import (
	"github.com/blocktree/algorand-adapter/algorand"
	"github.com/blocktree/openwallet/log"
)

var (
	WalletManager algorand.WalletManager
)

func init() {
	wm := algorand.WalletManager{}
	wm.Config = algorand.NewConfig(algorand.Symbol)
	wm.Blockscanner = algorand.NewAlgoBlockScanner(&wm)
	wm.Decoder = algorand.NewAddressDecoder(&wm)
	wm.TxDecoder = algorand.NewTransactionDecoder(&wm)
	wm.Log = log.NewOWLogger(wm.Symbol())
	WalletManager = wm
}

func main() {

}
