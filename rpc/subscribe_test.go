package rpc

import (
	"fmt"
	"testing"
	"time"
)

func TestSubsrice(t *testing.T) {
	client, err := NewClient("wss://mainnet.infura.io/ws/v3/", "")
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Subscribe("newPendingTransactions", func(data []byte) {
		fmt.Printf("data %s\n", data)
	})
	if err != nil {
		t.Fatal(err)
	}
	<-time.After(time.Minute)
}
