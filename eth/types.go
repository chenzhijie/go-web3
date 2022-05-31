package eth

import "math/big"

type EstimateFee struct {
	BaseFee              *big.Int
	MaxFeePerGas         *big.Int
	MaxPriorityFeePerGas *big.Int
}
