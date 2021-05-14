package ussd
//
//import (
//	"bytes"
//	"context"
//	"encoding/xml"
//	"fmt"
//	"github.com/techcraftt/tigosdk"
//	"io/ioutil"
//	"net/http"
//)
//
//
//type NameCheckHandleFunc func(context.Context, SubscriberNameRequest) SubscriberNameResponse
//
//type WalletToAccountFunc func(ctx context.Context, request WalletToAccountRequest) WalletToAccountResponse
//
//type Service interface {
//	QuerySubscriberName(ctx context.Context, request *http.Request) (resp SubscriberNameResponse, err error)
//
//	WalletToAccount(ctx context.Context, request *http.Request) (resp WalletToAccountResponse, err error)
//
//	AccountToWallet(ctx context.Context, req AccountToWalletRequest) (resp AccountToWalletResponse, err error)
//}
//
//
//
//
//
//var _ Service = (*Client)(nil)
//
//type Client struct {
//	Conf                      tigosdk.Config
//	HTTPClient                *http.Client
//	NameCheckReqHandler       NameCheckHandleFunc
//	WalletToAccountReqHandler WalletToAccountFunc
//}
//
//func NewClient(configs tigosdk.Config, client *http.Client, nameHandler NameCheckHandleFunc, w2aHandler WalletToAccountFunc) *Client {
//
//	//check if http client was provided if not use the default
//	if client == nil {
//		return &Client{
//			Conf:                      configs,
//			HTTPClient:                http.DefaultClient,
//			NameCheckReqHandler:       nameHandler,
//			WalletToAccountReqHandler: w2aHandler,
//		}
//	}
//	return &Client{
//		Conf:                      configs,
//		HTTPClient:                client,
//		NameCheckReqHandler:       nameHandler,
//		WalletToAccountReqHandler: w2aHandler,
//	}
//}
//
//// QuerySubscriberName handles the namecheck requests from Tigo, unmarshal the xml
//// into SubscriberNameRequest that is then passed to an injected NameCheckHandleFunc
//// that has the logic on what to do with the request and knows what should be the
//// SubscriberNameResponse for each request. So NameCheckHandleFunc computes and return
//// a proper response which is then marshalled to XML format and sent as a response to tigo
//func (c *Client) QuerySubscriberName(ctx context.Context, request *http.Request) (response SubscriberNameResponse, err error) {
//	var req SubscriberNameRequest
//	xmlBody, err := ioutil.ReadAll(request.Body)
//	if err != nil {
//		return response, err
//	}
//	// Try to decode the request body into the struct. If there is an error,
//	// respond to the Client with the error message and a 400 status code.
//	err = xml.Unmarshal(xmlBody, &req)
//	if err != nil {
//		return response, err
//	}
//	response = c.NameCheckReqHandler(ctx, req)
//	return
//}
//
//// WalletToAccountFunc is expected to replace WalletToAccount
//// This is very experimental ATM
//// Why experimental? because the response is handled internally by the library which is a bad practice
//// the ideal should have been
//// 1. library receive http.Request from tigo unmarshal it to consumable struct
//// 2. pass the request to the WalletToAccountReqHandler that check the req and
//// return a struct that should be sent back to tigo. but before then it is properly
//// marshalled to xml format and request headers are set properly before being sent to back
//// synchronously to tigo.
//func (c *Client) WalletToAccountFunc(writer http.ResponseWriter, request *http.Request) {
//
//	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
//	defer cancel()
//
//	var req WalletToAccountRequest
//	xmlBody, err := ioutil.ReadAll(request.Body)
//	if err != nil {
//		http.Error(writer, err.Error(), http.StatusInternalServerError)
//	}
//	// Try to decode the request body into the struct. If there is an error,
//	// respond to the Client with the error message and a 400 status code.
//	err = xml.Unmarshal(xmlBody, &req)
//	if err != nil {
//		http.Error(writer, err.Error(), http.StatusInternalServerError)
//	}
//
//	response := c.WalletToAccountReqHandler(ctx, req)
//
//	xmlResponse, err := xml.MarshalIndent(response, "", "  ")
//	if err != nil {
//		http.Error(writer, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	writer.Header().Set("Content-Type", "application/xml")
//	_, _ = writer.Write(xmlResponse)
//
//}
//
//func (c *Client) WalletToAccount(ctx context.Context, request *http.Request) (response WalletToAccountResponse, err error) {
//	var req WalletToAccountRequest
//	xmlBody, err := ioutil.ReadAll(request.Body)
//	if err != nil {
//		return response, err
//	}
//	// Try to decode the request body into the struct. If there is an error,
//	// respond to the Client with the error message and a 400 status code.
//	err = xml.Unmarshal(xmlBody, &req)
//	if err != nil {
//		return response, err
//	}
//	response = c.WalletToAccountReqHandler(ctx, req)
//	return
//}
//
//// AccountToWallet this is the implementation of disbursement
//func (c *Client) AccountToWallet(ctx context.Context, request AccountToWalletRequest) (response AccountToWalletResponse, err error) {
//
//	// Marshal the request body into application/xml
//	xmlStr, err := xml.MarshalIndent(request, "", "    ")
//	if err != nil {
//		return
//	}
//	xmlStr = []byte(xml.Header + string(xmlStr))
//
//	fmt.Printf("request is: %s\n", xmlStr)
//
//	// Create a HTTP Post Request to be sent to Tigo gateway
//	req, err := http.NewRequest(http.MethodPost, c.Conf.AccountToWalletRequestURL, bytes.NewBuffer(xmlStr)) // URL-encoded payload
//	if err != nil {
//		return
//	}
//	req.Header.Add("Content-Type", "application/xml")
//
//	// Create context with a timeout of 1 minute
//	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
//	defer cancel()
//
//	req = req.WithContext(ctx)
//
//	resp, err := c.HTTPClient.Do(req)
//	if err != nil {
//		return
//	}
//
//	defer resp.Body.Close()
//	xmlBody, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//
//		return
//	}
//
//	fmt.Printf("response is %s\n", string(xmlBody))
//
//	err = xml.Unmarshal(xmlBody, &response)
//	if err != nil {
//		return
//	}
//
//	return
//}
