package crypto

import (
	"bytes"

	"github.com/chenzhijie/go-web3/utils"
	"github.com/ethereum/go-ethereum/common"
)

func VerifyProof(proofs []common.Hash, root []byte, data []byte) (bool, error) {
	computedHash := make([]byte, len(data))
	copy(computedHash[:], data[:])
	util := utils.NewUtils()
	for _, proof := range proofs {
		var err error
		var abiEncodePacked []byte
		if bytes.Compare(computedHash, proof[:]) <= 0 {
			abiEncodePacked, err = util.AbiEncodePacked(computedHash, proof[:])
		} else {
			abiEncodePacked, err = util.AbiEncodePacked(proof[:], computedHash)
		}
		if err != nil {
			return false, err
		}
		computedHash = Keccak256Hash(abiEncodePacked)

	}
	return bytes.Equal(computedHash, root), nil

}
