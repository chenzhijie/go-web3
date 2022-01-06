package eth

import (
	"fmt"
	"testing"

	"github.com/chenzhijie/go-web3/rpc"
)

func TestContractCall(t *testing.T) {
	abi := `[
		{
			"inputs": [],
			"name": "decimals",
			"outputs": [
				{
					"internalType": "uint8",
					"name": "",
					"type": "uint8"
				}
			],
			"stateMutability": "pure",
			"type": "function"
		},
		{
			"inputs": [],
			"name": "getReserves",
			"outputs": [
				{
					"internalType": "uint112",
					"name": "_reserve0",
					"type": "uint112"
				},
				{
					"internalType": "uint112",
					"name": "_reserve1",
					"type": "uint112"
				},
				{
					"internalType": "uint32",
					"name": "_blockTimestampLast",
					"type": "uint32"
				}
			],
			"stateMutability": "view",
			"type": "function"
		}
	]`
	c, err := rpc.NewClient("https://rpc.flashbots.net", "")
	if err != nil {
		t.Fatal(err)
	}
	eth := NewEth(c)
	uniswapV2PairContr, err := eth.NewContract(abi, "0x0d4a11d5EEaaC28EC3F61d100daF4d40471f1852")
	if err != nil {
		t.Fatal(err)
	}
	reserves, err := uniswapV2PairContr.Call("getReserves")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("reserves %v, type %T\n", reserves, reserves)

	decimals, err := uniswapV2PairContr.Call("decimals")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("decimals %v, type %T\n", decimals, decimals)
}
