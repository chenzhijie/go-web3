package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

const MAX_GAS_LIMIT = 30000000

type CallMsgData []byte

func (a CallMsgData) MarshalText() ([]byte, error) {
	return hexutil.Bytes(a[:]).MarshalText()
}

type CallMsgBigInt big.Int

func (a CallMsgBigInt) MarshalText() ([]byte, error) {
	b := big.Int(a)
	return []byte(hexutil.EncodeBig(&b)), nil
}

func NewCallMsgBigInt(v *big.Int) *CallMsgBigInt {
	if v == nil {
		return nil
	}
	i := CallMsgBigInt(*v)
	return &i
}

type CallMsg struct {
	From     common.Address `json:"from,omitempty"`
	To       common.Address `json:"to"`
	Data     CallMsgData    `json:"data"`
	Gas      *CallMsgBigInt `json:"gas,omitempty"`
	GasPrice *CallMsgBigInt `json:"gasPrice,omitempty"`
	Value    *CallMsgBigInt `json:"value,omitempty"`
}

type ZeroValueCallMsg struct {
	From common.Address `json:"from,omitempty"`
	To   common.Address `json:"to"`
	Data CallMsgData    `json:"data"`
}
