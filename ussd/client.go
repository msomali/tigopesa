package ussd

import (
	"bytes"
	"context"
	"encoding/xml"
	"github.com/techcraftt/tigosdk"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	defaultTimeout  = time.Minute
)

var _ Service = (*Client)(nil)

type Client struct {
	Conf           tigosdk.Configs
	HTTPClient     *http.Client
	NameHandleFunc NameCheckHandleFunc
	W2AHandleFunc  WalletToAccountFunc
}


func NewClient(configs tigosdk.Configs, client *http.Client, nameHandler NameCheckHandleFunc, w2aHandler WalletToAccountFunc) *Client{

	//check if http client was provided if not use the default
	if client == nil{
		return &Client{
			Conf:           configs,
			HTTPClient:     http.DefaultClient,
			NameHandleFunc: nameHandler,
			W2AHandleFunc:  w2aHandler,
		}
	}
	return &Client{
		Conf:           configs,
		HTTPClient:     client,
		NameHandleFunc: nameHandler,
		W2AHandleFunc:  w2aHandler,
	}
}

//func (c Client) QuerySubscriberName(ctx context.Context, request SubscriberNameRequest) (response SubscriberNameResponse, err error) {
//	response, err = c.NameHandleFunc(ctx, request)
//	return
//}

func (c Client) QuerySubscriberName(ctx context.Context, request *http.Request) (response SubscriberNameResponse, err error) {
	var req SubscriberNameRequest
	xmlBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return response,err
	}
	// Try to decode the request body into the struct. If there is an error,
	// respond to the Client with the error message and a 400 status code.
	err = xml.Unmarshal(xmlBody, &req)
	if err != nil {
		return response,err
	}
	response = c.NameHandleFunc(ctx, req)
	return
}

//func (c Client) WalletToAccount(ctx context.Context, request WalletToAccountRequest) (response WalletToAccountResponse, err error) {
//	response, err =c.W2AHandleFunc(ctx, request)
//	return
//}

func (c Client) WalletToAccount(ctx context.Context, request *http.Request) (response WalletToAccountResponse, err error) {
	var req WalletToAccountRequest
	xmlBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return response,err
	}
	// Try to decode the request body into the struct. If there is an error,
	// respond to the Client with the error message and a 400 status code.
	err = xml.Unmarshal(xmlBody, &req)
	if err != nil {
		return response,err
	}
	response =c.W2AHandleFunc(ctx, req)
	return
}

// AccountToWallet this is the implementation of disbursement
func (c Client) AccountToWallet(ctx context.Context, request AccountToWalletRequest) (response AccountToWalletResponse, err error) {

	// Marshal the request body into application/xml
	xmlStr, err := xml.MarshalIndent(request, "", "    ")
	if err != nil {
		return
	}
	xmlStr = []byte(xml.Header + string(xmlStr))

	// Create a HTTP Post Request to be sent to Tigo gateway
	req, err := http.NewRequest(http.MethodPost, c.Conf.AccountToWalletRequestURL, bytes.NewBuffer(xmlStr)) // URL-encoded payload
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/xml")

	// Create context with a timeout of 1 minute
	ctx, cancel := context.WithTimeout(ctx,defaultTimeout)
	defer cancel()

	req = req.WithContext(ctx)

	resp, err := c.HTTPClient.Do(req)
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


