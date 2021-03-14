package utils

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestToWei(t *testing.T) {
	ethVal := 0.003
	u := NewUtils()
	wei := u.ToWei(ethVal)
	fmt.Printf("wei %v\n", wei)
	fmt.Printf("wei hex %v\n", u.ToHex(wei))
}

func TestFromWei(t *testing.T) {
	ethVal := 0.00000001
	u := NewUtils()
	wei := u.ToWei(ethVal)
	fmt.Printf("wei %v\n", wei)
	gwei := u.FromWeiWithUnit(wei, EtherUnitGWei)
	fmt.Printf("gwei %.4f\n", gwei)
}

func TestSignMethod(t *testing.T) {

	funcName := "transfer(address,uint256)"
	id := crypto.Keccak256([]byte(funcName))[:4]
	fmt.Printf("id %x\n", id)
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
