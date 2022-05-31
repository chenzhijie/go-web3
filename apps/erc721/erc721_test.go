package erc721

import (
	"fmt"
	"testing"

	"github.com/chenzhijie/go-web3"
	"github.com/ethereum/go-ethereum/common"
)

func TestIsApprovedForAll(t *testing.T) {
	web3, err := web3.NewWeb3WithProxy("https://rpc.flashbots.net", "http://127.0.0.1:7890")
	if err != nil {
		t.Fatal(err)
	}
	erc721, err := NewERC721(web3, common.HexToAddress("0x60e4d786628fea6478f785a6d7e704777c86a7c6"))
	if err != nil {
		t.Fatal(err)
	}
	// https://etherscan.io/tx/0x99b59ad25c4cc72e7606268e15a52aafdd4ff0a7a8f4f91a1ec72cd1653568fe
	owner := common.HexToAddress("0x1cb162bfd68a0757f8de8659ffe1768c42bd92cb")
	spender := common.HexToAddress("0x37D3341f56a03119c57EA3010cAD46b85cC560Ee")
	approved, err := erc721.IsApprovedForAll(owner, spender)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("approved %t\n", approved)
}
