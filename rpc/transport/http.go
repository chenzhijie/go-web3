package transport

import (
	"encoding/json"
	"net"
	"time"

	"github.com/chenzhijie/go-web3/rpc/codec"
	"github.com/valyala/fasthttp"
)

type HTTP struct {
	addr   string
	client *fasthttp.Client
}

func newHTTP(addr string) *HTTP {
	return &HTTP{
		addr: addr,
		client: &fasthttp.Client{
			Dial: func(addr string) (net.Conn, error) {
				return fasthttp.DialTimeout(addr, time.Duration(10)*time.Second)
			},
		},
	}
}

func (h *HTTP) Close() error {
	return nil
}

func (h *HTTP) Call(method string, out interface{}, params ...interface{}) error {
	request := codec.Request{
		Method: method,
	}
	if len(params) > 0 {
		data, err := json.Marshal(params)
		if err != nil {
			return err
		}
		request.Params = data
	}
	raw, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()

	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)

	req.SetRequestURI(h.addr)
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json")
	req.SetBody(raw)

	// fmt.Printf("req body %s\n", raw)

	if err := h.client.Do(req, res); err != nil {
		return err
	}

	var response codec.Response
	if err := json.Unmarshal(res.Body(), &response); err != nil {
		return err
	}
	if response.Error != nil {
		return response.Error
	}

	if err := json.Unmarshal(response.Result, out); err != nil {
		return err
	}
	return nil
}
