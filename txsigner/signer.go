package txsigner

import (
	"fmt"

	"github.com/blocktree/go-owcrypt"
)

var Default = &TransactionSigner{}

type TransactionSigner struct {
}

// SignTransactionHash 交易哈希签名算法
// required
func (singer *TransactionSigner) SignTransactionHash(msg []byte, privateKey []byte, eccType uint32) ([]byte, error) {
	sig, _, ret := owcrypt.Signature(privateKey, nil, msg, eccType)
	if ret != owcrypt.SUCCESS {
		return nil, fmt.Errorf("ECC sign hash failed")
	}
	return sig, nil
}
