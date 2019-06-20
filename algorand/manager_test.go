package algorand

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/astaxie/beego/config"

	"github.com/stellar/go/keypair"
)

var (
	tw *WalletManager
)

func testNewWalletManager() *WalletManager {
	wm := NewWalletManager()

	//读取配置
	absFile := filepath.Join("conf", "conf.ini")
	//log.Debug("absFile:", absFile)
	c, err := config.NewConfig("ini", absFile)
	if err != nil {
		return nil
	}
	wm.LoadAssetsConfig(c)
	return wm
}

func init() {
	tw = testNewWalletManager()
}

func TestAccount(t *testing.T) {

	const (
		account1Addr   = "GAVDK2OHFZ5B257PRTCOFYNGRIWV5JRCD5SINMLQJUMSSVYV4LVHI4CN"
		account1Secret = "SDNKCPIVRCS76DATVQUFXDO73DPSXVJ22YCIS46JOBV3UR47ONWFKEUX"
		//account2Secret = "SBOEFVTSQCFFTHHFAIPLOBMDY32JC4E4KEHR4TKCSUE2O5BSBTHOAANH"
	)

	sender, err := keypair.Parse(account1Secret)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	txid := "5d9d4712a05361619a4608a4e2560bbb6f941a8244364bd61c875bdb3945944a"
	txid = strings.Trim(txid, "\"")
	fmt.Printf("txid: %s\n", txid)
	fmt.Printf("pub: %s\n", sender.Address())

}
