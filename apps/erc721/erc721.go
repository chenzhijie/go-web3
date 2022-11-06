package erc721

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

type ERC721 struct {
	contr         *eth.Contract
	w3            *web3.Web3
	confirmation  int
	txPollTimeout int
}

func NewERC721(w3 *web3.Web3, contractAddress common.Address) (*ERC721, error) {
	contr, err := w3.Eth.NewContract(ERC721_ABI, contractAddress.String())
	if err != nil {
		return nil, err
	}
	e := &ERC721{
		contr:         contr,
		w3:            w3,
		txPollTimeout: 720,
	}
	return e, nil
}

func (e *ERC721) Address() common.Address {
	return e.contr.Address()
}

func (e *ERC721) SetConfirmation(blockCount int) {
	e.confirmation = blockCount
}

func (e *ERC721) SetTxPollTimeout(txPollTimeout int) {
	e.txPollTimeout = txPollTimeout
}

func (e *ERC721) TotalSupply() (*big.Int, error) {
	ret, err := e.contr.Call("totalSupply")
	if err != nil {
		return nil, err
	}
	supply, ok := ret.(*big.Int)
	if !ok {
		return nil, fmt.Errorf("invalid response %v type %T expect *big.Int", ret, ret)
	}
	return supply, nil
}

func (e *ERC721) BalanceOf(owner common.Address) (*big.Int, error) {

	ret, err := e.contr.Call("balanceOf", owner)
	if err != nil {
		return nil, err
	}

	bal, ok := ret.(*big.Int)
	if !ok {
		return nil, fmt.Errorf("invalid result %v, type %T", ret, ret)
	}
	return bal, nil
}

func (e *ERC721) OwnerOf(tokenId *big.Int) (common.Address, error) {
	ret, err := e.contr.Call("ownerOf", tokenId)
	if err != nil {
		return common.Address{}, err
	}

	owner, ok := ret.(common.Address)
	if !ok {
		return common.Address{}, fmt.Errorf("invalid result %v, type %T", ret, ret)
	}
	return owner, nil
}

func (e *ERC721) IsApprovedForAll(owner, operator common.Address) (bool, error) {
	ret, err := e.contr.Call("isApprovedForAll", owner, operator)
	if err != nil {
		return false, err
	}
	approved, ok := ret.(bool)
	if !ok {
		return false, fmt.Errorf("invalid response type %T", ret)
	}
	return approved, nil
}

func (e *ERC721) SetApprovalForAll(
	spender common.Address,
	approve bool,
	gasPrice *big.Int,
	gasTipCap *big.Int,
	gasFeeCap *big.Int,
) (common.Hash, error) {
	code, err := e.contr.EncodeABI("setApprovalForAll", spender, approve)
	if err != nil {
		return common.Hash{}, err
	}

	return e.invokeAndWait(code, gasPrice, gasTipCap, gasFeeCap)
}

func (e *ERC721) TransferFrom(
	from common.Address,
	to common.Address,
	tokenId *big.Int,
	gasPrice *big.Int,
	gasTipCap *big.Int,
	gasFeeCap *big.Int,
) (common.Hash, error) {
	code, err := e.contr.EncodeABI("transferFrom", from, to, tokenId)
	if err != nil {
		return common.Hash{}, err
	}
	return e.invokeAndWait(code, gasPrice, gasTipCap, gasFeeCap)
}

func (e *ERC721) IsApprovalForAll(owner, spender common.Address) (bool, error) {
	ret, err := e.contr.Call("isApprovedForAll", owner, spender)
	if err != nil {
		return false, err
	}
	approved, ok := ret.(bool)
	if !ok {
		return false, fmt.Errorf("invalid response %v type %T expect bool", ret, ret)
	}
	return approved, nil
}

func (e *ERC721) EstimateGasLimit(to common.Address, data []byte, gasPrice, wei *big.Int) (uint64, error) {
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
	if bytes.Compare(emptyAddr[:], from[:]) != 0 {
		call.From = from
	}

	gasLimit, err := e.w3.Eth.EstimateGas(call)
	if err != nil {
		return 0, err
	}
	return gasLimit, nil
}

func (e *ERC721) WaitBlock(blockCount uint64) error {
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

func (e *ERC721) SyncSendEIP1559Tx(
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

func (e *ERC721) SyncSendRawTransactionForTx(
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

func (e *ERC721) invokeAndWait(code []byte, gasPrice, gasTipCap, gasFeeCap *big.Int) (common.Hash, error) {
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

	if e.confirmation == 0 {
		return tx.TxHash, nil
	}

	if err := e.WaitBlock(uint64(e.confirmation)); err != nil {
		return common.Hash{}, err
	}

	return tx.TxHash, nil
}
