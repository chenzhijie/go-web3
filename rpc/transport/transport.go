package transport

import (
	"strings"
)

type Transport interface {
	Call(method string, out interface{}, params ...interface{}) error
	Close() error
}

type PubSubTransport interface {
	Subscribe(method string, callback func(b []byte)) (func() error, error)
}

const (
	wsPrefix  = "ws://"
	wssPrefix = "wss://"
)

func NewTransport(url, proxy string) (Transport, error) {
	if strings.HasPrefix(url, wsPrefix) || strings.HasPrefix(url, wssPrefix) {
		return newWebsocket(url)
	}
	// if _, err := os.Stat(url); !os.IsNotExist(err) {
	// 	return newIPC(url)
	// }
	return newHTTP(url, proxy), nil
}
