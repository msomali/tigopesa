package internal

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	ContentTypeXml  = "application/xml"
	ContentTypeJson = "application/json"
)

var (
	defaultResponseHeader = map[string]string{
		"Content-Type": ContentTypeXml,
	}
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
func Receive(r *http.Request, payloadType PayloadType, v interface{}) error {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if v == nil {
		return fmt.Errorf("v can not be nil")
	}
	// restore request body
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

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
		return xml.Unmarshal(body, v)
	}

	return err
}

//Reply respond to Tigo requests like callback request and namecheck
func Reply(r *Response, writer http.ResponseWriter) {

	if r.Payload == nil {

		for key, value := range r.Headers {
			writer.Header().Add(key, value)
		}
		writer.WriteHeader(r.StatusCode)
		_, _ = writer.Write([]byte("ok"))
		return
	}

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


