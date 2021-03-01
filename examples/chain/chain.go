package main

import (
	"fmt"

	"github.com/chenzhijie/go-web3"
)

func main() {
	var infuraURL = "https://kovan.infura.io/v3/68ee32bee1584c37a4fe742846dc9772"
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
