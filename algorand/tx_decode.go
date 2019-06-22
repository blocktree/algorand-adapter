package algorand

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/algorand/go-algorand-sdk/client/algod/models"
	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/transaction"
	"github.com/algorand/go-algorand-sdk/types"
	"github.com/blocktree/algorand-adapter/txsigner"
	"github.com/blocktree/go-owcrypt"
	"github.com/blocktree/openwallet/common"
	"github.com/blocktree/openwallet/log"
	"github.com/blocktree/openwallet/openwallet"
	"github.com/shopspring/decimal"
)

// txidPrefix is prepended to a transaction when computing its txid
var txidPrefix = []byte("TX")

type TransactionDecoder struct {
	openwallet.TransactionDecoderBase
	wm *WalletManager //钱包管理者
}

//NewTransactionDecoder 交易单解析器
func NewTransactionDecoder(wm *WalletManager) *TransactionDecoder {
	decoder := TransactionDecoder{}
	decoder.wm = wm
	return &decoder
}

//CreateRawTransaction 创建交易单
func (decoder *TransactionDecoder) CreateRawTransaction(wrapper openwallet.WalletDAI, rawTx *openwallet.RawTransaction) error {

	var (
		decimals  = decoder.wm.Decimal()
		accountID = rawTx.Account.AccountID
		// fixFees         = big.NewInt(0)
		findAddrBalance *AddrBalance
		feeInfo         *txFeeInfo
	)

	//获取wallet
	addresses, err := wrapper.GetAddressList(0, -1, "AccountID", accountID)
	if err != nil {
		return err
	}

	if len(addresses) == 0 {
		return fmt.Errorf("[%s] have not addresses", accountID)
	}

	var amountStr string
	for _, v := range rawTx.To {
		amountStr = v
		break
	}

	//todo cal fee
	amount := common.StringNumToBigIntWithExp(amountStr, decimals)

	for _, addr := range addresses {
		resp, _ := decoder.wm.Blockscanner.GetBalanceByAddress(addr.Address)
		if len(resp) == 0 {
			continue
		}

		b, err := strconv.ParseInt(resp[0].ConfirmBalance, 10, 64)
		if err != nil {
			continue
		}
		balanceAmount := big.NewInt(b)

		//总消耗数量 = 转账数量 + 手续费
		totalAmount := new(big.Int)
		totalAmount.Add(amount, feeInfo.Fee)

		//余额不足查找下一个地址
		if balanceAmount.Cmp(totalAmount) < 0 {
			continue
		}

		//只要找到一个合适使用的地址余额就停止遍历
		findAddrBalance = NewAddrBalance(resp[0])
		break
	}

	if findAddrBalance == nil {
		return fmt.Errorf("all address's balance of account is not enough")
	}

	parmas, err := decoder.wm.client.SuggestedParams()
	if err != nil {
		return fmt.Errorf("gets the suggested transaction parameters fail", err)
	}

	//最后创建交易单
	err = decoder.createRawTransaction(
		wrapper,
		rawTx,
		findAddrBalance,
		parmas)
	if err != nil {
		return err
	}

	return nil

}

//SignRawTransaction 签名交易单
func (decoder *TransactionDecoder) SignRawTransaction(wrapper openwallet.WalletDAI, rawTx *openwallet.RawTransaction) error {

	if rawTx.Signatures == nil || len(rawTx.Signatures) == 0 {
		return fmt.Errorf("transaction signature is empty")
	}

	key, err := wrapper.HDKey()
	if err != nil {
		return err
	}

	keySignatures := rawTx.Signatures[rawTx.Account.AccountID]
	if keySignatures != nil {
		for _, keySignature := range keySignatures {

			childKey, err := key.DerivedKeyWithPath(keySignature.Address.HDPath, keySignature.EccType)
			keyBytes, err := childKey.GetPrivateKeyBytes()
			if err != nil {
				return err
			}

			publicKey, _ := hex.DecodeString(keySignature.Address.PublicKey)

			msg, err := hex.DecodeString(keySignature.Message)
			if err != nil {
				return fmt.Errorf("decoder transaction hash failed, unexpected err: %v", err)
			}

			sig, err := txsigner.Default.SignTransactionHash(msg, keyBytes, keySignature.EccType)
			if err != nil {
				return fmt.Errorf("sign transaction hash failed, unexpected err: %v", err)
			}

			decoder.wm.Log.Debugf("message: %s", hex.EncodeToString(msg))
			decoder.wm.Log.Debugf("publicKey: %s", hex.EncodeToString(publicKey))
			decoder.wm.Log.Debugf("nonce : %s", keySignature.Nonce)
			decoder.wm.Log.Debugf("signature: %s", hex.EncodeToString(sig))

			keySignature.Signature = hex.EncodeToString(sig)
		}
	}

	decoder.wm.Log.Info("transaction hash sign success")

	rawTx.Signatures[rawTx.Account.AccountID] = keySignatures

	return nil
}

//VerifyRawTransaction 验证交易单，验证交易单并返回加入签名后的交易单
func (decoder *TransactionDecoder) VerifyRawTransaction(wrapper openwallet.WalletDAI, rawTx *openwallet.RawTransaction) error {

	var (
		txn types.Transaction
		sig types.Signature
	)

	if rawTx.Signatures == nil || len(rawTx.Signatures) == 0 {
		//this.wm.Log.Std.Error("len of signatures error. ")
		return fmt.Errorf("transaction signature is empty")
	}

	//支持多重签名
	for accountID, keySignatures := range rawTx.Signatures {
		decoder.wm.Log.Debug("accountID Signatures:", accountID)
		for _, keySignature := range keySignatures {

			messsage, _ := hex.DecodeString(keySignature.Message)
			signature, _ := hex.DecodeString(keySignature.Signature)
			publicKey, _ := hex.DecodeString(keySignature.Address.PublicKey)

			//验证签名
			ret := owcrypt.Verify(publicKey, nil, 0, messsage, uint16(len(messsage)), signature, keySignature.EccType)
			if ret != owcrypt.SUCCESS {
				return fmt.Errorf("transaction verify failed")
			}

			json.Unmarshal([]byte(rawTx.RawHex), txn)
			copy(sig[:], signature[:])
			// Construct the SignedTxn
			stx := types.SignedTxn{
				Sig: sig,
				Txn: txn,
			}

			// Encode the SignedTxn
			stxBytes := msgpack.Encode(stx)

			rawTx.IsCompleted = true
			rawTx.RawHex = string(stxBytes)
			break

		}
	}

	return nil
}

//SendRawTransaction 广播交易单
func (decoder *TransactionDecoder) SubmitRawTransaction(wrapper openwallet.WalletDAI, rawTx *openwallet.RawTransaction) (*openwallet.Transaction, error) {

	resp, err := decoder.wm.client.SendRawTransaction([]byte(rawTx.RawHex))
	if err != nil {
		return nil, err
	}

	log.Infof("Transaction [%s] submitted to the network successfully.", resp.TxID)

	rawTx.TxID = resp.TxID
	rawTx.IsSubmit = true

	decimals := decoder.wm.Decimal()

	//记录一个交易单
	tx := &openwallet.Transaction{
		From:       rawTx.TxFrom,
		To:         rawTx.TxTo,
		Amount:     rawTx.TxAmount,
		Coin:       rawTx.Coin,
		TxID:       rawTx.TxID,
		Decimal:    decimals,
		AccountID:  rawTx.Account.AccountID,
		Fees:       rawTx.Fees,
		SubmitTime: time.Now().Unix(),
	}

	tx.WxID = openwallet.GenTransactionWxID(tx)

	return tx, nil
}

//GetRawTransactionFeeRate 获取交易单的费率
func (decoder *TransactionDecoder) GetRawTransactionFeeRate() (feeRate string, unit string, err error) {
	suggestedFeeRate, err := decoder.wm.client.SuggestedFee()
	return strconv.FormatUint(suggestedFeeRate.Fee, 10), decoder.wm.Config.Symbol, err
}

// //CreateSummaryRawTransaction 创建汇总交易
// func (decoder *TransactionDecoder) CreateSummaryRawTransaction(wrapper openwallet.WalletDAI, sumRawTx *openwallet.SummaryRawTransaction) ([]*openwallet.RawTransaction, error) {

// 	var (
// 		decimals        = decoder.wm.Decimal()
// 		rawTxArray      = make([]*openwallet.RawTransaction, 0)
// 		accountID       = sumRawTx.Account.AccountID
// 		minTransfer     = common.StringNumToBigIntWithExp(sumRawTx.MinTransfer, decimals)
// 		retainedBalance = common.StringNumToBigIntWithExp(sumRawTx.RetainedBalance, decimals)
// 		fixFees         = big.NewInt(0)
// 		feeInfo         *txFeeInfo
// 	)

// 	if minTransfer.Cmp(retainedBalance) < 0 {
// 		return nil, fmt.Errorf("mini transfer amount must be greater than address retained balance")
// 	}

// 	//获取wallet
// 	addresses, err := wrapper.GetAddressList(sumRawTx.AddressStartIndex, sumRawTx.AddressLimit,
// 		"AccountID", sumRawTx.Account.AccountID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if len(addresses) == 0 {
// 		return nil, fmt.Errorf("[%s] have not addresses", accountID)
// 	}

// 	if len(sumRawTx.FeeRate) > 0 {
// 		fixFees = common.StringNumToBigIntWithExp(sumRawTx.FeeRate, decimals)
// 	} else {
// 		fixFees = common.StringNumToBigIntWithExp(decoder.wm.Config.FixFees, decimals)
// 	}

// 	//计算手续费
// 	//计算手续费
// 	feeInfo = &txFeeInfo{
// 		Fee:      fixFees,
// 		GasPrice: fixFees,
// 		GasUsed:  big.NewInt(1),
// 	}

// 	for _, addr := range addresses {

// 		account, exist, _ := decoder.wm.GetAccounts(addr.Address)
// 		if !exist {
// 			continue
// 		}

// 		//检查余额是否超过最低转账
// 		addrBalance_BI := account.Balance
// 		addrBalance := common.BigIntToDecimals(account.Balance, decoder.wm.Decimal())

// 		if addrBalance_BI.Cmp(minTransfer) < 0 || addrBalance_BI.Cmp(big.NewInt(0)) <= 0 {
// 			continue
// 		}
// 		//计算汇总数量 = 余额 - 保留余额
// 		sumAmount_BI := new(big.Int)
// 		sumAmount_BI.Sub(addrBalance_BI, retainedBalance)

// 		//减去手续费
// 		sumAmount_BI.Sub(sumAmount_BI, feeInfo.Fee)
// 		if sumAmount_BI.Cmp(big.NewInt(0)) <= 0 {
// 			continue
// 		}

// 		sumAmount := common.BigIntToDecimals(sumAmount_BI, decimals)
// 		feesAmount := common.BigIntToDecimals(feeInfo.Fee, decimals)

// 		decoder.wm.Log.Debugf("balance: %v", addrBalance.String())
// 		decoder.wm.Log.Debugf("fees: %v", feesAmount)
// 		decoder.wm.Log.Debugf("sumAmount: %v", sumAmount)

// 		//创建一笔交易单
// 		rawTx := &openwallet.RawTransaction{
// 			Coin:    sumRawTx.Coin,
// 			Account: sumRawTx.Account,
// 			To: map[string]string{
// 				sumRawTx.SummaryAddress: sumAmount.StringFixed(decoder.wm.Decimal()),
// 			},
// 			Required: 1,
// 		}

// 		createErr := decoder.createRawTransaction(
// 			wrapper,
// 			rawTx,
// 			account,
// 			feeInfo,
// 			"")
// 		if createErr != nil {
// 			return nil, createErr
// 		}

// 		//创建成功，添加到队列
// 		rawTxArray = append(rawTxArray, rawTx)

// 	}

// 	return rawTxArray, nil
// }

//createRawTransaction
func (decoder *TransactionDecoder) createRawTransaction(
	wrapper openwallet.WalletDAI,
	rawTx *openwallet.RawTransaction,
	addrBalance *AddrBalance,
	suggestedParams models.TransactionParams) error {

	var (
		accountTotalSent = decimal.Zero
		txFrom           = make([]string, 0)
		txTo             = make([]string, 0)
		keySignList      = make([]*openwallet.KeySignature, 0)
		amountStr        string
		destination      string
	)

	decimals := int32(0)
	if rawTx.Coin.IsContract {
		decimals = int32(rawTx.Coin.Contract.Decimals)
	} else {
		decimals = decoder.wm.Decimal()
	}

	for k, v := range rawTx.To {
		destination = k
		amountStr = v
		break
	}

	//计算账户的实际转账amount
	accountTotalSentAddresses, findErr := wrapper.GetAddressList(0, -1, "AccountID", rawTx.Account.AccountID, "Address", destination)
	if findErr != nil || len(accountTotalSentAddresses) == 0 {
		amountDec, _ := decimal.NewFromString(amountStr)
		accountTotalSent = accountTotalSent.Add(amountDec)
	}

	addr, err := wrapper.GetAddress(addrBalance.Address)
	if err != nil {
		return err
	}

	amount := common.StringNumToBigIntWithExp(amountStr, decimals)
	//存在直接转账
	txn, err := transaction.MakePaymentTxn(addrBalance.Address, destination, suggestedParams.Fee, amount.Uint64(), suggestedParams.LastRound, suggestedParams.LastRound+1000, []byte(""), "", suggestedParams.GenesisID, suggestedParams.GenesisHash)
	// todo add nonce

	toBeSigned := rawTransactionBytesToSign(txn)
	rawHex, _ := serialize(txn)
	rawTx.RawHex = string(rawHex)

	if rawTx.Signatures == nil {
		rawTx.Signatures = make(map[string][]*openwallet.KeySignature)
	}

	signature := openwallet.KeySignature{
		EccType: decoder.wm.Config.CurveType,
		Address: addr,
		Nonce:   "", // todo
		Message: hex.EncodeToString(toBeSigned),
	}
	keySignList = append(keySignList, &signature)

	feesAmount := common.IntToDecimals(int64(txn.Fee.ToAlgos()), decimals)
	gasPrice := common.IntToDecimals(int64(suggestedParams.Fee), decimals)
	accountTotalSent = accountTotalSent.Add(feesAmount)
	accountTotalSent = decimal.Zero.Sub(accountTotalSent)

	rawTx.Signatures[rawTx.Account.AccountID] = keySignList
	rawTx.FeeRate = gasPrice.String()
	rawTx.Fees = feesAmount.String()
	rawTx.IsBuilt = true
	rawTx.TxAmount = accountTotalSent.StringFixed(decimals)
	rawTx.TxFrom = txFrom
	rawTx.TxTo = txTo

	return nil
}

//CreateSummaryRawTransactionWithError 创建汇总交易，返回能原始交易单数组（包含带错误的原始交易单）
func (decoder *TransactionDecoder) CreateSummaryRawTransactionWithError(wrapper openwallet.WalletDAI, sumRawTx *openwallet.SummaryRawTransaction) ([]*openwallet.RawTransactionWithError, error) {
	raTxWithErr := make([]*openwallet.RawTransactionWithError, 0)
	rawTxs, err := decoder.CreateSummaryRawTransaction(wrapper, sumRawTx)
	if err != nil {
		return nil, err
	}
	for _, tx := range rawTxs {
		raTxWithErr = append(raTxWithErr, &openwallet.RawTransactionWithError{
			RawTx: tx,
			Error: nil,
		})
	}
	return raTxWithErr, nil
}

// //GetNewNonce  确定txdecode nonce值
// func (decoder *TransactionDecoder) GetNewNonce(wrapper openwallet.WalletDAI, addr *AddrBalance) uint64 {

// 	var (
// 		nonce        uint64
// 		nonce_submit uint64
// 	)
// 	//获取db记录的nonce并确认nonce值
// 	//nonce_cache, _ := wrapper.GetAddressExtParam(addr.Address, PESS_SEQUENCEID_KEY)
// 	////判断nonce_db是否为空,为空则说明当前nonce是0
// 	//if nonce_cache == nil {
// 	//	nonce = addr.SequenceID
// 	//} else {
// 	//	nonce = common.NewString(nonce_cache).UInt64()
// 	//}

// 	nonce_chain := addr.SequenceID

// 	//如果本地nonce_db > 链上nonce,采用本地nonce,否则采用链上nonce
// 	if nonce > nonce_chain {
// 		log.Debugf("use cache nonce")
// 		nonce_submit = nonce
// 	} else {
// 		log.Debugf("use chain nonce")
// 		nonce_submit = nonce_chain
// 	}

// 	return nonce_submit
// }

// rawTransactionBytesToSign returns the byte form of the tx that we actually sign
// and compute txID from.
func rawTransactionBytesToSign(tx types.Transaction) []byte {
	// Encode the transaction as msgpack
	encodedTx := msgpack.Encode(tx)

	// Prepend the hashable prefix
	msgParts := [][]byte{txidPrefix, encodedTx}
	return bytes.Join(msgParts, nil)
}

func serialize(obj interface{}) (encoded []byte, err error) {
	encoded, err = json.Marshal(obj)
	return
}
