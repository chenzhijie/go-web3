package main

import (
	"fmt"
	"os"

	"github.com/chenzhijie/go-web3"
)

func main() {
	// Open a terminal, and setup infura API key and ethereum privateKey to your env
	// $ export eth_infuraKey=YourInfuraAPIKey
	// $ export eth_privateKey=YourPrivateKey

	// change to your rpc provider
	var infuraURL = "https://rinkeby.infura.io/v3/" + os.Getenv("eth_infuraKey")

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
