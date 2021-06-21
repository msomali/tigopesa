package tigo

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/techcraftt/tigosdk/internal"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
	"time"
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
	Config struct {
		PayAccountName   string
		PayAccountMSISDN string
		PayBillerNumber  string
		PayRequestURL    string
		PayNamecheckURL  string

		DisburseAccountName   string
		DisburseAccountMSISDN string
		DisburseBrandID       string
		DisbursePIN           string
		DisburseRequestURL    string

		PushUsername              string
		PushPassword              string
		PushPasswordGrantType     string
		PushApiBaseURL            string
		PushGetTokenURL           string
		PushBillerMSISDN          string
		PushBillerCode            string
		PushPushPayURL            string
		PushReverseTransactionURL string
		PushHealthCheckURL        string
	}



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
		//Config:     config,
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

	errs := make(chan error)
	done := make(chan bool)

	go func() {
		buf, err := internal.MarshalPayload(t, payload)
		if err != nil {
			errs <- fmt.Errorf("could not marshal the payload: %s\n", err.Error())
		}
		_, err = client.Logger.Write([]byte(fmt.Sprintf("%s: %s\n\n", prefix, string(buf))))
		if err != nil {
			errs <- fmt.Errorf("could not print the payload: %s\n", err.Error())
		}

		done <- true
	}()

	select {
	case err := <-errs:
		fmt.Printf("error encountered: %s\n", err)
		return

	case <-done:
		fmt.Printf("done logging\n")
		return
	}
}

func (client *BaseClient) Log(request *http.Request, response *http.Response) {

	errs := make(chan error)
	done := make(chan bool)

	go func() {
		if request != nil {
			reqDump, err := httputil.DumpRequestOut(request, true)
			if err != nil {
				errs <- fmt.Errorf("could not dump request due to: %s\n", err.Error())
			}
			_, err = fmt.Fprintf(client.Logger, "request %s\n", reqDump)
			if err != nil {
				errs <- fmt.Errorf("could not print  request due to: %s\n", err.Error())
			}

			if response == nil {
				done <- true
			}
		}

		if response != nil {
			respDump, err := httputil.DumpResponse(response, true)
			if err != nil {
				errs <- fmt.Errorf("could not dump response due to: %s\n", err.Error())
			}
			_, err = fmt.Fprintf(client.Logger, "response:  %s\n", respDump)

			if err != nil {
				errs <- fmt.Errorf("could not print out response due to: %s\n", err.Error())
			}

			done <- true

		}
	}()

	select {
	case err := <-errs:
		fmt.Printf("error encountered: %s\n", err)
		return

	case <-done:
		fmt.Printf("done logging\n")
		return
	}

}


func (client *BaseClient) Send(_ context.Context, request *Request, v interface{}) error {
	var bodyBytes []byte

	//creates http request with context
	req, err := request.Transform()

	if err != nil{
		return err
	}

	if req.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(req.Body)
	}

	// Restore the request body content.
	req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	if v == nil {
		return errors.New("v interface can not be empty")
	}

	//// sending json request by default
	//if req.Header.Get("content-Type") == "" {
	//	req.Header.Set("Content-Type", "application/json")
	//	req.Header.Set("Cache-Control", "no-cache")
	//	req.Header.Set("username", client.Username)
	//	req.Header.Set("password", c.Password)
	//}

	resp, err := client.HttpClient.Do(req)

	go func(debugMode bool) {
		if debugMode{
			client.Log(req,resp)
		}
	}(client.DebugMode)

	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	switch resp.Header.Get("Content-Type") {
	case "application/json", "application/json;charset=UTF-8":
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			if err != io.EOF {
				return err
			}
		}
	case "text/xml", "text/xml;charset=UTF-8":
		if err := xml.NewDecoder(resp.Body).Decode(v); err != nil {
			if err != io.EOF {
				return err
			}
		}
	}

	return nil
}
