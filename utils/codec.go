package utils

import (
	"encoding/hex"
	"fmt"
	"strings"
)

func unmarshalTextByte(dst, src []byte, size int) error {
	str := string(src)

	str = strings.Trim(str, "\"")
	if !strings.HasPrefix(str, "0x") {
		return fmt.Errorf("0x prefix not found")
	}
	str = str[2:]
	b, err := hex.DecodeString(str)
	if err != nil {
		return err
	}
	if len(b) != size {
		return fmt.Errorf("length %d is not correct, expected %d", len(b), size)
	}
	copy(dst, b)
	return nil
}
