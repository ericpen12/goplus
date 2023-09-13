package http

import (
	"io"
	"net/http"
	"net/url"
)

type httpRequest struct {
	url      *url.URL
	client   *http.Client
	response any
	method   string
	body     io.Reader
	header   http.Header
}

type Option func(*httpRequest)

type Request interface {
	Do() (EasyResponse, error)
	URL() *url.URL
}

func EasyRequest(opt ...Option) Request {
	r := newHttpRequest()
	for _, o := range opt {
		o(r)
	}
	return r
}

func newHttpRequest() *httpRequest {
	return &httpRequest{
		client: &http.Client{},
		url:    &url.URL{},
		header: map[string][]string{},
		method: "POST",
	}
}

func (h *httpRequest) URL() *url.URL {
	return h.url
}
func (h *httpRequest) Do() (EasyResponse, error) {
	req, err := http.NewRequest(h.method, h.url.String(), h.body)
	if err != nil {
		return nil, err
	}
	req.Header = h.header
	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return Read(resp.Body, h.response)
}
