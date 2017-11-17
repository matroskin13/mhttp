package mhttp

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Request struct {
	URI       *url.URL
	Method    string
	Type      string
	Headers   http.Header
	cacheBody []byte
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

func NewRequest(path string, method string, headers map[string]string) (*Request, error) {
	uri, err := url.Parse(path)

	if err != nil {
		return nil, err
	}

	request := &Request{URI: uri, Method: method, Headers: make(http.Header)}

	for key, header := range headers {
		request.Headers.Add(key, header)
	}

	return request, nil
}

func (r *Request) Do(body []byte) (*Response, error) {
	r.URI.String()
	request, err := http.NewRequest(r.Method, r.URI.String(), bytes.NewBuffer(body))

	request.Header = r.Headers

	if err != nil {
		return nil, err
	}

	client := http.Client{}

	r.cacheBody = body

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

func (r Request) GetPrettyRequest() (string, error) {
	path := r.URI.Path

	if path == "" {
		path = "/"
	}

	prettyString := fmt.Sprintf("%s %s HTTP/1.1 \r\n", r.Method, path)

	for key := range r.Headers {
		prettyString = fmt.Sprintf(
			"%s%s: %s\r\n",
			prettyString, key,
			r.Headers.Get("Content-Type"),
		)
	}

	if len(r.cacheBody) > 0 {
		prettyString = fmt.Sprintf("%s\r\n%s", prettyString, r.cacheBody)
	}

	return prettyString, nil
}
