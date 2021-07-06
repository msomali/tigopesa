package tigo

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
	"time"

	"github.com/techcraftt/tigosdk/internal"
)

const (
	defaultTimeout = time.Minute
)

var (
	// defaultCtx is the context used by pkg when none is set
	// to override this one has to call WithContext method and supply
	// his her own context.Context
	defaultCtx = context.TODO()

	// defaultWriter is an io.Writer used for debugging. When debug mode is
	// set to true i.e DEBUG=true and no io.Writer is provided via
	// WithLogger method this is used.
	defaultWriter = os.Stderr

	// defaultHttpClient is the pkg used by library to send Http requests, specifically
	// disbursement requests in case a user does not specify one
	defaultHttpClient = &http.Client{
		Timeout: defaultTimeout,
	}
)

type (
	BaseClient struct {
		HttpClient *http.Client
		Ctx        context.Context
		Timeout    time.Duration
		Logger     io.Writer // for logging purposes
		DebugMode  bool
	}
)

func NewBaseClient(opts ...ClientOption) *BaseClient {
	client := &BaseClient{
		//Config:     conf,
		HttpClient: defaultHttpClient,
		Logger:     defaultWriter,
		Timeout:    defaultTimeout,
		Ctx:        defaultCtx,
		DebugMode:  false,
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

func (client *BaseClient) NewRequest(method, url string, payloadType internal.PayloadType, payload interface{}) (*http.Request, error) {
	var (
		ctx, _ = context.WithTimeout(context.Background(), client.Timeout)
	)

	request := &Request{
		Context:     ctx,
		HttpMethod:  method,
		URL:         url,
		PayloadType: payloadType,
		Payload:     payload,
	}
	return request.Transform()
}

func (client *BaseClient) LogPayload(t internal.PayloadType, prefix string, payload interface{}) {
	buf, _ := internal.MarshalPayload(t, payload)
	_, _ = client.Logger.Write([]byte(fmt.Sprintf("%s: %s\n\n", prefix, buf.String())))
}

func (client *BaseClient) Log(request *http.Request, response *http.Response) {

	if request != nil {
		reqDump, _:= httputil.DumpRequestOut(request, true)
		_, _ = fmt.Fprintf(client.Logger, "REQUEST: %s\n", reqDump)
	}

	if response != nil {
		respDump, _ := httputil.DumpResponse(response, true)
		_, _ = fmt.Fprintf(client.Logger, "RESPONSE: %s\n", respDump)
	}

}

func (client *BaseClient) Send(_ context.Context, request *Request, v interface{}) error {
	var bodyBytes []byte

	//creates http request with context
	req, err := request.Transform()

	if err != nil {
		return err
	}

	if req.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(req.Body)
	}

	if v == nil {
		return errors.New("v interface can not be empty")
	}

	// restore request body for logging
	req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	resp, err := client.HttpClient.Do(req)

	if err != nil {
		return err
	}

	//restore request body for logging
	req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	var respBodyBytes []byte
	if resp.Body != nil {
		respBodyBytes, _ = ioutil.ReadAll(resp.Body)
	}

	go func(debugMode bool) {
		if debugMode {
			client.Log(req, resp)
		}
	}(client.DebugMode)

	// restore response http.Response.Body after debugging
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(respBodyBytes))

	switch resp.Header.Get("Content-Type") {
	case "application/json", "application/json;charset=UTF-8":
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			if err != io.EOF {
				return err
			}
		}
	case "application/xml", "text/xml;charset=utf-8", "text/xml", "text/xml;charset=UTF-8":
		if err := xml.NewDecoder(resp.Body).Decode(v); err != nil {
			if err != io.EOF {
				return err
			}
		}
	}

	return nil
}
