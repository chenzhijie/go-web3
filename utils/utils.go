package utils

import (
	"bytes"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
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

var unitMap = map[EtherUnit]string{
	EtherUnitNoEther: "0",
	EtherUnitWei:     "1",
	EtherUnitKWei:    "1000",
	EtherUnitMWei:    "1000000",
	EtherUnitGWei:    "1000000000",
	EtherUnitSzabo:   "1000000000000",
	EtherUnitFinney:  "1000000000000000",
	EtherUnitEther:   "1000000000000000000",
}

type Utils struct{}

func NewUtils() *Utils {
	return &Utils{}
}

func (u *Utils) FromWei(wei *big.Int) *big.Float {
	if wei == nil {
		return big.NewFloat(0)
	}
	exp := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	expF := new(big.Float)
	expF.SetInt(exp)

	bigval := new(big.Float)
	bigval.SetInt(wei)
	ret := bigval.Quo(bigval, expF)
	return ret
}

func (u *Utils) FromGWei(wei *big.Int) *big.Float {
	if wei == nil {
		return big.NewFloat(0)
	}
	exp := new(big.Int).Exp(big.NewInt(10), big.NewInt(9), nil)
	expF := new(big.Float)
	expF.SetInt(exp)

	bigval := new(big.Float)
	bigval.SetInt(wei)
	ret := bigval.Quo(bigval, expF)
	return ret
}

func (u *Utils) FromWeiFloat(wei *big.Float) *big.Float {
	if wei == nil {
		return big.NewFloat(0)
	}
	exp := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	expF := new(big.Float)
	expF.SetInt(exp)

	bigval := new(big.Float)
	ret := bigval.Quo(wei, expF)
	return ret
}

func (u *Utils) FromDecimals(wei *big.Int, decimals int64) *big.Float {
	if wei == nil {
		return big.NewFloat(0)
	}
	exp := new(big.Int).Exp(big.NewInt(10), big.NewInt(decimals), nil)
	expF := new(big.Float)
	expF.SetInt(exp)

	bigval := new(big.Float)
	bigval.SetInt(wei)
	ret := bigval.Quo(bigval, expF)
	return ret
}

func (u *Utils) ToWei(val string) *big.Int {

	if !strings.Contains(val, ".") {
		whole, ok := big.NewInt(0).SetString(val, 10)
		if !ok {
			return big.NewInt(0)
		}
		return big.NewInt(1).Mul(whole, big.NewInt(1e18))
	}
	comps := strings.Split(val, ".")
	if len(comps) != 2 {
		return big.NewInt(0)
	}

	whole := comps[0]
	fraction := comps[1]
	baseLength := len(unitMap[EtherUnitEther]) - 1
	fractionLength := len(fraction)
	if fractionLength > baseLength {
		return big.NewInt(0)
	}
	fraction += strings.Repeat("0", baseLength-fractionLength)
	wholeInt, ok := big.NewInt(0).SetString(whole, 10)
	if !ok {
		return big.NewInt(0)
	}
	fractionInt, ok := big.NewInt(0).SetString(fraction, 10)
	if !ok {
		return big.NewInt(0)
	}

	wholeMulBase := big.NewInt(1).Mul(wholeInt, big.NewInt(1e18))
	wholeAddFraction := big.NewInt(1).Add(wholeMulBase, fractionInt)

	return wholeAddFraction
}

func (u *Utils) ToWeiInt(val int64, denominator int64) *big.Int {
	bigval := big.NewInt(val)
	exp := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)

	result := bigval.Mul(bigval, exp)
	result = big.NewInt(1).Div(result, big.NewInt(denominator))
	return result
}

func (u *Utils) FromWeiWithUnit(wei *big.Int, unit EtherUnit) *big.Float {
	if wei == nil {
		return big.NewFloat(0)
	}
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

	r := u.ToWei(roundfs)

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

// ToBlockNumArg. Wrap blockNumber arg from big.Int to string
func ToBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	pending := big.NewInt(-1)
	if number.Cmp(pending) == 0 {
		return "pending"
	}
	return hexutil.EncodeBig(number)
}
