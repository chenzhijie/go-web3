package main

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"os"

	"github.com/chenzhijie/go-web3"
	"github.com/ethereum/go-ethereum/crypto"
)

func main() {

	// change to your rpc provider
	var chainId = int64(97)
	var rpcProvider = "https://data-seed-prebsc-2-s3.binance.org:8545/"
	web3, err := web3.NewWeb3WithProxy(rpcProvider, os.Getenv("http_proxy"))
	if err != nil {
		panic(err)
	}
	web3.Eth.SetChainId(chainId)
	blockNumber, err := web3.Eth.GetBlockNumber()
	if err != nil {
		panic(err)
	}
	fmt.Println("Current block number: ", blockNumber)

	// only for test
	privateKey := "" // hex string format
	if len(privateKey) == 0 {
		panic("please replace to your privateKey and keep safe by yourself")
	}
	// setup your privateKey
	if err := web3.Eth.SetAccount(privateKey); err != nil {
		panic(err)
	}
	privateKeyData, err := hex.DecodeString(privateKey)
	if err != nil {
		panic(err)
	}
	ecdsaPrivateKey, err := crypto.ToECDSA(privateKeyData)
	if err != nil {
		panic(err)
	}

	addr := crypto.PubkeyToAddress(ecdsaPrivateKey.PublicKey)
	fmt.Printf("Address %s\n", addr)

	bnbBalance, err := web3.Eth.GetBalance(addr, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Bnb balance %v\n", bnbBalance)

	gasPrice, err := web3.Eth.GasPrice()
	if err != nil {
		panic(err)
	}
	fmt.Printf("GasPrice %v\n", gasPrice)
	nonce, err := web3.Eth.GetNonce(web3.Eth.Address(), nil)
	if err != nil {
		panic(err)
	}
	tx, err := web3.Eth.SendRawTransaction(
		addr,
		web3.Utils.ToWei("0.1"),
		nonce,
		21000,
		big.NewInt(int64(gasPrice)),
		nil,
	)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Send 0.1 BNB to self tx %s\n", tx)
}
