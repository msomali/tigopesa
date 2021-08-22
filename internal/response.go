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
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	ContentTypeXml  = "application/xml"
	ContentTypeJson = "application/json"
)

type (
	// Response contains details to be sent as a response to tigo
	// when tigo make callback request, name check request or payment
	// request.
	Response struct {
		StatusCode  int
		Payload     interface{}
		PayloadType PayloadType
		Headers     map[string]string
		Error       error
	}

	ResponseOption func(response *Response)
)

//NewResponse create a response to be sent back to Tigo. HTTP Status code, payload and its
// type need to be specified. Other fields like  Response.Error and Response.Headers can be
// changed using WithMoreResponseHeaders (add headers), WithResponseHeaders (replace all the
// existing ) and WithResponseError to add error its default value is nil, default value of
// Response.Headers is
// defaultResponseHeader = map[string]string{
//		"Content-Type": ContentTypeXml,
// }
func NewResponse(status int, payload interface{}, payloadType PayloadType, opts ...ResponseOption) *Response {

	var (
		defaultResponseHeader = map[string]string{
			"Content-Type": ContentTypeXml,
		}
	)
	response := &Response{
		StatusCode:  status,
		Payload:     payload,
		PayloadType: payloadType,
		Headers:     defaultResponseHeader,
		Error:       nil,
	}

	for _, opt := range opts {
		opt(response)
	}

	return response
}

//Receive takes *http.Request from Tigo like during push pay callback and name query
//It then unmarshal the provided request into given interface v
//The expected Content-Type should also be declared. If its application/json or
//application/xml
func (client *BaseClient) Receive(r *http.Request, rn RequestName, payloadType PayloadType, v interface{}) error {

	var reqBodyBytes []byte
	reqBodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	defer func(debug bool) {
		if debug {
			r.Body = io.NopCloser(bytes.NewBuffer(reqBodyBytes))
			name := strings.ToUpper(rn.String())
			client.log(name, r)
		}
	}(client.DebugMode)

	if v == nil {
		return fmt.Errorf("v can not be nil")
	}
	// restore request body
	r.Body = ioutil.NopCloser(bytes.NewBuffer(reqBodyBytes))

	switch payloadType {
	case JsonPayload:
		err := json.NewDecoder(r.Body).Decode(v)
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(r.Body)

		return err

	case XmlPayload:
		return xml.Unmarshal(reqBodyBytes, v)
	}

	return err
}

//Reply respond to Tigo requests like callback request and namecheck
func (client *BaseClient) Reply(name string, r *Response, writer http.ResponseWriter) {

	if r.Payload == nil {

		for key, value := range r.Headers {
			writer.Header().Add(key, value)
		}
		writer.WriteHeader(r.StatusCode)
		_, _ = writer.Write([]byte("ok"))
		return
	}

	defer func(debug bool) {
		client.logPayload(r.PayloadType, strings.ToUpper(name), r.Payload)
		return
	}(client.DebugMode)

	if r.Error != nil {
		http.Error(writer, r.Error.Error(), http.StatusInternalServerError)
		return
	}

	switch r.PayloadType {
	case XmlPayload:
		payload, err := xml.MarshalIndent(r.Payload, "", "  ")
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		for key, value := range r.Headers {
			writer.Header().Set(key, value)
		}

		writer.WriteHeader(r.StatusCode)
		_, err = writer.Write(payload)
		return

	case JsonPayload:
		payload, err := json.Marshal(r.Payload)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		for key, value := range r.Headers {
			writer.Header().Set(key, value)
		}
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(r.StatusCode)
		_, err = writer.Write(payload)
		return
	}

}

func WithResponseHeaders(headers map[string]string) ResponseOption {
	return func(response *Response) {
		response.Headers = headers
	}
}

func WithMoreResponseHeaders(headers map[string]string) ResponseOption {
	return func(response *Response) {
		for key, value := range headers {
			response.Headers[key] = value
		}
	}
}

func WithDefaultJsonHeader() ResponseOption {
	return func(response *Response) {
		response.Headers["Content-Type"] = ContentTypeJson
	}
}

func WithResponseError(err error) ResponseOption {
	return func(response *Response) {
		response.Error = err
	}
}
