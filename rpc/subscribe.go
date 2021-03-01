package rpc

import (
	"fmt"

	"github.com/chenzhijie/go-web3/rpc/transport"
)

func (c *Client) SubscriptionEnabled() bool {
	_, ok := c.transport.(transport.PubSubTransport)
	return ok
}

func (c *Client) Subscribe(method string, callback func(b []byte)) (func() error, error) {
	pub, ok := c.transport.(transport.PubSubTransport)
	if !ok {
		return nil, fmt.Errorf("Transport does not support the subscribe method")
	}
	close, err := pub.Subscribe(method, callback)
	return close, err
}
