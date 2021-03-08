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
- [EncodeABI(methodName string, args ...interface{}) ([]byte, error)](#EncodeABI)
- [SendRawTransaction(to common.Address,amount *big.Int,gasLimit uint64,gasPrice *big.Int,data []byte) (common.Hash, error) ](#SendRawTransaction)

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

### EncodeABI(methodName string, args ...interface{}) ([]byte, error)

EncodeABI data

```golang

data, err := contract.EncodeABI("balanceOf", web3.Eth.Address())
if err != nil {
    panic(err)
}
fmt.Printf("Data %x\n", data)

// => Data 70a08231000000000000000000000000c13a163dd812ed7eb8bb9152651054eae5ee0999 
```

### SendRawTransaction(to common.Address,amount *big.Int,gasLimit uint64,gasPrice *big.Int,data []byte) (common.Hash, error) 

Send transaction

```golang

txHash, err := web3.Eth.SendRawTransaction(
    common.HexToAddress(tokenAddr),
    big.NewInt(0),
    gasLimit,
    web3.Utils.ToGWei(1),
    approveInputData,
)
if err != nil {
    panic(err)
}
fmt.Printf("Send approve tx hash %v\n", txHash)

// => Send approve tx hash  0x837136c8b6f34b519c049d1cf703d3bba47d32f6801c25d83d0113bdc0e6936a 
```

## Examples

- **[Chain API](./examples/chain/chain.go)**
- **[Contract API](./examples/contract/erc20.go)**

## License

The go-web3 source code is available under the [LGPL-3.0](LICENSE) license.