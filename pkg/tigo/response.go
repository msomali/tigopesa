package tigo

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/techcraftt/tigosdk/internal"
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
	Response struct {
		StatusCode  int
		Payload     interface{}
		PayloadType internal.PayloadType
		Headers     map[string]string
	}

	ResponseOption func(response *Response)
)

//Receive takes *http.Request from Tigo like during push pay callback and name query
//It then unmarshal the provided request into given interface v
//The expected Content-Type should also be declared. If its application/json or
//application/xml
func Receive(r *http.Request, payloadType internal.PayloadType, v interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if v == nil {
		return fmt.Errorf("v can not be nil")
	}

	switch payloadType {
	case internal.JsonPayload:
		err = json.NewDecoder(r.Body).Decode(v)
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(r.Body)

		return err

	case internal.XmlPayload:
		return xml.Unmarshal(body, v)
	}

	return err
}

//Reply respond to Tigo requests like callback request and namecheck
func (r *Response) Reply(writer http.ResponseWriter) (err error) {

	var payload []byte
	errs := make(chan error)
	done := make(chan bool)

	switch r.PayloadType {
	case internal.XmlPayload:
		payload, err = xml.MarshalIndent(r.Payload, "", "  ")
		if err != nil {
			errs <- err
		}
		done <- true

	case internal.JsonPayload:
		payload, err = json.Marshal(r.Payload)
		if err != nil {
			errs <- err
		}
		<- done
	}

	select {
	case err := <-errs:
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return err

	case <-done:
		writer.WriteHeader(r.StatusCode)
		for key, value := range r.Headers {
			writer.Header().Add(key, value)
		}
		_,err = writer.Write(payload)
		return err
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
		headers := map[string]string{
			"Content-Type": ContentTypeJson,
		}
		response.Headers = headers
	}
}

func NewResponse(status int, payload interface{}, payloadType internal.PayloadType, opts ...ResponseOption) *Response {
	response := &Response{
		StatusCode:  status,
		Payload:     payload,
		PayloadType: payloadType,
		Headers:     defaultResponseHeader,
	}

	for _, opt := range opts {
		opt(response)
	}

	return response
}
