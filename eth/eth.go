package eth

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/chenzhijie/go-web3/rpc"
	"github.com/chenzhijie/go-web3/types"
	"github.com/chenzhijie/go-web3/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"

	eTypes "github.com/ethereum/go-ethereum/core/types"
)

// Eth is the eth namespace
type Eth struct {
	c             *rpc.Client
	privateKey    *ecdsa.PrivateKey
	address       common.Address
	chainId       *big.Int
	txPollTimeout int
}

func NewEth(c *rpc.Client) *Eth {
	return &Eth{
		c: c,
	}
}

func (e *Eth) SetAccount(privateKey string) error {
	if len(privateKey) == 0 {
		return fmt.Errorf("private key is empty")
	}
	privKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return err
	}

	e.privateKey = privKey

	addr := crypto.PubkeyToAddress(privKey.PublicKey)
	// fmt.Println("addr ", addr)
	copy(e.address[:], addr[:])

	return nil
}

func (e *Eth) SetChainId(chainId int64) {
	e.chainId = big.NewInt(chainId)
}

func (e *Eth) SetTxPollTimeout(timeout int) {
	if timeout == 0 {
		// default tx poll timeout is 720s
		e.txPollTimeout = 720
		return
	}
	e.txPollTimeout = timeout
}

func (e *Eth) Accounts() ([]common.Address, error) {
	var out []common.Address
	if err := e.c.Call("eth_accounts", &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (e *Eth) Address() common.Address {
	return e.address
}

func (e *Eth) GetBlockNumber() (uint64, error) {
	var out string
	if err := e.c.Call("eth_blockNumber", &out); err != nil {
		return 0, err
	}
	return utils.ParseUint64orHex(out)
}

func (e *Eth) GetBlockByNumber(i *big.Int, full bool) (*eTypes.Block, error) {
	var b *eTypes.Block
	if err := e.c.Call("eth_getBlockByNumber", &b, i.String(), full); err != nil {
		return nil, err
	}
	return b, nil
}

func (e *Eth) GetBlockByHash(hash common.Hash, full bool) (*eTypes.Block, error) {
	var b *eTypes.Block
	if err := e.c.Call("eth_getBlockByHash", &b, hash, full); err != nil {
		return nil, err
	}
	return b, nil
}

func (e *Eth) SendTransaction(txn *eTypes.Transaction) (common.Hash, error) {
	var hash common.Hash
	err := e.c.Call("eth_sendTransaction", &hash, txn)
	return hash, err
}

func (e *Eth) GetTransactionReceipt(hash common.Hash) (*eTypes.Receipt, error) {
	var receipt *eTypes.Receipt
	err := e.c.Call("eth_getTransactionReceipt", &receipt, hash)
	return receipt, err
}

func (e *Eth) GetNonce(addr common.Address, blockNumber *big.Int) (uint64, error) {
	var nonce string
	if err := e.c.Call("eth_getTransactionCount", &nonce, addr, toBlockNumArg(blockNumber)); err != nil {
		return 0, err
	}
	return utils.ParseUint64orHex(nonce)
}

func (e *Eth) GetBalance(addr common.Address, blockNumber *big.Int) (*big.Int, error) {
	var out string
	if err := e.c.Call("eth_getBalance", &out, addr, toBlockNumArg(blockNumber)); err != nil {
		return nil, err
	}
	b, ok := new(big.Int).SetString(out[2:], 16)
	if !ok {
		return nil, fmt.Errorf("failed to convert to big.int")
	}
	return b, nil
}

func (e *Eth) GasPrice() (uint64, error) {
	var out string
	if err := e.c.Call("eth_gasPrice", &out); err != nil {
		return 0, err
	}
	return utils.ParseUint64orHex(out)
}

func (e *Eth) Call(msg *types.CallMsg, block *big.Int) (string, error) {
	var out string
	if err := e.c.Call("eth_call", &out, msg, block.String()); err != nil {
		return "", err
	}
	return out, nil
}

func (e *Eth) EstimateGasContract(bin []byte) (uint64, error) {
	var out string
	msg := map[string]interface{}{
		"data": "0x" + hex.EncodeToString(bin),
	}
	if err := e.c.Call("eth_estimateGas", &out, msg); err != nil {
		return 0, err
	}
	return utils.ParseUint64orHex(out)
}

func (e *Eth) EstimateGas(msg *types.CallMsg) (uint64, error) {
	var out string
	if err := e.c.Call("eth_estimateGas", &out, msg); err != nil {
		return 0, err
	}
	return utils.ParseUint64orHex(out)
}

func (e *Eth) ChainID() (*big.Int, error) {
	if e.chainId != nil {
		return e.chainId, nil
	}
	var out string
	if err := e.c.Call("eth_chainId", &out); err != nil {
		return nil, err
	}
	return utils.ParseBigInt(out), nil
}

func (e *Eth) EncodeParams() {
}

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	pending := big.NewInt(-1)
	if number.Cmp(pending) == 0 {
		return "pending"
	}
	return hexutil.EncodeBig(number)
}
