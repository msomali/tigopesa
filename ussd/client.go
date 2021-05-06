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

var _ Service = (*client)(nil)

type client struct {
	Conf tigosdk.Configs
	HTTPClient *http.Client
	NameCheckHandler NameCheckHandleFunc
	WalletToAccountHandler WalletToAccountFunc
}



func NewClient(configs tigosdk.Configs, httpC *http.Client,nameCheckHandler NameCheckHandleFunc, w2aHandler WalletToAccountFunc) Service{
	return &client{
		Conf:                   configs,
		HTTPClient:             httpC,
		NameCheckHandler:       nameCheckHandler,
		WalletToAccountHandler: w2aHandler,
	}
}

func (c client) QuerySubscriberName(ctx context.Context, request SubscriberNameRequest) (response SubscriberNameResponse, err error) {

	response, err = c.NameCheckHandler(ctx, request)
	return
}

func (c client) QuerySubscriberNameL(ctx context.Context, request *http.Request) (response SubscriberNameResponse, err error) {
	var req SubscriberNameRequest
	xmlBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return response,err
	}
	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err = xml.Unmarshal(xmlBody, &req)
	if err != nil {
		return response,err
	}
	response, err = c.NameCheckHandler(ctx, req)
	return
}

func (c client) WalletToAccount(ctx context.Context, request WalletToAccountRequest) (response WalletToAccountResponse, err error) {
	response, err =c.WalletToAccountHandler(ctx, request)
	return
}

// AccountToWallet this is the implementation of disbursement
func (c client) AccountToWallet(ctx context.Context, request AccountToWalletRequest) (response AccountToWalletResponse, err error) {

	// Marshal the request body into application/xml
	xmlStr, err := xml.MarshalIndent(request, "", "    ")
	if err != nil {
		return
	}
	xmlStr = []byte(xml.Header + string(xmlStr))

	// Create a HTTP Post Request to be sent to Tigo gateway
	req, err := http.NewRequest(http.MethodPost, c.Conf.A2WReqURL, bytes.NewBuffer(xmlStr)) // URL-encoded payload
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


