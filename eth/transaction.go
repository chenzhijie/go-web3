package eth

import (
	"fmt"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	eTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

func (e *Eth) NewEIP1559Tx(
	to common.Address,
	amount *big.Int,
	gasLimit uint64,
	gasTipCap *big.Int,
	gasFeeCap *big.Int,
	data []byte,
	nonce uint64,
) (*eTypes.Transaction, error) {

	dynamicFeeTx := &eTypes.DynamicFeeTx{

		Nonce:     nonce,
		GasTipCap: gasTipCap,
		GasFeeCap: gasFeeCap,
		Gas:       gasLimit,
		To:        &to,
		Value:     amount,
		Data:      data,
	}
	if e.chainId != nil {
		dynamicFeeTx.ChainID = e.chainId
	}

	if e.privateKey == nil {
		return eTypes.NewTx(dynamicFeeTx), nil
	}

	signedTx, err := eTypes.SignNewTx(
		e.privateKey,
		eTypes.LatestSignerForChainID(e.chainId),
		dynamicFeeTx,
	)
	if err != nil {
		return nil, err
	}
	return signedTx, nil
}

func (e *Eth) SendRawEIP1559Transaction(
	to common.Address,
	amount *big.Int,
	nonce uint64,
	gasLimit uint64,
	gasTipCap *big.Int,
	gasFeeCap *big.Int,
	data []byte,
) (common.Hash, error) {
	var hash common.Hash
	dynamicFeeTx := &eTypes.DynamicFeeTx{
		Nonce:     nonce,
		GasTipCap: gasTipCap,
		GasFeeCap: gasFeeCap,
		Gas:       gasLimit,
		To:        &to,
		Value:     amount,
		Data:      data,
	}

	signedTx, err := eTypes.SignNewTx(e.privateKey, eTypes.LatestSignerForChainID(e.chainId), dynamicFeeTx)
	if err != nil {
		return hash, err
	}

	txData, err := signedTx.MarshalBinary()
	if err != nil {
		return hash, err
	}

	err = e.c.Call("eth_sendRawTransaction", &hash, hexutil.Encode(txData))

	return hash, err
}

func (e *Eth) SendRawTransaction(
	to common.Address,
	amount *big.Int,
	nonce uint64,
	gasLimit uint64,
	gasPrice *big.Int,
	data []byte,
) (common.Hash, error) {
	var hash common.Hash

	tx := eTypes.NewTransaction(nonce, to, amount, gasLimit, gasPrice, data)

	signedTx, err := eTypes.SignTx(tx, eTypes.NewEIP155Signer(e.chainId), e.privateKey)
	if err != nil {
		return hash, err
	}
	serializedTx, err := rlp.EncodeToBytes(signedTx)
	if err != nil {
		return hash, err
	}

	err = e.c.Call("eth_sendRawTransaction", &hash, fmt.Sprintf("0x%x", serializedTx))
	return hash, err

}

func (e *Eth) SyncSendRawTransaction(
	to common.Address,
	amount *big.Int,
	nonce uint64,
	gasLimit uint64,
	gasPrice *big.Int,
	data []byte,
) (*eTypes.Receipt, error) {

	tx := eTypes.NewTransaction(nonce, to, amount, gasLimit, gasPrice, data)

	signedTx, err := eTypes.SignTx(tx, eTypes.NewEIP155Signer(e.chainId), e.privateKey)
	if err != nil {
		return nil, err
	}
	serializedTx, err := rlp.EncodeToBytes(signedTx)
	if err != nil {
		return nil, err
	}
	var hash common.Hash
	err = e.c.Call("eth_sendRawTransaction", &hash, fmt.Sprintf("0x%x", serializedTx))
	if err != nil {
		return nil, err
	}

	type ReceiptCh struct {
		ret *eTypes.Receipt
		err error
	}

	var timeoutFlag int32
	ch := make(chan *ReceiptCh, 1)

	go func() {
		for {
			time.Sleep(time.Second)
			receipt, _ := e.GetTransactionReceipt(hash)
			if receipt != nil {
				ch <- &ReceiptCh{
					ret: receipt,
					err: nil,
				}
				break
			}
			if atomic.LoadInt32(&timeoutFlag) == 1 {
				break
			}
		}
	}()

	select {
	case result := <-ch:
		return result.ret, result.err
	case <-time.After(time.Duration(e.txPollTimeout) * time.Second):
		atomic.StoreInt32(&timeoutFlag, 1)
		return nil, fmt.Errorf("Transaction was not mined within %v seconds, "+
			"please make sure your transaction was properly sent. Be aware that it might still be mined!", e.txPollTimeout)
	}
}

func (e *Eth) SyncSendEIP1559RawTransaction(
	to common.Address,
	amount *big.Int,
	nonce uint64,
	gasLimit uint64,
	gasTipCap *big.Int,
	gasFeeCap *big.Int,
	data []byte,
) (*eTypes.Receipt, error) {

	dynamicFeeTx := &eTypes.DynamicFeeTx{
		Nonce:     nonce,
		GasTipCap: gasTipCap,
		GasFeeCap: gasFeeCap,
		Gas:       gasLimit,
		To:        &to,
		Value:     amount,
		Data:      data,
	}

	signedTx, err := eTypes.SignNewTx(e.privateKey, eTypes.LatestSignerForChainID(e.chainId), dynamicFeeTx)
	if err != nil {
		return nil, err
	}

	txData, err := signedTx.MarshalBinary()
	if err != nil {
		return nil, err
	}
	var hash common.Hash
	err = e.c.Call("eth_sendRawTransaction", &hash, hexutil.Encode(txData))
	if err != nil {
		return nil, err
	}

	type ReceiptCh struct {
		ret *eTypes.Receipt
		err error
	}

	var timeoutFlag int32
	ch := make(chan *ReceiptCh, 1)

	go func() {
		for {
			time.Sleep(time.Second)
			receipt, _ := e.GetTransactionReceipt(hash)
			if receipt != nil {
				ch <- &ReceiptCh{
					ret: receipt,
					err: nil,
				}
				break
			}
			if atomic.LoadInt32(&timeoutFlag) == 1 {
				break
			}
		}
	}()

	select {
	case result := <-ch:
		return result.ret, result.err
	case <-time.After(time.Duration(e.txPollTimeout) * time.Second):
		atomic.StoreInt32(&timeoutFlag, 1)
		return nil, fmt.Errorf("Transaction was not mined within %v seconds, "+
			"please make sure your transaction was properly sent. Be aware that it might still be mined!", e.txPollTimeout)
	}
}
