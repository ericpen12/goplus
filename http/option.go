package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
)

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
