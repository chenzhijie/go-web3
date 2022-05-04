package utils

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func (u *Utils) EncodeFunctionSignature(funcName string) []byte {
	return crypto.Keccak256([]byte(funcName))[:4]
}

func (u *Utils) DecodeParameters(parameters []string, data []byte) ([]interface{}, error) {

	args := make(abi.Arguments, 0)

	for _, p := range parameters {
		arg := abi.Argument{}
		var err error
		arg.Type, err = abi.NewType(p, "", nil)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}

	return args.Unpack(data)
}

func (u *Utils) EncodeParameters(parameters []string, data []interface{}) ([]byte, error) {

	args := make(abi.Arguments, 0)

	for _, p := range parameters {
		arg := abi.Argument{}
		var err error
		arg.Type, err = abi.NewType(p, "", nil)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}

	return args.Pack(data...)
}

func (u *Utils) PackCode(signature string, args []string, params []interface{}) []byte {
	methodSig := u.EncodeFunctionSignature(signature)
	if len(args) == 0 {
		return methodSig
	}
	inputCode, err := u.EncodeParameters(args, params)
	if err != nil {
		panic(err)
	}
	code := append(methodSig, inputCode...)
	return code
}

// Equal to solidity `abi.encodePacked(args)`
func (u *Utils) AbiEncodePacked(args ...interface{}) ([]byte, error) {
	bytes := make([]byte, 0)
	for _, arg := range args {
		switch val := arg.(type) {
		case *big.Int:
			bytes = append(bytes, common.LeftPadBytes(val.Bytes(), 32)...)
		case bool:
			if val {
				bytes = append(bytes, []byte{0x0, 0x1}...)
			}
		case common.Hash:
			bytes = append(bytes, val[:]...)
		case []byte:
			bytes = append(bytes, val...)
		case common.Address:
			bytes = append(bytes, val[:]...)
		default:
			return nil, fmt.Errorf("unsupport type %T", arg)
		}
	}
	return bytes, nil
}
