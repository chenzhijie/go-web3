package flashbots

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/chenzhijie/go-web3/rpc/transport"
	"github.com/ethereum/go-ethereum/common/hexutil"
	eTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/valyala/fasthttp"
)

type FlashBot struct {
	httpClient       *transport.HTTP
	providerURL      string
	signerPrivateKey *ecdsa.PrivateKey
}

// NewFlashBot, init flashbot instance
// providerURL: flashbot provider http url
// rpcRequestSigner: private key for signing request message
// proxyAddr: http proxy or socks5 proxy addr
func NewFlashBot(providerURL, rpcRequestSigner string) (*FlashBot, error) {

	signerKeyData, err := hex.DecodeString(strings.TrimPrefix(rpcRequestSigner, "0x"))
	if err != nil {
		return nil, err
	}
	signerPrivateKey, err := crypto.ToECDSA(signerKeyData)
	if err != nil {
		return nil, err
	}

	httpClient := transport.NewHTTP(providerURL, os.Getenv("http_proxy"))
	fb := &FlashBot{
		httpClient:       httpClient,
		providerURL:      providerURL,
		signerPrivateKey: signerPrivateKey,
	}
	return fb, nil
}

func (fb *FlashBot) SendBundle(txs []*eTypes.Transaction, targetBlockNumber *big.Int) (*BundleResult, error) {
	bundles := make([]string, 0)
	for _, tx := range txs {
		txData, err := tx.MarshalBinary()
		if err != nil {
			return nil, err
		}
		bundles = append(bundles, hexutil.Encode(txData))
	}
	return fb.SendRawBundle(bundles, targetBlockNumber)
}

func (fb *FlashBot) SendRawBundle(transactions []string, targetBlockNumber *big.Int) (*BundleResult, error) {

	type SendBundleParams struct {
		Transactions      []string `json:"txs"`
		BlockNumber       string   `json:"blockNumber"`
		MinTimestamp      int64    `json:"minTimestamp,omitempty"`
		MaxTimestamp      int64    `json:"maxTimestamp,omitempty"`
		RevertingTxHashes []string `json:"revertingTxHashes,omitempty"`
	}

	params := SendBundleParams{
		Transactions: transactions,
		BlockNumber:  fmt.Sprintf("0x%x", targetBlockNumber),
	}

	httpResp, err := fb.sendRequest(
		fb.providerURL,
		"eth_sendBundle",
		[]interface{}{params},
	)
	if err != nil {
		return nil, err
	}
	var resp struct {
		ID     uint64       `json:"id"`
		Error  interface{}  `json:"error,omitempty"`
		Result BundleResult `json:"result"`
	}

	err = json.Unmarshal(httpResp, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		errMsg, ok := resp.Error.(string)
		if ok {
			return nil, errors.New(errMsg)
		}

		return nil, fmt.Errorf("send request err %v", resp.Error)

	}
	return &resp.Result, nil
}

func (fb *FlashBot) CallRawBundle(transactions []string, blockNumber *big.Int, stateBlockNumber string) (*CallResult, error) {

	type CallBundleParams struct {
		Transactions     []string `json:"txs"`
		BlockNumber      string   `json:"blockNumber"`
		StateBlockNumber string   `json:"stateBlockNumber"`
	}

	params := CallBundleParams{
		Transactions:     transactions,
		BlockNumber:      fmt.Sprintf("0x%x", blockNumber.Uint64()),
		StateBlockNumber: stateBlockNumber,
	}

	httpResp, err := fb.sendRequest(fb.providerURL, "eth_callBundle", []interface{}{params})
	if err != nil {
		return nil, err
	}

	var resp struct {
		ID     uint64      `json:"id"`
		Error  interface{} `json:"error,omitempty"`
		Result CallResult  `json:"result"`
	}

	err = json.Unmarshal(httpResp, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		errMsg, ok := resp.Error.(string)
		if ok {
			return nil, errors.New(errMsg)
		}

		return nil, fmt.Errorf("send request err %v", resp.Error)

	}
	return &resp.Result, nil
}

func (fb *FlashBot) GetUserStats(blockNumber *big.Int) (*UserStats, error) {

	httpResp, err := fb.sendRequest(fb.providerURL, "flashbots_getUserStats", []interface{}{fmt.Sprintf("0x%x", blockNumber.Uint64())})
	if err != nil {
		return nil, err
	}

	var resp struct {
		ID     uint64      `json:"id"`
		Error  interface{} `json:"error,omitempty"`
		Result UserStats   `json:"result"`
	}

	err = json.Unmarshal(httpResp, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		errMsg, ok := resp.Error.(string)
		if ok {
			return nil, errors.New(errMsg)
		}

		return nil, fmt.Errorf("send request err %v", resp.Error)

	}
	return &resp.Result, nil
}

func (fb *FlashBot) GetBunderStats(bundleHash string, blockNumber *big.Int) (*BundleStats, error) {

	type reqParam struct {
		BundleHash  string `json:"bundleHash"`
		BlockNumber string `json:"blockNumber"`
	}
	param := &reqParam{
		BundleHash:  bundleHash,
		BlockNumber: fmt.Sprintf("0x%x", blockNumber.Uint64()),
	}

	httpResp, err := fb.sendRequest(fb.providerURL,
		"flashbots_getBundleStats",
		[]interface{}{
			param,
		})
	if err != nil {
		return nil, err
	}

	fmt.Printf("httpResp %s\n", httpResp)
	var resp struct {
		ID     uint64      `json:"id"`
		Error  interface{} `json:"error,omitempty"`
		Result BundleStats `json:"result"`
	}

	err = json.Unmarshal(httpResp, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		errMsg, ok := resp.Error.(string)
		if ok {
			return nil, errors.New(errMsg)
		}

		return nil, fmt.Errorf("send request err %v", resp.Error)

	}
	return &resp.Result, nil
}

func (fb *FlashBot) Simulate(
	txs []*eTypes.Transaction,
	blockNumber *big.Int,
	stateBlockNumber string,
) (*CallResult, error) {
	transactions := make([]string, 0)
	for _, tx := range txs {
		txData, err := tx.MarshalBinary()
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, hexutil.Encode(txData))
	}

	return fb.CallRawBundle(transactions, blockNumber, stateBlockNumber)
}

func (fb *FlashBot) SimulateRaw(
	transactions []string,
	blockNumber *big.Int,
	stateBlockNumber string,
) (*CallResult, error) {
	return fb.CallRawBundle(transactions, blockNumber, stateBlockNumber)
}

func (fb *FlashBot) sendRequest(relay string, method string, params []interface{}) ([]byte, error) {

	payload, err := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  method,
		"params":  params,
	})
	if err != nil {
		return nil, err
	}

	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()

	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)

	req.SetRequestURI(relay)
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json")
	req.SetBody(payload)

	fbHeader, err := fb.flashbotHeader(payload)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-Flashbots-Signature", fbHeader)

	return fb.httpClient.Do(req, res)

}

func (fb *FlashBot) flashbotHeader(payload []byte) (string, error) {

	hashedPayload := crypto.Keccak256Hash(payload).Hex()
	signature, err := crypto.Sign(
		crypto.Keccak256([]byte("\x19Ethereum Signed Message:\n"+strconv.Itoa(len(hashedPayload))+hashedPayload)),
		fb.signerPrivateKey,
	)
	if err != nil {
		return "", err
	}

	return crypto.PubkeyToAddress(fb.signerPrivateKey.PublicKey).Hex() +
		":" + hexutil.Encode(signature), nil
}
