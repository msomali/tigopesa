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
	"strings"
	"time"

	"github.com/techcraftt/tigosdk/internal"
)

const (
	defaultTimeout          = time.Minute
	jsonContentTypeString   = "application/json"
	xmlContentTypeString    = "text/xml"
	appXMLContentTypeString = "application/xml"
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

// Log is called to print the details of http.Request sent from Tigo during
// callback, namecheck or ussd payment. It is used for debugging purposes
func (client *BaseClient) Log(name string, request *http.Request) {

	if request != nil {
		reqDump, _ := httputil.DumpRequest(request, true)
		_, err := fmt.Fprintf(client.Logger, "%s REQUEST: %s\n", name, reqDump)
		if err != nil {
			fmt.Printf("error while logging %s request: %v\n",
				strings.ToLower(name), err)
			return
		}

		return
	}

	return

}

// LogOut is like Log except this is for outgoing client requests:
// http.Request that is supposed to be sent to tigo
func (client *BaseClient) LogOut(name string, request *http.Request, response *http.Response) {

	if request != nil {
		reqDump, _ := httputil.DumpRequestOut(request, true)
		_, err := fmt.Fprintf(client.Logger, "%s REQUEST: %s\n", name, reqDump)
		if err != nil {
			fmt.Printf("error while logging %s request: %v\n",
				strings.ToLower(name), err)
		}
	}

	if response != nil {
		respDump, _ := httputil.DumpResponse(response, true)
		_, err := fmt.Fprintf(client.Logger, "%s RESPONSE: %s\n", name, respDump)
		if err != nil {
			fmt.Printf("error while logging %s response: %v\n",
				strings.ToLower(name), err)
		}
	}

}

func (client *BaseClient) Send(_ context.Context, rn RequestName, request *Request, v interface{}) error {
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

	name := strings.ToUpper(rn.String())

	go func(debugMode bool) {
		if debugMode {
			client.LogOut(name, req, resp)
		}
	}(client.DebugMode)

	// restore response http.Response.Body after debugging
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(respBodyBytes))

	contentType := resp.Header.Get("Content-Type")

	if strings.Contains(contentType, jsonContentTypeString) {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			if err != io.EOF {
				return err
			}
		}
	}

	if strings.Contains(contentType, xmlContentTypeString) ||
		strings.Contains(contentType, appXMLContentTypeString) {
		if err := xml.NewDecoder(resp.Body).Decode(v); err != nil {
			if err != io.EOF {
				return err
			}
		}
	}
	return nil
}
