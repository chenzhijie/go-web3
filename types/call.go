package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type CallMsgData []byte

func (a CallMsgData) MarshalText() ([]byte, error) {
	return hexutil.Bytes(a[:]).MarshalText()
}

type CallMsgBigInt big.Int

func (a CallMsgBigInt) MarshalText() ([]byte, error) {
	b := big.Int(a)
	return []byte(hexutil.EncodeBig(&b)), nil
}

func NewCallMsgBigInt(v *big.Int) CallMsgBigInt {
	if v == nil {
		return CallMsgBigInt(*big.NewInt(0))
	}
	i := *v
	return CallMsgBigInt(i)
}

type CallMsg struct {
	From     common.Address `json:"from,omitempty"`
	To       common.Address `json:"to"`
	Data     CallMsgData    `json:"data"`
	GasPrice CallMsgBigInt  `json:"gasPrice,omitempty"`
	Value    CallMsgBigInt  `json:"value,omitempty"`
}
