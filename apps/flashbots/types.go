package flashbots

import (
	"encoding/json"
	"errors"
	"math/big"
)

const (
	DefaultRelayURL = "https://relay.flashbots.net"
	TestRelayURL    = "https://relay-goerli.flashbots.net"
)

type BundleResult struct {
	BundleHash string `json:"bundleHash"`
}

type txResult struct {
	CoinbaseDiff      string `json:"coinbaseDiff"`
	EthSentToCoinbase string `json:"ethSentToCoinbase"`
	FromAddress       string `json:"fromAddress"`
	GasFees           string `json:"gasFees"`
	GasPrice          string `json:"gasPrice"`
	GasUsed           uint64 `json:"gasUsed"`
	ToAddress         string `json:"toAddress"`
	TxHash            string `json:"txHash"`
	Value             string `json:"value"`
	Error             string `json:"error,omitempty"`
}

type CallResult struct {
	BundleGasPrice    string     `json:"bundleGasPrice"`
	BundleHash        string     `json:"bundleHash"`
	CoinbaseDiff      string     `json:"coinbaseDiff"`
	EthSentToCoinbase string     `json:"ethSentToCoinbase"`
	GasFees           string     `json:"gasFees"`
	Results           []txResult `json:"results"`
	StateBlockNumber  uint64     `json:"stateBlockNumber"`
	TotalGasUsed      uint64     `json:"totalGasUsed"`
}

func (cr *CallResult) String() string {
	r, err := json.Marshal(cr)
	if err != nil {
		return err.Error()
	}
	return string(r)
}

func (r *CallResult) EffectiveGasPrice() (*big.Int, error) {
	gu := new(big.Int).SetUint64(r.TotalGasUsed)
	gp, ok := new(big.Int).SetString(r.CoinbaseDiff, 10)
	if !ok {
		return nil, errors.New("invalid value returned for CoinbaseDiff")
	}
	wei := new(big.Int).Div(gp, gu)
	return wei, nil
}

type UserStats struct {
	IsHighPriority       bool   `json:"is_high_priority"`
	AllTimeMinerPayments string `json:"all_time_miner_payments"`
	AllTimeGasSimulated  string `json:"all_time_gas_simulated"`
	Last7dMinerPayments  string `json:"last_7d_miner_payments"`
	Last7dGasSimulated   string `json:"last_7d_gas_simulated"`
	Last1dMinerPayments  string `json:"last_1d_miner_payments"`
	Last1dGasSimulated   string `json:"last_1d_gas_simulated"`
}

type BundleStats struct {
	IsHighPriority bool   `json:"is_high_priority"`
	IsSimulated    bool   `json:"is_simulated"`
	IsSentToMiners bool   `json:"is_sent_to_miners"`
	SimulatedAt    string `json:"simulated_at"`
	SubmittedAt    string `json:"submitted_at"`
	SentToMinersAt string `json:"sent_to_miners_at"`
}
