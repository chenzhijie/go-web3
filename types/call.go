package types

import (
	"fmt"
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
	return []byte(fmt.Sprintf("0x%x", &b)), nil
}

func NewCallMsgBigInt(v *big.Int) CallMsgBigInt {
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
