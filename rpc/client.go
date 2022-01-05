package rpc

import (
	"github.com/chenzhijie/go-web3/rpc/transport"
)

type Client struct {
	transport transport.Transport
	addr      string
}

func NewClient(addr, proxy string) (*Client, error) {
	c := &Client{
		addr: addr,
	}

	t, err := transport.NewTransport(addr, proxy)
	if err != nil {
		return nil, err
	}
	c.transport = t
	return c, nil
}

func (c *Client) Close() error {
	return c.transport.Close()
}

func (c *Client) Call(method string, out interface{}, params ...interface{}) error {
	return c.transport.Call(method, out, params...)
}
