package eth

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"

	"github.com/chenzhijie/go-web3/rpc"
	"github.com/chenzhijie/go-web3/types"
	"github.com/chenzhijie/go-web3/utils"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type Contract struct {
	abi      abi.ABI
	addr     common.Address
	provider *rpc.Client
}

func (c *Contract) AllMethods() []string {
	methodNames := make([]string, 0)
	for methodName := range c.abi.Methods {
		methodNames = append(methodNames, methodName)
	}
	return methodNames
}

func (c *Contract) Methods(methodName string) abi.Method {
	m, _ := c.abi.Methods[methodName]
	return m
}

func (c *Contract) Address() common.Address {
	return c.addr
}

func (c *Contract) Call(methodName string, args ...interface{}) (interface{}, error) {

	data, err := c.EncodeABI(methodName, args...)

	if err != nil {
		return nil, err
	}

	msg := &types.ZeroValueCallMsg{
		To:   c.addr,
		Data: data,
	}

	var out string
	if err := c.provider.Call("eth_call", &out, msg, "latest"); err != nil {
		return nil, err
	}

	outputBytes, err := hexutil.Decode(out)
	if err != nil {
		return nil, err
	}

	response, err := c.abi.Unpack(methodName, outputBytes)
	if err != nil {
		return nil, err
	}
	if len(response) != 1 {
		return response, nil
	}
	return response[0], nil
}

func (c *Contract) CallWithMultiReturns(methodName string, args ...interface{}) ([]interface{}, error) {
	return c.CallAtWithMultiReturns(nil, methodName, args...)
}

func (c *Contract) CallAtWithMultiReturns(blockNumber *big.Int, methodName string, args ...interface{}) ([]interface{}, error) {

	data, err := c.EncodeABI(methodName, args...)

	if err != nil {
		return nil, err
	}

	msg := &types.CallMsg{
		To:   c.addr,
		Data: data,
		Gas:  types.NewCallMsgBigInt(big.NewInt(types.MAX_GAS_LIMIT)),
	}

	var out string
	if err := c.provider.Call("eth_call", &out, msg, utils.ToBlockNumArg(blockNumber)); err != nil {
		return nil, err
	}

	outputBytes, err := hexutil.Decode(out)
	if err != nil {
		return nil, err
	}

	response, err := c.abi.Unpack(methodName, outputBytes)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Contract) CallWithFromAndValue(
	methodName string,
	from common.Address,
	value *big.Int,
	args ...interface{},
) ([]interface{}, error) {

	data, err := c.EncodeABI(methodName, args...)

	if err != nil {
		return nil, err
	}

	msg := &types.CallMsg{
		From: from,
		To:   c.addr,
		Data: data,
		Gas:  types.NewCallMsgBigInt(big.NewInt(types.MAX_GAS_LIMIT)),
	}
	if value != nil {
		msg.Value = types.NewCallMsgBigInt(value)
	}

	var out string
	if err := c.provider.Call("eth_call", &out, msg, "latest"); err != nil {
		return nil, err
	}

	outputBytes, err := hexutil.Decode(out)
	if err != nil {
		return nil, err
	}

	response, err := c.abi.Unpack(methodName, outputBytes)
	if err != nil {
		return nil, err
	}
	if len(response) == 0 {
		return nil, fmt.Errorf("invalid response %v", response)
	}
	return response, nil
}

func (c *Contract) EncodeABI(methodName string, args ...interface{}) ([]byte, error) {
	m := c.Methods(methodName)
	if len(m.ID) == 0 {
		return nil, fmt.Errorf("method %v not found", methodName)
	}
	data := m.ID

	inputData, err := m.Inputs.Pack(args...)
	if err != nil {
		return nil, err
	}
	return append(data, inputData...), nil
}

func NewContract(abiString string, contractAddr ...string) (*Contract, error) {
	if len(abiString) == 0 {
		return nil, errors.New("invalid abi json string")
	}
	a, err := abi.JSON(bytes.NewReader([]byte(abiString)))
	if err != nil {
		return nil, err
	}

	var addr common.Address
	if len(contractAddr) == 1 {
		addr = common.HexToAddress(contractAddr[0])
	}
	c := &Contract{
		abi:  a,
		addr: addr,
	}
	return c, nil
}

func (e *Eth) NewContract(abiString string, contractAddr ...string) (*Contract, error) {
	c, err := NewContract(abiString, contractAddr...)
	if err != nil {
		return nil, err
	}
	c.provider = e.c

	return c, nil
}
