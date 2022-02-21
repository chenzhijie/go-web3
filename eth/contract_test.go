package eth

import (
	"fmt"
	"testing"

	"github.com/chenzhijie/go-web3/rpc"
	"github.com/chenzhijie/go-web3/types"
	"github.com/chenzhijie/go-web3/utils"
	"github.com/ethereum/go-ethereum/common"
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

func TestCallWithMethodSignature(t *testing.T) {
	web3Utils := &utils.Utils{}
	methodSignature := web3Utils.EncodeFunctionSignature("factory()")
	c, err := rpc.NewClient("https://emerald.oasis.dev", "http://127.0.0.1:7890")
	if err != nil {
		t.Fatal(err)
	}
	e := NewEth(c)
	result, err := e.Call(&types.CallMsg{
		To:   common.HexToAddress("0x250d48C5E78f1E85F7AB07FEC61E93ba703aE668"),
		Data: methodSignature,
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	addr := common.HexToAddress(result)
	fmt.Printf("addr %v, type %T\n", addr, result)
}
