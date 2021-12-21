package utils

import (
	"bytes"
	"fmt"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
)

type EtherUnit int

const (
	EtherUnitNoEther EtherUnit = iota
	EtherUnitWei
	EtherUnitKWei
	EtherUnitMWei
	EtherUnitGWei
	EtherUnitSzabo
	EtherUnitFinney
	EtherUnitEther
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
	bigval.SetInt(wei)
	ret := bigval.Quo(bigval, expF)
	return ret
}

func (u *Utils) FromWeiFloat(wei *big.Float) *big.Float {
	exp := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	expF := new(big.Float)
	expF.SetInt(exp)

	bigval := new(big.Float)
	ret := bigval.Quo(wei, expF)
	return ret
}

func (u *Utils) FromDecimals(wei *big.Int, decimals int64) *big.Float {
	exp := new(big.Int).Exp(big.NewInt(10), big.NewInt(decimals), nil)
	expF := new(big.Float)
	expF.SetInt(exp)

	bigval := new(big.Float)
	bigval.SetInt(wei)
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

func (u *Utils) ToWeiInt(val int64, denominator int64) *big.Int {
	bigval := big.NewInt(val)
	exp := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)

	result := bigval.Mul(bigval, exp)
	result = big.NewInt(1).Div(result, big.NewInt(denominator))
	return result
}

func (u *Utils) FromWeiWithUnit(wei *big.Int, unit EtherUnit) *big.Float {
	unitInt := 0
	switch unit {
	case EtherUnitNoEther:
		unitInt = 0
	case EtherUnitWei:
		unitInt = 1
	case EtherUnitKWei:
		unitInt = 3
	case EtherUnitMWei:
		unitInt = 6
	case EtherUnitGWei:
		unitInt = 9
	case EtherUnitSzabo:
		unitInt = 12
	case EtherUnitFinney:
		unitInt = 15
	case EtherUnitEther:
		unit = 18
	}
	exp := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(unitInt)), nil)
	expF := new(big.Float)
	expF.SetInt(exp)

	bigval := new(big.Float)
	bigval.SetInt(wei)

	ret := bigval.Quo(bigval, expF)
	return ret
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
	return bytes.Equal(a[:], b[:])
}

func (u *Utils) DifferentAddress(a, b common.Address) bool {
	return !bytes.Equal(a[:], b[:])
}

func (u *Utils) RoundNWei(wei *big.Int, n int) (*big.Int, error) {
	af := u.FromWei(wei)
	aff, _ := af.Float64()

	roundfs := ""
	if n > 6 {
		return nil, fmt.Errorf("round n not support bigger than 6")
	}
	switch n {
	case 1:
		roundfs = fmt.Sprintf("%.1f", aff)

	case 2:
		roundfs = fmt.Sprintf("%.2f", aff)

	case 3:
		roundfs = fmt.Sprintf("%.3f", aff)

	case 4:
		roundfs = fmt.Sprintf("%.4f", aff)

	case 5:
		roundfs = fmt.Sprintf("%.5f", aff)

	case 6:
		roundfs = fmt.Sprintf("%.6f", aff)

	}

	roundf, err := strconv.ParseFloat(roundfs, 64)
	if err != nil {
		return nil, err
	}
	r := u.ToWei(roundf)

	return r, nil
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
