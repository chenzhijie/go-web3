package main

import (
	"fmt"

	"github.com/chenzhijie/go-web3"
)

func main() {
	// change to your rpc provider
	var infuraURL = "https://mainnet.infura.io/v3/7238211010344719ad14a89db874158c"
	web3, err := web3.NewWeb3(infuraURL)
	if err != nil {
		panic(err)
	}
	blockNumber, err := web3.Eth.GetBlockNumber()
	if err != nil {
		panic(err)
	}
	fmt.Println("Current block number: ", blockNumber)
}
