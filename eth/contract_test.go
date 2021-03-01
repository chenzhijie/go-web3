package eth

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestDecodeData(t *testing.T) {
	data := "0x18160ddd"
	ret, err := hex.DecodeString(data)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("ret %x\n", ret)
}
