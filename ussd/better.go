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
)

var _ http.RoundTripper = (*loggingTransport)(nil)

type NameChecker  func(ctx context.Context, request SubscriberNameRequest) (SubscriberNameResponse, error)
type WalletToAccountHandler func(ctx context.Context, request WalletToAccountRequest) (WalletToAccountResponse, error)

type loggingTransport struct {
	logger io.Writer
	next http.RoundTripper
}

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

type BetterClient struct {
	tigosdk.Config
	HttpClient             *http.Client
	NameCheckHandler       func(ctx context.Context, request SubscriberNameRequest) (SubscriberNameResponse, error)
	WalletToAccountHandler func(ctx context.Context, request WalletToAccountRequest) (WalletToAccountResponse, error)
	Logger                 io.Writer // for logging purposes
}

//var defaultNameChecker = func(ctx context.Context, request SubscriberNameRequest) (SubscriberNameResponse, error){
//	return SubscriberNameResponse{
//		XMLName:   xml.Name{},
//		Text:      "",
//		Type:      "",
//		Result:    "",
//		ErrorCode: "",
//		ErrorDesc: "",
//		Msisdn:    "",
//		Flag:      "",
//		Content:   "",
//	},nil
//}
//
//var defaultClient = &BetterClient{
//	Config:                 tigosdk.Config{},
//	HttpClient:             nil,
//	NameCheckHandler:       defaultNameChecker,
//	WalletToAccountHandler: nil,
//	Logger:                 nil,
//}

func NewClientWithLogging(configs tigosdk.Config, logger io.Writer, nameKeeper NameChecker, accountant WalletToAccountHandler) *BetterClient {
	httpClient := &http.Client{
		Transport:     loggingTransport{
			logger: logger,
			next:   http.DefaultTransport,
		},
		Timeout:       defaultTimeout,
	}

	return &BetterClient{
		Config:                 configs,
		HttpClient:             httpClient,
		NameCheckHandler:       nameKeeper,
		WalletToAccountHandler: accountant,
		Logger:                 os.Stdout,
	}
}

var _ BetterService = (*BetterClient)(nil)

type BetterService interface {

	QuerySubscriberNameX(writer http.ResponseWriter, request *http.Request)

	WalletToAccountX(writer http.ResponseWriter,request *http.Request)

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
