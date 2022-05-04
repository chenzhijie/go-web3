package crypto

import (
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
)

func Keccak256Hash(data []byte) []byte {
	return ethCrypto.Keccak256(data)
}
