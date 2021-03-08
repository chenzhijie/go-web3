package utils

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Utils struct{}

func NewUtils() *Utils {
	return &Utils{}
}

func (u *Utils) FromWei(wei *big.Int) *big.Float {
	exp := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	expF := new(big.Float)
	expF.SetInt(exp)

	bigval := new(big.Float)
	bigval.SetUint64(wei.Uint64())

	ret := bigval.Quo(bigval, expF)
	return ret
}

func (u *Utils) ToWei(val float64) *big.Int {
	bigval := new(big.Float)
	bigval.SetFloat64(val)
	exp := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)

	expF := new(big.Float)
	expF.SetInt(exp)

	bigval.Mul(bigval, expF)

	result := new(big.Int)
	bigval.Int(result) // store converted number in result

	return result
}

func (u *Utils) ToGWei(val float64) *big.Int {
	bigval := new(big.Float)
	bigval.SetFloat64(val)
	exp := new(big.Int).Exp(big.NewInt(10), big.NewInt(9), nil)

	expF := new(big.Float)
	expF.SetInt(exp)

	bigval.Mul(bigval, expF)

	result := new(big.Int)
	bigval.Int(result) // store converted number in result

	return result
}

func (u *Utils) ToHex(n *big.Int) string {
	return fmt.Sprintf("0x%x", n) // or %X or upper case
}

func (u *Utils) HexToUint64(str string) (uint64, error) {
	return ParseUint64orHex(str)
}

func (u *Utils) ToDecimals(val uint64, decimals int64) *big.Int {
	return convert(val, decimals)
}

func (u *Utils) SameAddress(a, b common.Address) bool {
	return bytes.Compare(a[:], b[:]) == 0
}

func (u *Utils) DifferentAddress(a, b common.Address) bool {
	return bytes.Compare(a[:], b[:]) != 0
}

// Ether converts a value to the ether unit with 18 decimals
func Ether(i uint64) *big.Int {
	return convert(i, 18)
}

func convert(val uint64, decimals int64) *big.Int {
	v := big.NewInt(int64(val))
	exp := new(big.Int).Exp(big.NewInt(10), big.NewInt(decimals), nil)
	return v.Mul(v, exp)
}
