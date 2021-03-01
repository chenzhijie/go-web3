package utils

import (
	"fmt"
	"testing"
)

func TestToWei(t *testing.T) {
	ethVal := 0.003

	u := NewUtils()

	wei := u.ToWei(ethVal)
	fmt.Printf("wei %v\n", wei)
	fmt.Printf("wei hex %v\n", u.ToHex(wei))
}
