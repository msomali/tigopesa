package ussd

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
	"time"
)

type (
	// RequestType represent a request Type that needs to be handled
	// At the moment there are only two types of requests that are expected
	// one is the name-check request and the other is the W2A request
	// this is an experimental API still not implemented
	// Expected implementation, though it may change may involve something like
	// func (client *Client) HandleRequest(ctx context.Context, requestType RequestType) http.HandlerFunc {
	//   return func(writer http.ResponseWriter, request *http.Request) {
	//      switch requestType {
	//      case SubscriberName:
	//      // implementation logic
	//
	//      case WalletToAccount
	//      //implementation logic
	//
	//      default:
	//         http.Error(writer, "unknown request type", http.StatusInternalServerError)
	//      }
	//   }
	// }
	RequestType int

	QuerySubscriberFunc func(ctx context.Context, request SubscriberNameRequest) (SubscriberNameResponse, error)
	WalletToAccountFunc func(ctx context.Context, request WalletToAccountRequest) (WalletToAccountResponse, error)

	// Client is a ussd client that assembles together necessary parts
	// needed in performing 3 operations
	// 1. handling of NameCheckRequest, 2. Handling of Collection (WalletToAccount)
	// requests and 3. Making disbursement requests
	// it contains Config to store all needed credentials and vars obtained during
	// integration stage. QuerySubscriberFunc an injected func to handle all namecheck
	// requests and a WalletToAccountFunc to handle all WalletToAccount requests
	Client struct {
		Config
		httpClient          *http.Client
		QuerySubscriberFunc QuerySubscriberFunc
		WalletToAccountFunc WalletToAccountFunc
		ctx                 context.Context
		timeout             time.Duration
		logger              io.Writer // for logging purposes
	}

	// ClientOption is a setter func to set Client details like
	// timeout, context, httpClient and logger
	ClientOption func(client *Client)

	loggingTransport struct {
		//ctx    context.Context
		logger io.Writer
		next   http.RoundTripper
	}

	// Config contains configurations that are only needed by Client
	// a USSD client. This is different from tigosdk.Config which contains
	// other configurations that are not needed by the USSD Client but
	// are needed by the pushpay client. This will make it easier for those
	// that needs a single client implementation, they wont necessary need to
	// import code they dont want to use.
	Config struct {
		Username                  string `json:"username"`
		Password                  string `json:"password"`
		AccountName               string `json:"account_name"`
		AccountMSISDN             string `json:"account_msisdn"`
		BrandID                   string `json:"brand_id"`
		BillerCode                string `json:"biller_code"`
		AccountToWalletRequestURL string `json:"disbursement_request_url"`
		AccountToWalletRequestPIN string `json:"disbursement_request_pin"`
		WalletToAccountRequestURL string `json:"collection_request_url"`
		NameCheckRequestURL       string `json:"namecheck_request_url"`
	}



	Service interface {
		SubscriberNameHandler(writer http.ResponseWriter, request *http.Request)

		WalletToAccountHandler(writer http.ResponseWriter, request *http.Request)

		AccountToWalletHandler(ctx context.Context, req AccountToWalletRequest) (resp AccountToWalletResponse, err error)
	}
)

const (
	SubscriberName RequestType = iota
	WalletToAccount

	//debugKey is the value that stores the debugging key is env file
	debugKey            = "DEBUG"
	defaultTimeout      = time.Minute
	SyncLookupResponse  = "SYNC_LOOKUP_RESPONSE"
	SyncBillPayResponse = "SYNC_BILLPAY_RESPONSE"
	REQMFCI             = "REQMFCI"
)

var (

	// defaultCtx is the context used by client when none is set
	// to override this one has to call WithContext method and supply
	// his her own context.Context
	defaultCtx = context.TODO()

	// defaultWriter is an io.Writer used for debugging. When debug mode is
	// set to true i.e DEBUG=true and no io.Writer is provided via
	// WithLogger method this is used.
	defaultWriter = os.Stderr


	// loggingTransport implements http.RoundTripper
	_ http.RoundTripper = (*loggingTransport)(nil)

	// Client implements Service
	_ Service = (*Client)(nil)

	// defaultLoggerTransport is a modified http.Transport that is responsible
	// for logging all requests and responses  that a HTTPClient owned by Client
	// sent and receives
	defaultLoggerTransport = loggingTransport{
		logger: defaultWriter,
		next:   http.DefaultTransport,
	}

	// defaultHttpClient is the client used by library to send Http requests, specifically
	// disbursement requests in case a user does not specify one
	defaultHttpClient = &http.Client{
		Transport: defaultLoggerTransport,
		Timeout:   defaultTimeout,
	}
)

func (l loggingTransport) RoundTrip(request *http.Request) (response *http.Response, err error) {

	if os.Getenv(debugKey) == "true" && request != nil{
		reqDump, err := httputil.DumpRequestOut(request, true)
		if err != nil {
			return nil, err
		}
		_, err = fmt.Fprint(l.logger, fmt.Sprintf("Request %s\n", string(reqDump)))
		if err != nil {
			return nil, err
		}

	}
	defer func() {
		if response != nil && os.Getenv(debugKey) == "true"{
			respDump, err := httputil.DumpResponse(response, true)
			_, err = l.logger.Write([]byte(fmt.Sprintf("Response %s\n",string(respDump))))
			if err != nil {
				return
			}
		}
	}()
	response, err = l.next.RoundTrip(request)
	return
}




// NewClient creates a *Client from the specified Config, WalletToAccountFunc and QuerySubscriberFunc
// You can add numerous available ClientOption like timeout using WithTimeout, logger (which is essentially
// io.Writer) using WithLogger, context using WithContext and httpClient using WithHTTPClient
// if those options are not set explicitly NewClient go with default values
// timeout = time.Minute, context = context.TODO() , logger = os.StdErr and a default http.Client which is
// has been modified to be bale to log requests and responses when debug mode is enabled.
func NewClient(conf Config, collector WalletToAccountFunc, namesHandler QuerySubscriberFunc,
	options ...ClientOption) *Client {
	client := &Client{
		Config:              conf,
		httpClient:          defaultHttpClient,
		QuerySubscriberFunc: namesHandler,
		WalletToAccountFunc: collector,
		ctx:                 defaultCtx,
		timeout:             defaultTimeout,
		logger:              defaultWriter,
	}
	for _, option := range options {
		option(client)
	}

	return client
}

// WithContext set the context to be used by Client in its ops
// this unset the default value which is context.TODO()
// This context value is mostly used by Handlers
func WithContext(ctx context.Context) ClientOption{
	return func(client *Client) {
		client.ctx = ctx
	}
}

// WithTimeout used to set the timeout used by handlers like sending requests to
// Tigo Gateway and back in case of Disbursement or to set the max time for
// handlers QuerySubscriberFunc and WalletToAccountFunc while handling requests from tigo
// the default value is 1 minute
func WithTimeout(timeout time.Duration) ClientOption {
	return func(client *Client) {
		client.timeout = timeout
	}
}

// WithLogger set a logger of user preference but of type io.Writer
// that will be used for debugging use cases. A default value is os.Stderr
// it can be replaced by any io.Writer unless its nil which in that case
// it will be ignored
func WithLogger(out io.Writer) ClientOption {
	return func(client *Client) {
		if out == nil {
			return
		}
		client.logger = out
	}
}

// WithHTTPClient when called unset the present http.Client and replace it
// with c. In case user tries to pass a nil value referencing the client
// i.e WithHTTPClient(nil), it will be ignored and the client wont be replaced
// Note: the new client Transport will be modified. It will be wrapped by another
// middleware that enables client to
func WithHTTPClient(c *http.Client) ClientOption {

	// TODO check if its really necessary to set the default timeout to 1 minute
	//if c.Timeout == 0 {
	//	c.Timeout = defaultTimeout
	//}
	return func(client *Client) {
		if c == nil {
			return
		}

		lt := loggingTransport{
			logger: client.logger,
			next:   c.Transport,
		}

		hc := &http.Client{
			Transport:     lt,
			CheckRedirect: c.CheckRedirect,
			Jar:           c.Jar,
			Timeout:       c.Timeout,
		}
		client.httpClient = hc
	}
}

// SubscriberNameHandler is an http.HandleFunc that handles TigoPesa Namecheck requests
// it uses internally an injected QuerySubscriberFunc.
// It accepts the *http.Request unmarshal it to SubscriberNameRequest before passed to
// QuerySubscriberFunc that returns SubscriberNameResponse and error if error is not null
// this is termed as http.StatusInternalServerError. If its nil the response is marshalled
// to XML formatted SubscriberNameResponse and returned to TigoPesa.
func (client *Client) SubscriberNameHandler(writer http.ResponseWriter, request *http.Request) {
	ctx, cancel := context.WithTimeout(client.ctx, defaultTimeout)
	defer cancel()

	var req SubscriberNameRequest
	xmlBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
	// Try to decode the request body into the struct. If there is an error,
	// respond to the Client with the error message and a 400 status code.
	err = xml.Unmarshal(xmlBody, &req)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

	response, err := client.QuerySubscriberFunc(ctx, req)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

	xmlResponse, err := xml.MarshalIndent(response, "", "  ")
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	//todo: inject logger
	if client.logger != nil && os.Getenv(debugKey) == "true" {
		reqDump, _ := httputil.DumpRequestOut(request, true)
		_, _ = client.logger.Write([]byte(fmt.Sprintf("Request: %s\nResponse: %s\n", string(reqDump), string(xmlResponse))))
	}

	writer.Header().Set("Content-Type", "application/xml")
	_, _ = writer.Write(xmlResponse)

}

func (client *Client) WalletToAccountHandler(writer http.ResponseWriter, request *http.Request) {
	ctx, cancel := context.WithTimeout(client.ctx, defaultTimeout)
	defer cancel()

	var req WalletToAccountRequest
	xmlBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
	// Try to decode the request body into the struct. If there is an error,
	// respond to the Client with the error message and a 400 status code.
	err = xml.Unmarshal(xmlBody, &req)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

	response, err := client.WalletToAccountFunc(ctx, req)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

	xmlResponse, err := xml.MarshalIndent(response, "", "  ")
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	//todo: inject logger
	if client.logger != nil && os.Getenv("DEBUG") == "true" {
		reqDump, _ := httputil.DumpRequestOut(request, true)
		_, _ = client.logger.Write([]byte(fmt.Sprintf("Request: %s\nResponse: %s\n", string(reqDump), string(xmlResponse))))
	}

	writer.Header().Set("Content-Type", "application/xml")
	_, _ = writer.Write(xmlResponse)

}

func (client *Client) AccountToWalletHandler(ctx context.Context, request AccountToWalletRequest) (response AccountToWalletResponse, err error) {
	// Marshal the request body into application/xml
	xmlStr, err := xml.MarshalIndent(request, "", "    ")
	if err != nil {
		return
	}
	xmlStr = []byte(xml.Header + string(xmlStr))

	// Create a HTTP Post Request to be sent to Tigo gateway
	req, err := http.NewRequest(http.MethodPost, client.AccountToWalletRequestURL, bytes.NewBuffer(xmlStr)) // URL-encoded payload
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/xml")

	// Create context with a timeout of 1 minute
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	req = req.WithContext(ctx)

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	xmlBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {

		return
	}

	err = xml.Unmarshal(xmlBody, &response)
	if err != nil {
		return
	}

	return
}
