package eth

import (
	"fmt"

	"github.com/chenzhijie/go-web3/abi"
	"github.com/chenzhijie/go-web3/rpc"
	"github.com/chenzhijie/go-web3/types"
	"github.com/ethereum/go-ethereum/common"
)

type Contract struct {
	abi      *abi.ABI
	addr     common.Address
	provider *rpc.Client
}

func (c *Contract) Methods(methodName string) *abi.Method {
	m, _ := c.abi.Methods[methodName]
	return m
}

func (c *Contract) Address() common.Address {
	return c.addr
}

func (c *Contract) Call(methodName string, args ...interface{}) (string, error) {
	m := c.Methods(methodName)
	if m == nil {
		return "", fmt.Errorf("method %v not found", methodName)
	}
	data, err := m.EncodeABI(args...)
	if err != nil {
		return "", err
	}
	// fmt.Printf("data %x\n", data)
	msg := &types.CallMsg{
		To:   c.addr,
		Data: data,
	}
	// fmt.Printf("msg %v\n", msg)

	var out string
	if err := c.provider.Call("eth_call", &out, msg, "latest"); err != nil {
		return "", err
	}
	return out, nil

}

func (e *Eth) NewContract(abiString string, contractAddr ...string) (*Contract, error) {
	a, err := abi.NewABI(abiString)
	if err != nil {
		return nil, err
	}

	var addr common.Address
	if contractAddr != nil && len(contractAddr) == 1 {
		addr = common.HexToAddress(contractAddr[0])
	}
	c := &Contract{
		abi:      a,
		addr:     addr,
		provider: e.c,
	}
	return c, nil
}
