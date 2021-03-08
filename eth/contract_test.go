package eth

import (
	"encoding/hex"
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

func TestDecodeData(t *testing.T) {
	data := "0x18160ddd"
	ret, err := hex.DecodeString(data)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("ret %x\n", ret)
}

func TestForWaitCh(t *testing.T) {

	ch := make(chan bool, 0)

	var isTimeout int32

	go func() {
		time.Sleep(time.Duration(3) * time.Second)
		if atomic.LoadInt32(&isTimeout) == 1 {
			return
		}
		ch <- false
	}()

	select {
	case ret := <-ch:
		if !ret {
			fmt.Println("continue")
		}
	case <-time.After(time.Duration(10) * time.Second):
		atomic.StoreInt32(&isTimeout, 1)
		fmt.Println("timeout")
	}
	fmt.Println("done")
}
