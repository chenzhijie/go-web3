<h1 >go-web3</h1>

Ethereum Golang API

## Development Environment
The requirements to develop are:

- [Golang](https://golang.org/doc/install) version 1.11 or later



## API

- [NewWeb3()](#NewWeb3)
- [SetChainId(chainId int64)](#SetChainId)
- [SetAccount(privateKey string) error](#SetAccount)
- [GetBlockNumber()](#GetBlockNumber)

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
// => Current block number:  11997285
```


### SetChainId(chainId int64)

Setup chainId for different network.

```golang
// change to your rpc provider
var infuraURL = "https://mainnet.infura.io/v3/7238211010344719ad14a89db874158c"
web3, err := web3.NewWeb3(infuraURL)
if err != nil {
    panic(err)
}
web3.Eth.SetChainId(1)
```


### SetAccount(privateKey string) error

Setup default account with privateKey (hex format)

```golang
// change to your rpc provider
web3, err := web3.NewWeb3(infuraURL)
if err != nil {
    panic(err)
}
err := web3.Eth.SetAccount("610ca682d9b48e079e9017bb000a503071a158941674d304efccc68d9b8756f9")
if err != nil {
    panic(err)
}
```



## License

The go-web3 source code is available under the [LGPL-3.0](LICENSE) license.