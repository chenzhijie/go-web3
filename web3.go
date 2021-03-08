package web3

import (
	"github.com/chenzhijie/go-web3/eth"
	"github.com/chenzhijie/go-web3/rpc"
	"github.com/chenzhijie/go-web3/utils"
)

type Web3 struct {
	Eth   *eth.Eth
	Utils *utils.Utils
	c     *rpc.Client
}

func NewWeb3(provider string) (*Web3, error) {
	c, err := rpc.NewClient(provider)
	if err != nil {
		return nil, err
	}
	e := eth.NewEth(c)
	e.SetChainId(1)
	u := utils.NewUtils()
	w := &Web3{
		Eth:   e,
		Utils: u,
		c:     c,
	}
	return w, nil
}

func (w *Web3) Version() (string, error) {
	var out string
	err := w.c.Call("web3_clientVersion", &out)
	return out, err
}
