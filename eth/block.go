package eth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type rpcBlock struct {
	Hash         common.Hash          `json:"hash"`
	Transactions []*types.Transaction `json:"transactions"`
	UncleHashes  []common.Hash        `json:"uncles"`
}

func (e *Eth) getBlock(method string, args ...interface{}) (*types.Block, error) {
	var raw json.RawMessage
	err := e.c.Call(method, &raw, args...)
	if err != nil {
		return nil, err
	} else if len(raw) == 0 {
		return nil, ethereum.NotFound
	} else if bytes.Equal([]byte(raw), []byte("null")) {
		return nil, ethereum.NotFound
	}

	// Decode header and transactions.
	var head *types.Header
	var body rpcBlock
	if err := json.Unmarshal(raw, &head); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(raw, &body); err != nil {
		return nil, err
	}
	if head == nil {
		return nil, errors.New("json.Unmarshal header failed")
	}
	if head.UncleHash == types.EmptyUncleHash && len(body.UncleHashes) > 0 {
		return nil, fmt.Errorf("server returned non-empty uncle list but block header indicates no uncles")
	}
	if head.UncleHash != types.EmptyUncleHash && len(body.UncleHashes) == 0 {
		return nil, fmt.Errorf("server returned empty uncle list but block header indicates uncles")
	}
	if head.TxHash == types.EmptyRootHash && len(body.Transactions) > 0 {
		return nil, fmt.Errorf("server returned non-empty transaction list but block header indicates no transactions")
	}
	if head.TxHash != types.EmptyRootHash && len(body.Transactions) == 0 {
		return nil, fmt.Errorf("server returned empty transaction list but block header indicates transactions")
	}
	// Load uncles because they are not included in the block response.
	var uncles []*types.Header
	// TODO
	// if len(body.UncleHashes) > 0 {
	// 	uncles = make([]*types.Header, len(body.UncleHashes))
	// 	reqs := make([]rpc.BatchElem, len(body.UncleHashes))
	// 	for i := range reqs {
	// 		reqs[i] = rpc.BatchElem{
	// 			Method: "eth_getUncleByBlockHashAndIndex",
	// 			Args:   []interface{}{body.Hash, hexutil.EncodeUint64(uint64(i))},
	// 			Result: &uncles[i],
	// 		}
	// 	}
	// 	if err := e.c.BatchCallContext(ctx, reqs); err != nil {
	// 		return nil, err
	// 	}
	// 	for i := range reqs {
	// 		if reqs[i].Error != nil {
	// 			return nil, reqs[i].Error
	// 		}
	// 		if uncles[i] == nil {
	// 			return nil, fmt.Errorf("got null header for uncle %d of block %x", i, body.Hash[:])
	// 		}
	// 	}
	// }
	// Fill the sender cache of transactions in the block.
	return types.NewBlockWithHeader(head).WithBody(body.Transactions, uncles), nil
}
