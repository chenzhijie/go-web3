package utils

import "github.com/ethereum/go-ethereum/crypto"

func (u *Utils) EncodeFunctionSignature(funcName string) []byte {
	return crypto.Keccak256([]byte(funcName))[:4]
}
