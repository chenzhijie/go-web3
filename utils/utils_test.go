package utils

import (
	"bytes"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestToWei(t *testing.T) {
	u := NewUtils()
	values := []string{
		"12",
		"1",
		"1.0",
		"13.1111111111111111",
		"13.111111111111111111",
		"0.141",
		"13.1111111111111111111", // wrong
	}
	expects := []string{
		"12000000000000000000",
		"1000000000000000000",
		"1000000000000000000",
		"13111111111111111100",
		"13111111111111111111",
		"141000000000000000",
		"0",
	}
	for i, val := range values {
		bigVal, ok := big.NewInt(0).SetString(expects[i], 10)
		if !ok {
			t.Fatal(fmt.Sprintf("convert %s failed", val))
		}
		if !bytes.Equal(u.ToWei(val).Bytes(), bigVal.Bytes()) {
			t.Fatal(fmt.Sprintf("%s and %s not equal", val, bigVal))
		}
	}

}

func TestFromWei(t *testing.T) {
	ethVal := "0.00000001"
	u := NewUtils()
	wei := u.ToWei(ethVal)
	fmt.Printf("wei %v\n", wei)
	gwei := u.FromWeiWithUnit(wei, EtherUnitGWei)
	fmt.Printf("gwei %.4f\n", gwei)
}

func TestSignMethod(t *testing.T) {

	funcName := "freeMint()"
	id := crypto.Keccak256([]byte(funcName))[:4]
	fmt.Printf("id 0x%x\n", id)

	funcName = "mint()"
	id = crypto.Keccak256([]byte(funcName))[:4]
	fmt.Printf("id 0x%x\n", id)
}

func TestSameAddr(t *testing.T) {
	addr1 := common.HexToAddress("0x0000000000fC95fD88A4c46d9d7984A56289c52A")
	addr2 := common.HexToAddress("0x0000000000fc95fd88a4c46d9d7984a56289c52a")

	fmt.Printf("addr1 %v, addr2 %v\n", addr1, addr2)

	fmt.Printf("addr1 == addr2 %t\n", addr1 == addr2)

	fmt.Printf("bytes cmp addr1 == addr2 %t\n", bytes.Compare(addr1[:], addr2[:]) == 0)
}

func BenchmarkTestCompare1(b *testing.B) {
	addr1 := common.HexToAddress("0x0000000000fC95fD88A4c46d9d7984A56289c52A")
	addr2 := common.HexToAddress("0x0000000000fc95fd88a4c46d9d7984a56289c52a")

	for i := 0; i < b.N; i++ {
		if addr1 == addr2 {
			continue
		}
	}
}

func BenchmarkTestCompare2(b *testing.B) {
	addr1 := common.HexToAddress("0x0000000000fC95fD88A4c46d9d7984A56289c52A")
	addr2 := common.HexToAddress("0x0000000000fc95fd88a4c46d9d7984a56289c52a")

	for i := 0; i < b.N; i++ {
		if bytes.Compare(addr1[:], addr2[:]) == 0 {
			continue
		}
	}
}

func BenchmarkTestCompare3(b *testing.B) {
	addr1 := common.HexToAddress("0x0000000000fC95fD88A4c46d9d7984A56289c52A")
	addr2 := common.HexToAddress("0x0000000000fc95fd88a4c46d9d7984a56289c52a")

	for i := 0; i < b.N; i++ {
		if addr1.Hex() == addr2.Hex() {
			continue
		}
	}
}

func BenchmarkTestCompare4(b *testing.B) {
	addr1 := common.HexToAddress("0x0000000000fC95fD88A4c46d9d7984A56289c52A")
	addr2 := common.HexToAddress("0x0000000000fc95fd88a4c46d9d7984a56289c52a")

	for i := 0; i < b.N; i++ {
		if string(addr1[:]) == string(addr2[:]) {
			continue
		}
	}
}

func TestRoundNWei(t *testing.T) {

	u := Utils{}
	v := u.ToWei("0.073229374492")
	fmt.Printf(" %v\n", v)
	ret, err := u.RoundNWei(v, 5)
	if err != nil {
		panic(err)
	}
	fmt.Printf("ret %v\n", ret)
}
