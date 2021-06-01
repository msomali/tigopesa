package ussd

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"github.com/techcraftt/tigosdk/pkg/tigo"
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
	// func (pkg *Client) HandleRequest(ctx context.Context, requestType RequestType) http.HandlerFunc {
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

	// QuerySubscriberNameProvider is an  interface that has the signature of QuerySubscriberFunc
	// which also implements it.
	QuerySubscriberNameProvider interface {
		QuerySubscriberName(ctx context.Context, request SubscriberNameRequest) (SubscriberNameResponse, error)
	}

	// QuerySubscriberFunc takes a SubscriberNameRequest and return SubscriberNameResponse or error
	// There is no default implementation of this. It should be implemented and injected to Client as
	// Client.QuerySubscriberFunc. This carry the whole logic that response from TigoPesa Namecheck
	// requests.
	// The likely scenario may be taking SubscriberNameRequest and run it against the user database
	// or auth system and the respond appropriately.
	// an example may look like
	//      func (app *App)QuerySubscriber(ctx context.Context, request SubscriberNameRequest) (resp SubscriberNameResponse, err error){
	//           err = app.database.find(ctx, request.CustomerReferenceID)
	//           if err != nil {
	//               return resp, err
	//           }
	//      }
	// the injected function will then be called by Client.SubscriberNameHandler and its implementation gets
	// executed in terms of err != nil . http.StatusInternalServerError will be sent to TigoPesa. So make
	// sure the errors reported by this func are system failure only like database connection but in terms of
	// things like Customer Not Found or Invalid Reference Number err should be nil and appropriate SubscriberNameResponse
	// should be used. example
	//
	// resp := SubscriberNameResponse {
	//		Type:      REQMFCI,
	//		Result:    "TF",
	//		ErrorCode: NAMECHECK_NOT_REGISTERED,
	//		ErrorDesc: "Customer Details Not Found",
	//		Msisdn:    "",
	//		Flag:      "N",
	//		Content:   "The provided reference number is unknown",
	//	}
	QuerySubscriberFunc func(ctx context.Context, request SubscriberNameRequest) (SubscriberNameResponse, error)

	WalletToAccountProvider interface {
		WalletToAccount(ctx context.Context, request WalletToAccountRequest) (WalletToAccountResponse, error)
	}
	// WalletToAccountFunc handles all the collection requests sent from TigoPesa. The http requests are firstly
	// received by Client.WalletToAccountHandler properly marshalled to WalletToAccountRequest then passed to
	// this func which after processing the requests returns an error in case there is a system failure. An error
	// should not be returned for cases like Insufficient funds, Account is barred etc etc. in case of those errors
	// just like in QuerySubscriberFunc the error should be nil and these other "errors" should be include in
	// WalletToAccountResponse
	WalletToAccountFunc func(ctx context.Context, request WalletToAccountRequest) (WalletToAccountResponse, error)

	// Client is a ussd pkg that assembles together necessary parts
	// needed in performing 3 operations
	// 1. handling of NameCheckRequest, 2. Handling of Collection (WalletToAccount)
	// requests and 3. Making disbursement requests
	// it contains Config to store all needed credentials and vars obtained during
	// integration stage. QuerySubscriberFunc an injected func to handle all namecheck
	// requests and a WalletToAccountFunc to handle all WalletToAccount requests
	Client struct {
		tigo.BaseClient
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
		logger io.Writer
		next   http.RoundTripper
	}

	// Config contains configurations that are only needed by Client
	// a USSD pkg. This is different from tigosdk.Config which contains
	// other configurations that are not needed by the USSD Client but
	// are needed by the pushpay pkg. This will make it easier for those
	// that needs a single pkg implementation, they wont necessary need to
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

	// loggingTransport implements http.RoundTripper
	_ http.RoundTripper = (*loggingTransport)(nil)

	// Client implements Service
	_ Service = (*Client)(nil)

	_ WalletToAccountProvider = (*WalletToAccountFunc)(nil)

	//QuerySubscriberFunc implements QuerySubscriberNameProvider
	_ QuerySubscriberNameProvider = (*QuerySubscriberFunc)(nil)

	// defaultCtx is the context used by pkg when none is set
	// to override this one has to call WithContext method and supply
	// his her own context.Context
	defaultCtx = context.TODO()

	// defaultWriter is an io.Writer used for debugging. When debug mode is
	// set to true i.e DEBUG=true and no io.Writer is provided via
	// WithLogger method this is used.
	defaultWriter = os.Stderr

	// defaultLoggerTransport is a modified http.Transport that is responsible
	// for logging all requests and responses  that a HTTPClient owned by Client
	// sent and receives
	defaultLoggerTransport = loggingTransport{
		logger: defaultWriter,
		next:   http.DefaultTransport,
	}

	// defaultHttpClient is the pkg used by library to send Http requests, specifically
	// disbursement requests in case a user does not specify one
	defaultHttpClient = &http.Client{
		Transport: defaultLoggerTransport,
		Timeout:   defaultTimeout,
	}
)

func (l loggingTransport) RoundTrip(request *http.Request) (response *http.Response, err error) {

	if os.Getenv(debugKey) == "true" && request != nil {
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
		if response != nil && os.Getenv(debugKey) == "true" {
			respDump, err := httputil.DumpResponse(response, true)
			_, err = l.logger.Write([]byte(fmt.Sprintf("Response %s\n", string(respDump))))
			if err != nil {
				return
			}
		}
	}()
	response, err = l.next.RoundTrip(request)
	return
}

func (w WalletToAccountFunc) WalletToAccount(ctx context.Context, request WalletToAccountRequest) (WalletToAccountResponse, error) {
	return w(ctx, request)
}

func (q QuerySubscriberFunc) QuerySubscriberName(ctx context.Context, request SubscriberNameRequest) (SubscriberNameResponse, error) {
	return q(ctx, request)
}

// NewClient creates a *Client from the specified Config, WalletToAccountFunc and QuerySubscriberFunc
// You can add numerous available ClientOption like timeout using WithTimeout, logger (which is essentially
// io.Writer) using WithLogger, context using WithContext and httpClient using WithHTTPClient
// if those options are not set explicitly NewClient go with default values
// timeout = time.Minute, context = context.TODO() , logger = os.Stderr and a default http.Client which is
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
func WithContext(ctx context.Context) ClientOption {
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
// with c. In case user tries to pass a nil value referencing the pkg
// i.e WithHTTPClient(nil), it will be ignored and the pkg wont be replaced
// Note: the new pkg Transport will be modified. It will be wrapped by another
// middleware that enables pkg to
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
	ctx, cancel := context.WithTimeout(client.ctx, client.timeout)
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
	ctx, cancel := context.WithTimeout(client.ctx, client.timeout)
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

// AccountToWalletHandler handles all disbursement requests
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
	ctx, cancel := context.WithTimeout(ctx, client.timeout)
	defer cancel()

	req = req.WithContext(ctx)

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)
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

// HandleRequest is experimental no guarantees
// For reliability use SubscriberNameHandler and WalletToAccountHandler
func (client *Client) HandleRequest(ctx context.Context, requestType RequestType) http.HandlerFunc {
	ctx, cancel := context.WithTimeout(ctx, client.timeout)
	defer cancel()
	return func(writer http.ResponseWriter, request *http.Request) {
		switch requestType {
		case SubscriberName:
			client.SubscriberNameHandler(writer, request)
		case WalletToAccount:
			client.WalletToAccountHandler(writer, request)
		default:
			http.Error(writer, "unknown request type", http.StatusInternalServerError)
		}
	}
}
