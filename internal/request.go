/*
 * MIT License
 *
 * Copyright (c) 2021 TechCraft Technologies Co. Ltd
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 *
 */

package internal

import (
	"context"
	"net/http"
)

const (
	PushPayRequest RequestName = iota
	DisburseRequest
	GetTokenRequest
	RefundRequest
	HealthCheckRequest
	CallbackRequest
	NameQueryRequest
	PaymentRequest
)

func (rn RequestName) String() string {
	states := [...]string{
		"PushPayRequest",
		"DISBURSE_REQUEST",
		"GetTokenRequest",
		"RefundRequest",
		"Healthcheck",
		"Callbackrequest",
		"NameQueryRequest",
		"PaymentRequest",
	}
	if len(states) < int(rn) {
		return ""
	}

	return states[rn]
}

type (
	//RequestName is used to identify the type of request being saved
	// important in debugging or switch cases where a number of different
	// requests can be served.
	RequestName int

	// Request encapsulate details of a request to be sent to Tigo.
	Request struct {
		Name        RequestName
		Context     context.Context
		Method      string
		URL         string
		PayloadType PayloadType
		Payload     interface{}
		Headers     map[string]string
		Params      map[string]string
	}

	RequestOption func(request *Request)
)

func NewRequest(ctx context.Context, method, url string, payloadType PayloadType, payload interface{}, opts ...RequestOption) *Request {
	var (
		defaultRequestHeaders = map[string]string{
			"Content-Type":  "application/json",
			"Cache-Control": "no-cache",
		}
	)

	request := &Request{
		Context:     ctx,
		Method:      method,
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
		for key, value := range headers {
			request.Headers[key] = value
		}
	}
}

//WithAuthHeaders add password and username to request headers
func WithAuthHeaders(username, password string) RequestOption {
	return func(request *Request) {
		request.Headers["Username"] = username
		request.Headers["Password"] = password
	}
}

func (request *Request) AddHeader(key, value string) {
	request.Headers[key] = value
}

//newRequestWithContext takes a *Request and transform into *http.Request with a context
func (request *Request) newRequestWithContext() (*http.Request, error) {

	buffer, err := MarshalPayload(request.PayloadType, request.Payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(request.Context, request.Method, request.URL, buffer)
	if err != nil {
		return nil, err
	}
	for key, value := range request.Headers {
		req.Header.Add(key, value)
	}

	return req, nil
}
