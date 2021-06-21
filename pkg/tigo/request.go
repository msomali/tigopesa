package tigo

import (
	"bytes"
	"context"
	"github.com/techcraftt/tigosdk/internal"
	"io"
	"net/http"
)

var (
	defaultRequestHeaders = map[string]string{
		"Content-Type":  "application/json",
		"Cache-Control": "no-cache",
	}
)

type (
	Request struct {
		Context     context.Context
		HttpMethod  string
		URL         string
		PayloadType internal.PayloadType
		Payload     interface{}
		Headers     map[string]string
	}

	RequestOption func(request *Request)
)

func NewRequest(method, url string, payloadType internal.PayloadType, payload interface{}, opts ...RequestOption) *Request {
	request := &Request{
		Context:     context.TODO(),
		HttpMethod:  method,
		URL:         url,
		PayloadType: payloadType,
		Payload:     payload,
		Headers:     defaultRequestHeaders,
	}

	for _, opt := range opts {
		opt(request)
	}

	return request
}

func WithRequestContext(ctx context.Context) RequestOption {
	return func(request *Request) {
		request.Context = ctx
	}
}

func WithRequestHeaders(headers map[string]string) RequestOption {
	return func(request *Request) {
		request.Headers = headers
	}
}

func WithMoreHeaders(headers map[string]string) RequestOption {
	return func(request *Request) {
		for key, value := range headers{
			request.Headers[key] = value
		}
	}
}

func WithAuthHeaders(username, password string) RequestOption {
	return func(request *Request) {
		request.Headers["username"] = username
		request.Headers["password"] = password
	}
}


//Transform takes a *Request and transform into *http.Request with a context
func (request *Request) Transform() (*http.Request, error) {
	var (
		buffer io.Reader
	)

	if request.Payload != nil {
		buf, err := internal.MarshalPayload(request.PayloadType, request.Payload)
		if err != nil {
			return nil, err
		}

		buffer = bytes.NewBuffer(buf)
	}

	req, err := http.NewRequestWithContext(request.Context, request.HttpMethod, request.URL, buffer)
	if err != nil{
		return nil, err
	}
	for key, value := range request.Headers{
		req.Header.Set(key,value)
	}
	return req, nil
}
