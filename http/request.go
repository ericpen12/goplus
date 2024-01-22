package http

import (
	"bytes"
	"encoding/json"
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

func WithHost(host string) Option {
	return func(r *httpRequest) {
		r.url.Host = host
	}
}

func WithUrl(URL string) Option {
	return func(r *httpRequest) {
		r.url, _ = url.Parse(URL)
	}
}

func WithUri(uri string) Option {
	return func(r *httpRequest) {
		r.url = r.url.JoinPath(uri)
	}
}

func WithUrlParams(urlParams map[string]string) Option {
	return func(r *httpRequest) {
		queryParams := url.Values{}
		for k, v := range urlParams {
			queryParams.Add(k, v)
		}
		r.url.RawQuery = queryParams.Encode()
	}
}

func WithJson(data map[string]interface{}) Option {
	return func(r *httpRequest) {
		r.header.Set("Content-Type", "application/json;charset=UTF-8")
		b, _ := json.Marshal(data)
		r.body = bytes.NewReader(b)
	}
}

func WithMethodPost() Option {
	return func(r *httpRequest) {
		r.method = http.MethodPost
	}
}

func WithResponse(res any) Option {
	return func(r *httpRequest) {
		r.response = res
	}
}

func WithHeader(headers map[string]string) Option {
	return func(r *httpRequest) {
		for k, v := range headers {
			r.header.Set(k, v)
		}
	}
}

func WithToken(token string) Option {
	return func(r *httpRequest) {
		r.header.Set("Authorization", token)
	}
}
