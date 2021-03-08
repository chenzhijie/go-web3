<h1 >go-web3</h1>

Ethereum Golang API

## Development Environment
The requirements to develop are:

- [Golang](https://golang.org/doc/install) version 1.14 or later



## API

- [NewWeb3()](#NewWeb3)
- [SetChainId(chainId int64)](#SetChainId)
- [SetAccount(privateKey string) error](#SetAccount)
- [GetBlockNumber()](#GetBlockNumber)
- [GetNonce(addr common.Address, blockNumber *big.Int) (uint64, error)](#GetNonce)
- [NewContract(abiString string, contractAddr ...string) (*Contract, error)](#NewContract)
- [Call(methodName string, args ...interface{}) (interface{}, error)](#Call)

### NewWeb3()

Creates a new web3 instance with http provider.

```golang
// change to your rpc provider
var infuraURL = "https://mainnet.infura.io/v3/7238211010344719ad14a89db874158c"
web3, err := web3.NewWeb3(infuraURL)
if err != nil {
    panic(err)
}
```


### GetBlockNumber()

Get current block number.

```golang
blockNumber, err := web3.Eth.GetBlockNumber()
if err != nil {
    panic(err)
}
fmt.Println("Current block number: ", blockNumber)
// => Current block number:  11997285
```


### SetChainId(chainId int64)

Setup chainId for different network.

```golang
web3.Eth.SetChainId(1)
```


### SetAccount(privateKey string) error

Setup default account with privateKey (hex format)

```golang
err := web3.Eth.SetAccount("610ca682d9b48e079e9017bb000a503071a158941674d304efccc68d9b8756f9")
if err != nil {
    panic(err)
}
```


### GetNonce(addr common.Address, blockNumber *big.Int) (uint64, error)

Get transaction nonce for address

```golang
nonce, err := web3.Eth.GetNonce(web3.Eth.Address(), nil)
if err != nil {
    panic(err)
}
fmt.Println("Latest nonce: ", nonce)
// => Latest nonce: 1 
```

### NewContract(abiString string, contractAddr ...string) (*Contract, error)

Init contract api

```golang
abiString := "" // abi string
contractAddr := "" // contract address
contract, err := web3.Eth.NewContract(abiString, contractAddr )
if err != nil {
    panic(err)
}
```

### Call(methodName string, args ...interface{}) (interface{}, error)

Contract call method

```golang

totalSupply, err := contract.Call("totalSupply")
if err != nil {
    panic(err)
}
fmt.Printf("Total supply %v\n", totalSupply)

// => Total supply  10000000000
```



## License

The go-web3 source code is available under the [LGPL-3.0](LICENSE) license.