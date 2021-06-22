package tigosdk

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
	//      case NameQuery:
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
		*tigo.BaseClient
		QuerySubscriberFunc QuerySubscriberFunc
		WalletToAccountFunc WalletToAccountFunc
	}

	Service interface {
		SubscriberNameHandler(writer http.ResponseWriter, request *http.Request)

		WalletToAccountHandler(writer http.ResponseWriter, request *http.Request)

		AccountToWalletHandler(ctx context.Context, req AccountToWalletRequest) (resp AccountToWalletResponse, err error)
	}
)

const (
	NameQuery RequestType = iota
	WalletToAccount
	PushPay
	Disbursement

	//debugKey is the value that stores the debugging key is env file
	debugKey = "DEBUG"
)

var (

	// loggingTransport implements http.RoundTripper
	//_ http.RoundTripper = (*loggingTransport)(nil)

	// Client implements Service
	_ Service = (*Client)(nil)

	_ WalletToAccountProvider = (*WalletToAccountFunc)(nil)

	//QuerySubscriberFunc implements QuerySubscriberNameProvider
	_ QuerySubscriberNameProvider = (*QuerySubscriberFunc)(nil)
)

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
func NewClient(client *tigo.BaseClient, collector WalletToAccountFunc, namesHandler QuerySubscriberFunc) *Client {
	return &Client{
		BaseClient:          client,
		QuerySubscriberFunc: namesHandler,
		WalletToAccountFunc: collector,
	}
}

// SubscriberNameHandler is an http.HandleFunc that handles TigoPesa Namecheck requests
// it uses internally an injected QuerySubscriberFunc.
// It accepts the *http.Request unmarshal it to SubscriberNameRequest before passed to
// QuerySubscriberFunc that returns SubscriberNameResponse and error if error is not null
// this is termed as http.StatusInternalServerError. If its nil the response is marshalled
// to XML formatted SubscriberNameResponse and returned to TigoPesa.
func (client *Client) SubscriberNameHandler(writer http.ResponseWriter, request *http.Request) {
	ctx, cancel := context.WithTimeout(client.Ctx, client.Timeout)
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
	if client.Logger != nil && os.Getenv(debugKey) == "true" {
		reqDump, _ := httputil.DumpRequestOut(request, true)
		_, _ = client.Logger.Write([]byte(fmt.Sprintf("Request: %s\nResponse: %s\n", string(reqDump), string(xmlResponse))))
	}

	writer.Header().Set("Content-Type", "application/xml")
	_, _ = writer.Write(xmlResponse)

}

func (client *Client) WalletToAccountHandler(writer http.ResponseWriter, request *http.Request) {

	ctx, cancel := context.WithTimeout(client.Ctx, client.Timeout)
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
	if client.Logger != nil && os.Getenv("DEBUG") == "true" {
		reqDump, _ := httputil.DumpRequestOut(request, true)
		_, _ = client.Logger.Write([]byte(fmt.Sprintf("Request: %s\nResponse: %s\n", string(reqDump), string(xmlResponse))))
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
	ctx, cancel := context.WithTimeout(ctx, client.Timeout)
	defer cancel()

	req = req.WithContext(ctx)

	resp, err := client.HttpClient.Do(req)
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
	ctx, cancel := context.WithTimeout(ctx, client.Timeout)
	defer cancel()
	return func(writer http.ResponseWriter, request *http.Request) {
		switch requestType {
		case NameQuery:
			client.SubscriberNameHandler(writer, request)
		case WalletToAccount:
			client.WalletToAccountHandler(writer, request)
		default:
			http.Error(writer, "unknown request type", http.StatusInternalServerError)
		}
	}
}
