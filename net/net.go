package net

import (
	"github.com/chenzhijie/go-web3/rpc"
	"github.com/chenzhijie/go-web3/utils"
)

type Net struct {
	c *rpc.Client
}

func (n *Net) Version() (uint64, error) {
	var out string
	if err := n.c.Call("net_version", &out); err != nil {
		return 0, err
	}
	return utils.ParseUint64orHex(out)
}

func (n *Net) Listening() (bool, error) {
	var out bool
	err := n.c.Call("net_listening", &out)
	return out, err
}

func (n *Net) PeerCount() (uint64, error) {
	var out string
	if err := n.c.Call("net_peerCount", &out); err != nil {
		return 0, err
	}
	return utils.ParseUint64orHex(out)
}
