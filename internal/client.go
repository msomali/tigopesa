package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"
)

const (
	defaultTimeout          = 60 * time.Second
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
		Http      *http.Client
		Ctx       context.Context
		Timeout   time.Duration
		Logger    io.Writer // for logging purposes
		DebugMode bool
	}
)

func NewBaseClient(opts ...ClientOption) *BaseClient {
	client := &BaseClient{
		//Config:     conf,
		Http:      defaultHttpClient,
		Logger:    defaultWriter,
		Timeout:   defaultTimeout,
		Ctx:       defaultCtx,
		DebugMode: false,
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

//func (client *BaseClient) NewRequest(method, url string, payloadType PayloadType, payload interface{}) (*http.Request, error) {
//	var (
//		ctx, _ = context.WithTimeout(context.Background(), client.Timeout)
//	)
//
//	request := &Request{
//		Context:     ctx,
//		Method:      method,
//		URL:         url,
//		PayloadType: payloadType,
//		Payload:     payload,
//	}
//	return request.newRequestWithContext()
//}

func (client *BaseClient) LogPayload(t PayloadType, prefix string, payload interface{}) {
	buf, _ := MarshalPayload(t, payload)
	_, _ = client.Logger.Write([]byte(fmt.Sprintf("%s: %s\n\n", prefix, buf.String())))
}

// log is called to print the details of http.Request sent from Tigo during
// callback, namecheck or ussd payment. It is used for debugging purposes
func (client *BaseClient) log(name string, request *http.Request) {

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

// logOut is like log except this is for outgoing client requests:
// http.Request that is supposed to be sent to tigo
func (client *BaseClient) logOut(name string, request *http.Request, response *http.Response) {

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

	return
}

func (client *BaseClient) Send(ctx context.Context, rn RequestName, request *Request, v interface{}) error {
	var (
		_, cancel = context.WithTimeout(context.Background(), client.Timeout)
	)

	defer cancel()
	var req *http.Request
	var res *http.Response

	var reqBodyBytes []byte
	var resBodyBytes []byte
	defer func(debug bool) {
		if debug {
			req.Body = io.NopCloser(bytes.NewBuffer(reqBodyBytes))
			res.Body = io.NopCloser(bytes.NewBuffer(resBodyBytes))
			name := strings.ToUpper(rn.String())
			client.logOut(name, req, res)
		}
	}(client.DebugMode)

	//creates http request with context
	req, err := request.newRequestWithContext()

	if err != nil {
		return err
	}

	if req.Body != nil {
		reqBodyBytes, _ = io.ReadAll(req.Body)
	}

	if v == nil {
		return errors.New("v interface can not be empty")
	}

	req.Body = io.NopCloser(bytes.NewBuffer(reqBodyBytes))
	res, err = client.Http.Do(req)

	if err != nil {
		return err
	}

	if res.Body != nil {
		resBodyBytes, _ = io.ReadAll(res.Body)
	}

	contentType := res.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		if err := json.NewDecoder(bytes.NewBuffer(resBodyBytes)).Decode(v); err != nil {
			if err != io.EOF {
				return err
			}
		}
	}

	if strings.Contains(contentType, xmlContentTypeString) ||
		strings.Contains(contentType, appXMLContentTypeString) {
		if err := xml.NewDecoder(bytes.NewBuffer(resBodyBytes)).Decode(v); err != nil {
			if err != io.EOF {
				return err
			}
		}
	}
	return nil
}
