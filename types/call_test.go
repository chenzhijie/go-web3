package types

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"
)

func TestInputData(t *testing.T) {

	data, err := hex.DecodeString("18160ddd")
	if err != nil {
		panic(err)
	}

	call := &CallMsg{
		Data:     data,
		GasPrice: NewCallMsgBigInt(big.NewInt(1000)),
		Value:    NewCallMsgBigInt(big.NewInt(10)),
	}

	ret, err := json.Marshal(call)
	if err != nil {
		panic(err)
	}

	// out := make(InputData, 4)
	// if err := json.Unmarshal(ret, &out); err != nil {
	// 	panic(err)
	// }
	fmt.Printf("ret %s\n", ret)

	// fmt.Printf("out %v\n", out)

}

func BenchmarkTestCallMsgMarshal(b *testing.B) {

	data, err := hex.DecodeString("18160ddd")
	if err != nil {
		panic(err)
	}

	call := &CallMsg{
		Data:     data,
		GasPrice: NewCallMsgBigInt(big.NewInt(1000)),
		Value:    NewCallMsgBigInt(big.NewInt(10)),
	}

	for i := 0; i < b.N; i++ {
		_, err = json.Marshal(call)
		if err != nil {
			panic(err)
		}
	}

}
