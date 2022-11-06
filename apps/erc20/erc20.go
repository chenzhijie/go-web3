package erc20

import (
	"bytes"
	"fmt"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/chenzhijie/go-web3"
	"github.com/chenzhijie/go-web3/eth"
	"github.com/chenzhijie/go-web3/types"
	"github.com/ethereum/go-ethereum/common"
	eTypes "github.com/ethereum/go-ethereum/core/types"
)

type ERC20 struct {
	contr         *eth.Contract
	w3            *web3.Web3
	confirmation  int
	txPollTimeout int
}

func NewERC20(w3 *web3.Web3, contractAddress common.Address) (*ERC20, error) {
	contr, err := w3.Eth.NewContract(ERC20_ABI, contractAddress.String())
	if err != nil {
		return nil, err
	}
	e := &ERC20{
		contr:         contr,
		w3:            w3,
		txPollTimeout: 720,
	}
	return e, nil
}

func (e *ERC20) Address() common.Address {
	return e.contr.Address()
}

func (e *ERC20) SetConfirmation(blockCount int) {
	e.confirmation = blockCount
}

func (e *ERC20) SetTxPollTimeout(txPollTimeout int) {
	e.txPollTimeout = txPollTimeout
}

func (e *ERC20) Allowance(owner, spender common.Address) (*big.Int, error) {

	ret, err := e.contr.Call("allowance", owner, spender)
	if err != nil {
		return nil, err
	}

	allow, ok := ret.(*big.Int)
	if !ok {
		return nil, fmt.Errorf("invalid result %v, type %T", ret, ret)
	}
	return allow, nil
}

func (e *ERC20) Decimals() (uint8, error) {
	ret, err := e.contr.Call("decimals")
	if err != nil {
		return 0, err
	}

	decimals, ok := ret.(uint8)
	if !ok {
		return 0, fmt.Errorf("invalid result %v, type %T", ret, ret)
	}
	return decimals, nil
}

func (e *ERC20) Symbol() (string, error) {
	ret, err := e.contr.Call("symbol")
	if err != nil {
		return "", err
	}

	symbol, ok := ret.(string)
	if !ok {
		return "", fmt.Errorf("invalid result %v, type %T", ret, ret)
	}
	return symbol, nil
}

func (e *ERC20) BalanceOf(owner common.Address) (*big.Int, error) {

	ret, err := e.contr.Call("balanceOf", owner)
	if err != nil {
		return nil, err
	}

	allow, ok := ret.(*big.Int)
	if !ok {
		return nil, fmt.Errorf("invalid result %v, type %T", ret, ret)
	}
	return allow, nil
}

func (e *ERC20) Approve(spender common.Address, limit, gasPrice, gasTipCap, gasFeeCap *big.Int) (common.Hash, error) {

	code, err := e.contr.EncodeABI("approve", spender, limit)
	if err != nil {
		return common.Hash{}, err
	}

	return e.invokeAndWait(code, gasPrice, gasTipCap, gasFeeCap)
}

func (e *ERC20) Transfer(to common.Address, amount, gasPrice, gasTipCap, gasFeeCap *big.Int) (common.Hash, error) {
	code, err := e.contr.EncodeABI("transfer", to, amount)
	if err != nil {
		return common.Hash{}, err
	}

	return e.invokeAndWait(code, gasPrice, gasTipCap, gasFeeCap)
}

func (e *ERC20) EstimateGasLimit(to common.Address, data []byte, gasPrice, wei *big.Int) (uint64, error) {
	call := &types.CallMsg{
		To:    to,
		Data:  data,
		Gas:   types.NewCallMsgBigInt(big.NewInt(types.MAX_GAS_LIMIT)),
		Value: types.NewCallMsgBigInt(wei),
	}
	if gasPrice != nil {
		call.GasPrice = types.NewCallMsgBigInt(gasPrice)
	}

	var emptyAddr common.Address
	from := e.w3.Eth.Address()
	if !bytes.Equal(emptyAddr[:], from[:]) {
		call.From = from
	}

	gasLimit, err := e.w3.Eth.EstimateGas(call)
	if err != nil {
		return 0, err
	}
	return gasLimit, nil
}

func (e *ERC20) WaitBlock(blockCount uint64) error {
	num, err := e.w3.Eth.GetBlockNumber()
	if err != nil {
		return err
	}
	ti := time.NewTicker(time.Second)
	defer ti.Stop()
	for {
		<-ti.C
		nextNum, err := e.w3.Eth.GetBlockNumber()
		if err != nil {
			return err
		}
		if nextNum >= num+blockCount {
			return nil
		}
	}
}

func (e *ERC20) SyncSendRawTransactionForTx(
	gasPrice *big.Int, gasLimit uint64, to common.Address, data []byte, wei *big.Int,
) (*eTypes.Receipt, error) {
	nonce, err := e.w3.Eth.GetNonce(e.w3.Eth.Address(), nil)
	if err != nil {
		return nil, err
	}
	hash, err := e.w3.Eth.SendRawTransaction(to, wei, nonce, gasLimit, gasPrice, data)
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
			receipt, err := e.w3.Eth.GetTransactionReceipt(hash)
			if err != nil && err.Error() != "not found" {
				ch <- &ReceiptCh{
					err: err,
				}
				break
			}
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
		// fmt.Println("send tx done")
	}()

	select {
	case result := <-ch:
		if result.err != nil {
			return nil, err
		}

		return result.ret, nil
	case <-time.After(time.Duration(e.txPollTimeout) * time.Second):
		atomic.StoreInt32(&timeoutFlag, 1)
		return nil, fmt.Errorf("transaction was not mined within %v seconds, "+
			"please make sure your transaction was properly sent. Be aware that it might still be mined!", e.txPollTimeout)
	}
}

func (e *ERC20) SyncSendEIP1559Tx(
	gasTipCap *big.Int,
	gasFeeCap *big.Int,
	gasLimit uint64,
	to common.Address,
	data []byte,
	wei *big.Int,
) (*eTypes.Receipt, error) {
	nonce, err := e.w3.Eth.GetNonce(e.w3.Eth.Address(), nil)
	if err != nil {
		return nil, err
	}
	hash, err := e.w3.Eth.SendRawEIP1559Transaction(to, wei, nonce, gasLimit, gasTipCap, gasFeeCap, data)
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
			receipt, err := e.w3.Eth.GetTransactionReceipt(hash)
			if err != nil && err.Error() != "not found" {
				ch <- &ReceiptCh{
					err: err,
				}
				break
			}
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
		// fmt.Println("send tx done")
	}()

	select {
	case result := <-ch:
		if result.err != nil {
			return nil, err
		}

		return result.ret, nil
	case <-time.After(time.Duration(e.txPollTimeout) * time.Second):
		atomic.StoreInt32(&timeoutFlag, 1)
		return nil, fmt.Errorf("transaction was not mined within %v seconds, "+
			"please make sure your transaction was properly sent. Be aware that it might still be mined!", e.txPollTimeout)
	}
}

func (e *ERC20) invokeAndWait(code []byte, gasPrice, gasTipCap, gasFeeCap *big.Int) (common.Hash, error) {
	gasLimit, err := e.EstimateGasLimit(e.contr.Address(), code, nil, nil)
	if err != nil {
		return common.Hash{}, err
	}

	var tx *eTypes.Receipt
	if gasPrice != nil {
		tx, err = e.SyncSendRawTransactionForTx(gasPrice, gasLimit, e.contr.Address(), code, nil)
	} else {
		tx, err = e.SyncSendEIP1559Tx(gasTipCap, gasFeeCap, gasLimit, e.contr.Address(), code, nil)
	}

	if err != nil {
		return common.Hash{}, err
	}

	if tx == nil {
		return common.Hash{}, nil
	}

	if e.confirmation == 0 {
		return tx.TxHash, nil
	}

	if err := e.WaitBlock(uint64(e.confirmation)); err != nil {
		return common.Hash{}, err
	}

	return tx.TxHash, nil
}
