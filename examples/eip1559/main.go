package main

import (
	"fmt"
	"os"

	"github.com/chenzhijie/go-web3"
	"github.com/chenzhijie/go-web3/utils"
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

	web3.Eth.SetAccount(os.Getenv("eth_privateKey"))
	blockNumber, err := web3.Eth.GetBlockNumber()
	if err != nil {
		panic(err)
	}
	fmt.Println("Current block number: ", blockNumber)
	fee, err := web3.Eth.EstimateFee()
	if err != nil {
		panic(err)
	}
	util := utils.Utils{}
	fmt.Printf("base fee %v, %.3f Gwei\n", fee.BaseFee, util.FromWeiWithUnit(fee.BaseFee, utils.EtherUnitGWei))
	fmt.Printf("max priority fee per gas %v, %.3f Gwei\n", fee.MaxPriorityFeePerGas, util.FromWeiWithUnit(fee.MaxPriorityFeePerGas, utils.EtherUnitGWei))
	fmt.Printf("max fee per gas %v, %.3f Gwei\n", fee.MaxFeePerGas, util.FromWeiWithUnit(fee.MaxFeePerGas, utils.EtherUnitGWei))

	fmt.Println("current account ", web3.Eth.Address())

	receipt, err := web3.Eth.SyncSendEIP1559RawTransaction(web3.Eth.Address(), util.ToWei(0.01), 21000, fee.MaxPriorityFeePerGas, fee.MaxFeePerGas, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("eip 1559 tx %v\n", receipt)

}
