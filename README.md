<h1 >go-web3</h1>

Ethereum Golang API

## Development Environment
The requirements to develop are:

- [Golang](https://golang.org/doc/install) version 1.14 or later



## API

- [NewWeb3()](#NewWeb3)
- [SetChainId(chainId int64)](#setchainidchainid-int64)
- [SetAccount(privateKey string) error](#setaccountprivatekey-string-error)
- [GetBlockNumber()](#GetBlockNumber)
- [GetNonce(addr common.Address, blockNumber *big.Int) (uint64, error)](#getnonceaddr-commonaddress-blocknumber-bigint-uint64-error)
- [NewContract(abiString string, contractAddr ...string) (*Contract, error)](#newcontractabistring-string-contractaddr-string-contract-error)
- [Call(methodName string, args ...interface{}) (interface{}, error)](#callmethodname-string-args-interface-interface-error)
- [EncodeABI(methodName string, args ...interface{}) ([]byte, error)](#encodeabimethodname-string-args-interface-byte-error)
- [SendRawTransaction(to common.Address,amount *big.Int,gasLimit uint64,gasPrice *big.Int,data []byte) (common.Hash, error) ](#sendrawtransactionto-commonaddressamount-bigintgaslimit-uint64gasprice-bigintdata-byte-commonhash-error)

### NewWeb3()

Creates a new web3 instance with http provider.

```golang
// change to your rpc provider
var rpcProviderURL = "https://rpc.flashbots.net"
web3, err := web3.NewWeb3(rpcProviderURL)
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
pv, err := crypto.GenerateKey()
if err != nil {
    panic(err)
}
privateKey := hex.EncodeToString(crypto.FromECDSA(pv))
err := web3.Eth.SetAccount(privateKey)
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
abiString := `[
	{
		"constant": true,
		"inputs": [],
		"name": "totalSupply",
		"outputs": [
			{
				"name": "",
				"type": "uint256"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	}
]`
contractAddr := "0x6B175474E89094C44Da98b954EedeAC495271d0F" // contract address
contract, err := web3.Eth.NewContract(abiString, contractAddr)
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
- **[EIP1559 API](./examples/eip1559/main.go)**

## Contributing

Go-Web3 welcome contributions. Please follow the guidelines when opening issues and contributing code to the repo.

### Contributing

We follow the [fork-and-pull](https://help.github.com/en/articles/about-collaborative-development-models) Git workflow:

 1. **Fork** the repo on GitHub
 2. **Clone** it to your own machine
 3. **Commit** changes to your fork
 4. **Push** changes to your fork
 5. Submit a **Pull request** for review

**NOTE:** Be sure to merge the latest changes before making a pull request!

### Pull Requests

As outlined in Keavy McMinn's article ["How to write the perfect pull request"](https://github.blog/2015-01-21-how-to-write-the-perfect-pull-request/), you should include:

  1. The purpose of the PR
  2. A brief overview of what you did
  3. Tag any issues that the PR relates to and [close issues with a keyword](https://help.github.com/en/articles/closing-issues-using-keywords)
  4. What kind of feedback you're looking for (if any)
  5. Tag individuals you want feedback from (if any)

### Issues

Feel free to submit issues and enhancement requests [here](https://github.com/chenzhijie/go-web3/issues/new). Please consider [how to ask a good question](https://stackoverflow.com/help/how-to-ask) and take the time to research your issue before asking for help.

Duplicate questions will be closed.

## License

The go-web3 source code is available under the [LGPL-3.0](LICENSE) license.