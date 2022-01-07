package eth

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/chenzhijie/go-web3/rpc"
)

func TestGetBlockByNumber(t *testing.T) {
	c, err := rpc.NewClient("https://rpc.flashbots.net", "http://127.0.0.1:7890")
	if err != nil {
		t.Fatal(err)
	}
	eth := NewEth(c)
	blockNumber, err := eth.GetBlockNumber()
	if err != nil {
		t.Fatal(err)
	}
	block, err := eth.GetBlocByNumber(big.NewInt(int64(blockNumber)), true)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("block hash %s has %v txs\n", block.Hash(), len(block.Transactions()))
}

func TestPollBlock(t *testing.T) {
	c, err := rpc.NewClient("https://rpc.flashbots.net", "http://127.0.0.1:7890")
	if err != nil {
		t.Fatal(err)
	}
	eth := NewEth(c)
	for {
		blockNumber, err := eth.GetBlockNumber()
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("get block %v\n", blockNumber)
		block, err := eth.GetBlocByNumber(big.NewInt(int64(blockNumber)), true)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("block hash %s has %v txs\n", block.Hash(), len(block.Transactions()))
		time.Sleep(time.Duration(5) * time.Second)
	}

}
