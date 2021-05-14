package ussd

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"github.com/techcraftt/tigosdk"
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
	// func HandleRequest(ctx context.Context, requestType RequestType) http.HandlerFunc {
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

	NameCheckHandler       func(ctx context.Context, request SubscriberNameRequest) (SubscriberNameResponse, error)
	WalletToAccountHandler func(ctx context.Context, request WalletToAccountRequest) (WalletToAccountResponse, error)

	BetterClient struct {
		tigosdk.Config
		HttpClient             *http.Client
		NameCheckHandler       NameCheckHandler
		WalletToAccountHandler WalletToAccountHandler
		logger                 io.Writer // for logging purposes
	}

	loggingTransport struct {
		logger io.Writer
		next   http.RoundTripper
	}

	// Config contains configurations that are only needed by BetterClient
	// a USSD client. This is different from tigosdk.Config which contains
	// other configurations that are not needed by the USSD BetterClient but
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
)

const (
	SubscriberName RequestType = iota
	WalletToAccount

	//
	defaultTimeout        = time.Minute
	SYNC_LOOKUP_RESPONSE  = "SYNC_LOOKUP_RESPONSE"
	SYNC_BILLPAY_RESPONSE = "SYNC_BILLPAY_RESPONSE"
	REQMFCI               = "REQMFCI"
)

// loggingTransport implements http.RoundTripper
var _ http.RoundTripper = (*loggingTransport)(nil)

// BetterClient implements BetterService
var _ BetterService = (*BetterClient)(nil)

func (l loggingTransport) RoundTrip(request *http.Request) (response *http.Response, err error) {
	defer func() {
		var reqDump, respDump []byte

		if request != nil {
			reqDump, _ = httputil.DumpRequestOut(request, true)
		}
		if response != nil {
			respDump, _ = httputil.DumpResponse(response, true)
		}
		_, _ = l.logger.Write([]byte(fmt.Sprintf("Request: %s\nResponse: %s\n", string(reqDump), string(respDump))))
	}()
	response, err = l.next.RoundTrip(request)
	return
}

func NewBetterClient(configs tigosdk.Config, client *http.Client, out io.Writer,
	nameKeeper NameCheckHandler, accountant WalletToAccountHandler) *BetterClient {

	c := &BetterClient{
		Config:                 configs,
		NameCheckHandler:       nameKeeper,
		WalletToAccountHandler: accountant,
	}

	if client == nil {

		// No *http.Client is set and neither is logger (io.Writer) in our case
		// So the http.DefaultClient will then be set  to be used
		if out == nil {
			c.SetHTTPClient(http.DefaultClient)
		}

		// here the http.Client is not set but the logger (io.Writer) is set
		// so a a new http.Client will be spun with a modified Transport that
		// enables logging of both requests and responses
		httpClient := &http.Client{
			Transport: loggingTransport{
				logger: out,
				next:   http.DefaultTransport,
			},
			Timeout: defaultTimeout,
		}
		c.SetLogger(out)
		c.SetHTTPClient(httpClient)
	}

	if client != nil {
		if out != nil {
			c.SetLogger(out)
		}
		c.SetHTTPClient(client)
	}

	return c
}

func NewClientWithLogging(configs tigosdk.Config, out io.Writer, nameKeeper NameCheckHandler, accountant WalletToAccountHandler) *BetterClient {
	httpClient := &http.Client{
		Transport: loggingTransport{
			logger: out,
			next:   http.DefaultTransport,
		},
		Timeout: defaultTimeout,
	}

	return &BetterClient{
		Config:                 configs,
		HttpClient:             httpClient,
		NameCheckHandler:       nameKeeper,
		WalletToAccountHandler: accountant,
		logger:                 os.Stdout,
	}
}

func (client *BetterClient) SetHTTPClient(httpClient *http.Client) {
	client.HttpClient = httpClient
}

// SetLogger set custom logs destination.
func (client *BetterClient) SetLogger(out io.Writer) {
	client.logger = out
}

type BetterService interface {
	QuerySubscriberNameX(writer http.ResponseWriter, request *http.Request)

	WalletToAccountX(writer http.ResponseWriter, request *http.Request)

	AccountToWalletX(ctx context.Context, req AccountToWalletRequest) (resp AccountToWalletResponse, err error)
}

func (client BetterClient) QuerySubscriberNameX(writer http.ResponseWriter, request *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
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

	response, err := client.NameCheckHandler(ctx, req)

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

func (client BetterClient) WalletToAccountX(writer http.ResponseWriter, request *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
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

	response, err := client.WalletToAccountHandler(ctx, req)

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

func (client BetterClient) AccountToWalletX(ctx context.Context, request AccountToWalletRequest) (response AccountToWalletResponse, err error) {
	// Marshal the request body into application/xml
	xmlStr, err := xml.MarshalIndent(request, "", "    ")
	if err != nil {
		return
	}
	xmlStr = []byte(xml.Header + string(xmlStr))

	fmt.Printf("request is: %s\n", xmlStr)

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

	resp, err := client.HttpClient.Do(req)
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
