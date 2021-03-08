<h1 align="center">go-web3</h1>
<h4 align="center">Version 0.1.0</h4>

Welcome to the Ethereum Golang API

## Development Environment
The requirements to develop are:

- [Golang](https://golang.org/doc/install) version 1.11 or later



## API

- [NewWeb3()](#NewWeb3)

### NewWeb3()

Creates a new web3 instance with http provider.

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


## License

The go-web3 source code is available under the [LGPL-3.0](LICENSE) license.