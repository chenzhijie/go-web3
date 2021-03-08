package utils

import "github.com/ethereum/go-ethereum/common"

func (u *Utils) LeftPadBytes(slice []byte, l int) []byte {
	return common.LeftPadBytes(slice, l)
}
