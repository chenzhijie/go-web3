package types

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type FeeHistory struct {
	BaseFeePerGas []*hexutil.Big   `json:"baseFeePerGas"`
	GasUsedRatio  []float64        `json:"gasUsedRatio"`
	OldestBlock   *hexutil.Big     `json:"oldestBlock"`
	Reward        [][]*hexutil.Big `json:"reward"`
}

type Bigs []*hexutil.Big

func (s Bigs) Len() int {
	return len(s)
}

func (s Bigs) Less(i, j int) bool {
	return s[i].ToInt().Cmp(s[j].ToInt()) < 0
}

func (s Bigs) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
