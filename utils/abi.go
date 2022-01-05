package utils

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
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
