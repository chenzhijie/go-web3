package eth

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

var privateKeyUsedForTest = "1b734ae16eb3b7470d99780dff19bc7e2d8ce5b04785a7390d7363e78d37c6e8"

func TestSignText(t *testing.T) {
	eth := NewEth(nil)
	if err := eth.SetAccount(privateKeyUsedForTest); err != nil {
		t.Fatal(err)
	}

	signature, err := eth.SignText([]byte("hello"))
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("signature %x\n", signature)
}

func TestSignTypedData(t *testing.T) {

	testTypedData := apitypes.TypedData{
		Domain: apitypes.TypedDataDomain{
			Name:    "Hashflow - Identity Verification",
			Version: "1.0",
		},
		Types: apitypes.Types{
			"EIP712Domain": {
				{
					Name: "name",
					Type: "string",
				},
				{
					Name: "version",
					Type: "string",
				},
			},
			"Identity": {
				{
					Name: "wallet",
					Type: "address",
				},
				{
					Name: "timestamp",
					Type: "uint256",
				},
			},
		},
		Message: apitypes.TypedDataMessage{
			"wallet":    "0xb1c0d8c7ca1a5bb05c57b99bb5acdc498062b060",
			"timestamp": "1667274062",
		},
		PrimaryType: "Identity",
	}

	eth := NewEth(nil)
	if err := eth.SetAccount(privateKeyUsedForTest); err != nil {
		t.Fatal(err)
	}

	signature, err := eth.SignTypedData(testTypedData)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("signature %x\n", signature)
}
