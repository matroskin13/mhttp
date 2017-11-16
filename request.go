package mhttp

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Request struct {
	URI     *url.URL
	Method  string
	Type    string
	Headers http.Header
}

type Response struct {
	http.Response

	BodyRaw []byte
}

const (
	TypeHTML = "html"
	TypeJSON = "json"
)

func GetTypeByAlias(alias string) string {
	switch alias {
	case TypeJSON:
		return "application/json"
	case TypeHTML:
		return "text/html"
	}

	return alias
}

func NewRequest(path string, method string, httpType string) (*Request, error) {
	uri, err := url.Parse(path)

	if err != nil {
		return nil, err
	}

	request := &Request{URI: uri, Method: method, Headers: make(http.Header)}

	if httpType != "" {
		request.Headers.Add("Content-Type", httpType)
	}

	return request, nil
}

func (r Request) Do(body []byte) (*Response, error) {
	r.URI.String()
	request, err := http.NewRequest(r.Method, r.URI.String(), bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	client := http.Client{}

	resp, err := client.Do(request)

	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return &Response{
		BodyRaw:  b,
		Response: *resp,
	}, nil
}
