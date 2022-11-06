package flashbots

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/chenzhijie/go-web3"
	"github.com/chenzhijie/go-web3/types"
	"github.com/ethereum/go-ethereum/common"
	eTypes "github.com/ethereum/go-ethereum/core/types"
)

const (
	mainnetInfuraProvider     = "https://mainnet.infura.io/v3/91ffab09868d430f9ce744c78d7ff427"
	goerliInfuraProvider      = "https://goerli.infura.io/v3/91ffab09868d430f9ce744c78d7ff427"
	goerliFlashbotMintNFTAddr = "0x20EE855E43A7af19E407E39E5110c2C1Ee41F64D"
)

func TestFlashbotSendBundleTx(t *testing.T) {

	signerKey := os.Getenv("signerKey")

	if len(signerKey) == 0 {
		t.Fatal("signer key or sender key is empty")
	}

	web3, err := web3.NewWeb3(goerliInfuraProvider)
	if err != nil {
		t.Fatal(err)
	}

	err = web3.Eth.SetAccount(signerKey)
	if err != nil {
		t.Fatal(err)
	}

	web3.Eth.SetChainId(5)

	currentBlockNumber, err := web3.Eth.GetBlockNumber()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("currentBlockNumber %v\n", currentBlockNumber)

	mintValue := web3.Utils.ToWei("0.03")

	mintNFTData, err := hex.DecodeString("1249c58b")
	if err != nil {
		t.Fatal(err)
	}

	bundleTxs := make([]*eTypes.Transaction, 0)

	gasLimit, err := web3.Eth.EstimateGas(&types.CallMsg{
		From:  web3.Eth.Address(),
		To:    common.HexToAddress(goerliFlashbotMintNFTAddr),
		Data:  mintNFTData,
		Value: types.NewCallMsgBigInt(mintValue),
	})
	if err != nil {
		t.Fatal(err)
	}
	nonce, err := web3.Eth.GetNonce(web3.Eth.Address(), nil)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("gaslimit %v nonce %v\n", gasLimit, nonce)

	mintNFTtx, err := web3.Eth.NewEIP1559Tx(
		common.HexToAddress(goerliFlashbotMintNFTAddr),
		mintValue, // 6a94d74f430000
		gasLimit,
		web3.Utils.ToGWei(0),  //
		web3.Utils.ToGWei(30), // b2d05e00
		mintNFTData,
		nonce,
	)
	if err != nil {
		t.Fatal(err)
	}

	bundleTxs = append(bundleTxs, mintNFTtx)

	fb, err := NewFlashBot(TestRelayURL, signerKey)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := fb.Simulate(
		bundleTxs,
		big.NewInt(int64(currentBlockNumber)),
		"latest",
	)

	if err != nil {
		t.Fatal(err)
	}
	egp, err := resp.EffectiveGasPrice()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Resp %s EffectiveGasPrice %v\n", resp, web3.Utils.FromGWei(egp))
	bundleResp, err := fb.SendBundle(
		bundleTxs,
		big.NewInt(int64(currentBlockNumber)+1),
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("bundle resp %v\n", bundleResp)

}

func TestGetBundleStats(t *testing.T) {

	signerKey := os.Getenv("signerKey")

	if len(signerKey) == 0 {
		t.Fatal("signer key or sender key is empty")
	}

	fb, err := NewFlashBot(DefaultRelayURL, signerKey)
	if err != nil {
		t.Fatal(err)
	}

	bundleHash := "0x4a11aa0e0bdc321a7bbe5c96f9952cc38e38d8843b293379761d736222f8635b"
	targetBlockNumber := big.NewInt(int64(6974433))
	stat, err := fb.GetBunderStats(bundleHash, targetBlockNumber)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("bundle stat %v\n", stat)

}

func TestGetUserStats(t *testing.T) {

	signerKey := os.Getenv("signerKey")

	if len(signerKey) == 0 {
		t.Fatal("signer key or sender key is empty")
	}

	fb, err := NewFlashBot(TestRelayURL, signerKey)
	if err != nil {
		t.Fatal(err)
	}

	targetBlockNumber := big.NewInt(int64(6974433))
	fmt.Printf("%x\n", targetBlockNumber.Int64())
	stat, err := fb.GetUserStats(targetBlockNumber)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("user stat %v\n", stat)

}
