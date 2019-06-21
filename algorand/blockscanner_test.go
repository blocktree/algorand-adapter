package algorand

import (
	"testing"

	"github.com/blocktree/openwallet/log"
	"github.com/blocktree/openwallet/openwallet"
)

func TestALGOBlockScanner_GetBalanceByAddress(t *testing.T) {
	wm := testNewWalletManager()
	b, err := wm.Blockscanner.GetBalanceByAddress("5TSQNIL54GB545B3WLC6OVH653SHAELMHU6MSVNGTUNMOEHAMWG7EC3AA4")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	log.Infof("block = %+v", b)
}

func TestGetBlockHeight(t *testing.T) {
	wm := testNewWalletManager()
	height, err := wm.Blockscanner.GetBlockHeight()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	log.Infof("height: %+v", height)
}

func TestALGOBlockScanner_GetCurrentBlock(t *testing.T) {
	wm := testNewWalletManager()
	b, err := wm.Blockscanner.GetCurrentBlock()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	log.Infof("block = %+v", b)
}

func TestALGOBlockScanner_GetCurrentBlockHeader(t *testing.T) {
	wm := testNewWalletManager()
	header, err := wm.Blockscanner.GetCurrentBlockHeader()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	log.Infof("header: %v", header)
}

func TestALGOBlockScanner_GetBlockByHeight(t *testing.T) {
	wm := testNewWalletManager()
	block, err := wm.Blockscanner.GetBlockByHeight(173745)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	log.Infof("block: %v", block)
}

func TestALGOBlockScanner_GetTransaction(t *testing.T) {
	wm := testNewWalletManager()
	tx, err := wm.Blockscanner.GetTransaction("YRSG7IKDPCK4XMKFFTFFFYMIHF6SJOMHUOIE4FFUWNLEQ4WG2ZOQ")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	log.Infof("tx: %+v", tx)
}

func TestALGOBlockScanner_ExtractTransactionData(t *testing.T) {

	//GetSourceKeyByAddress 获取地址对应的数据源标识
	scanTargetFunc := func(target openwallet.ScanTarget) (string, bool) {
		if target.Address == "GGUMZYT7GHGTOUOMBXVY3AY754UKOKABBD4732COI7IVXMOCR4P4X23YYA" {
			return "sender", true
		} else if target.Address == "JQHFIDCG6KAEZ4DCOPPISLCHXUAXUUPZO4MQOHYWDGC5ALOKTOTIHGHDZM" {
			return "recipient", true
		}
		return "", false
	}

	txs, err := tw.Blockscanner.ExtractTransactionData("EW6IZCOPYGSTTXPQUMD25RLKCMOU66MHU4IFYXH6ZTL3RPQ5G4RA", scanTargetFunc)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	for key, datas := range txs {
		log.Notice("key:", key)
		for _, data := range datas {
			for i, input := range data.TxInputs {
				log.Infof("data.TxInputs[%d]: %+v", i, input)
			}

			for i, output := range data.TxOutputs {
				log.Infof("data.TxOutputs[%d]: %+v", i, output)
			}

			log.Infof("data.Transaction: %+v", data.Transaction)
		}
	}
}
