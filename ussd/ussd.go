/*
 * MIT License
 *
 * Copyright (c) 2021 TechCraft Technologies Co. Ltd
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 *
 */

package ussd

import (
	"context"
	"encoding/xml"
	"github.com/techcraftlabs/base"
	"net/http"
	"time"
)

const (
	syncLookupResponse  = "SYNC_LOOKUP_RESPONSE"
	syncBillPayResponse = "SYNC_BILLPAY_RESPONSE"

	ErrNameNotRegistered = "error010"
	ErrNameInvalidFormat = "error030"
	ErrNameUserSuspended = "error030"
	NoNamecheckErr       = "error000"

	ErrSuccessTxn               = "error000"
	ErrServiceNotAvailable      = "error001"
	ErrInvalidCustomerRefNumber = "error010"
	ErrCustomerRefNumLocked     = "error011"
	ErrInvalidAmount            = "error012"
	ErrAmountInsufficient       = "error013"
	ErrAmountTooHigh            = "error014"
	ErrAmountTooLow             = "error015"
	ErrInvalidPayment           = "error016"
	ErrGeneralError             = "error100"
	ErrRetryConditionNoResponse = "error111"
)

var (
	_ PaymentHandler   = (*PaymentHandleFunc)(nil)
	_ NameQueryHandler = (*NameQueryFunc)(nil)
	_ Service          = (*Client)(nil)
)

type (
	Service interface {
		NameQueryServeHTTP(writer http.ResponseWriter, request *http.Request)
		PaymentServeHTTP(writer http.ResponseWriter, request *http.Request)
	}

	PaymentHandler interface {
		HandlePayRequest(ctx context.Context, request PayRequest) (PayResponse, error)
	}

	PaymentHandleFunc func(ctx context.Context, request PayRequest) (PayResponse, error)

	NameQueryHandler interface {
		HandleNameQuery(ctx context.Context, request NameRequest) (NameResponse, error)
	}

	NameQueryFunc func(ctx context.Context, request NameRequest) (NameResponse, error)

	nameRequest struct {
		XMLName             xml.Name `xml:"COMMAND"`
		Text                string   `xml:",chardata"`
		Type                string   `xml:"TYPE"`
		Msisdn              string   `xml:"MSISDN"`
		CompanyName         string   `xml:"COMPANYNAME"`
		CustomerReferenceID string   `xml:"CUSTOMERREFERENCEID"`
	}

	NameRequest struct {
		Msisdn              string `xml:"MSISDN"`
		CompanyName         string `xml:"COMPANYNAME"`
		CustomerReferenceID string `xml:"CUSTOMERREFERENCEID"`
	}

	NameResponse struct {
		Result    string `xml:"RESULT"`
		ErrorCode string `xml:"ERRORCODE"`
		ErrorDesc string `xml:"ERRORDESC"`
		Msisdn    string `xml:"MSISDN"`
		Flag      string `xml:"FLAG"`
		Content   string `xml:"CONTENT"`
	}

	nameResponse struct {
		XMLName   xml.Name `xml:"COMMAND"`
		Text      string   `xml:",chardata"`
		Type      string   `xml:"TYPE"`
		Result    string   `xml:"RESULT"`
		ErrorCode string   `xml:"ERRORCODE"`
		ErrorDesc string   `xml:"ERRORDESC"`
		Msisdn    string   `xml:"MSISDN"`
		Flag      string   `xml:"FLAG"`
		Content   string   `xml:"CONTENT"`
	}

	PayRequest struct {
		TxnID               string  `xml:"TXNID"`
		Msisdn              string  `xml:"MSISDN"`
		Amount              float64 `xml:"AMOUNT"`
		CompanyName         string  `xml:"COMPANYNAME"`
		CustomerReferenceID string  `xml:"CUSTOMERREFERENCEID"`
		SenderName          string  `xml:"SENDERNAME"`
	}

	payRequest struct {
		XMLName             xml.Name `xml:"COMMAND"`
		Text                string   `xml:",chardata"`
		TYPE                string   `xml:"TYPE"`
		TxnID               string   `xml:"TXNID"`
		Msisdn              string   `xml:"MSISDN"`
		Amount              float64  `xml:"AMOUNT"`
		CompanyName         string   `xml:"COMPANYNAME"`
		CustomerReferenceID string   `xml:"CUSTOMERREFERENCEID"`
		SenderName          string   `xml:"SENDERNAME"`
	}

	PayResponse struct {
		TxnID            string `xml:"TXNID"`
		RefID            string `xml:"REFID"`
		Result           string `xml:"RESULT"`
		ErrorCode        string `xml:"ERRORCODE"`
		ErrorDescription string `xml:"ERRORDESCRIPTION"`
		Msisdn           string `xml:"MSISDN"`
		Flag             string `xml:"FLAG"`
		Content          string `xml:"CONTENT"`
	}

	payResponse struct {
		XMLName          xml.Name `xml:"COMMAND"`
		Text             string   `xml:",chardata"`
		Type             string   `xml:"TYPE"`
		TxnID            string   `xml:"TXNID"`
		RefID            string   `xml:"REFID"`
		Result           string   `xml:"RESULT"`
		ErrorCode        string   `xml:"ERRORCODE"`
		ErrorDescription string   `xml:"ERRORDESCRIPTION"`
		Msisdn           string   `xml:"MSISDN"`
		Flag             string   `xml:"FLAG"`
		Content          string   `xml:"CONTENT"`
	}

	Config struct {
		AccountName   string
		AccountMSISDN string
		BillerNumber  string
		RequestURL    string
		NamecheckURL  string
	}

	Client struct {
		rv   base.Receiver
		rp   base.Replier
		base *base.Client
		*Config
		ph PaymentHandler
		nh NameQueryHandler
	}
)

func transformNameRequest(request nameRequest) NameRequest {
	return NameRequest{
		Msisdn:              request.Msisdn,
		CompanyName:         request.CompanyName,
		CustomerReferenceID: request.CustomerReferenceID,
	}
}

func transformPayRequest(request payRequest) PayRequest {
	return PayRequest{
		TxnID:               request.TxnID,
		Msisdn:              request.Msisdn,
		Amount:              request.Amount,
		CompanyName:         request.CompanyName,
		CustomerReferenceID: request.CustomerReferenceID,
		SenderName:          request.SenderName,
	}
}

func transformToXMLNameResponse(response NameResponse) nameResponse {
	return nameResponse{
		Type:      syncLookupResponse,
		Result:    response.Result,
		ErrorCode: response.ErrorCode,
		ErrorDesc: response.ErrorDesc,
		Msisdn:    response.Msisdn,
		Flag:      response.Flag,
		Content:   response.Content,
	}
}

func transformToXMLPayResponse(response PayResponse) payResponse {
	return payResponse{
		Type:             syncBillPayResponse,
		TxnID:            response.TxnID,
		RefID:            response.RefID,
		Result:           response.Result,
		ErrorCode:        response.ErrorCode,
		ErrorDescription: response.ErrorDescription,
		Msisdn:           response.Msisdn,
		Flag:             response.Flag,
		Content:          response.Content,
	}
}

func NewClient(config *Config, handler PaymentHandler, queryHandler NameQueryHandler, opts ...ClientOption) *Client {
	client := &Client{

		Config: config,
		ph:     handler,
		nh:     queryHandler,
		base:   base.NewClient(),
	}

	for _, opt := range opts {
		opt(client)
	}

	lg, dm := client.base.Logger, client.base.DebugMode

	client.rp = base.NewReplier(lg, dm)
	client.rv = base.NewReceiver(lg, dm)

	return client
}

func (c *Client) SetNameQueryHandler(nh NameQueryHandler) {
	c.nh = nh
}

func (c *Client) SetPaymentRequestHandler(ph PaymentHandler) {
	c.ph = ph
}

func (handler PaymentHandleFunc) HandlePayRequest(ctx context.Context, request PayRequest) (PayResponse, error) {
	return handler(ctx, request)
}

func (handler NameQueryFunc) HandleNameQuery(ctx context.Context, request NameRequest) (NameResponse, error) {
	return handler(ctx, request)
}

func (c *Client) NameQueryServeHTTP(writer http.ResponseWriter, request *http.Request) {

	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()
	var req nameRequest

	_, err := c.rv.Receive(ctx, "name query", request, &req)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

	response, err := c.nh.HandleNameQuery(ctx, transformNameRequest(req))

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

	var opts []base.ResponseOption
	headers := map[string]string{
		"Content-Type": "application/xml",
	}
	headersOpts := base.WithResponseHeaders(headers)
	opts = append(opts, headersOpts)

	payload := transformToXMLNameResponse(response)
	res := base.NewResponse(200, payload, opts...)

	c.rp.Reply(writer, res)

}

func (c *Client) PaymentServeHTTP(writer http.ResponseWriter, request *http.Request) {

	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()

	var req payRequest

	_, err := c.rv.Receive(ctx, "payment request", request, &req)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := c.ph.HandlePayRequest(ctx, transformPayRequest(req))

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

	var opts []base.ResponseOption
	headers := map[string]string{
		"Content-Type": "application/xml",
	}
	headersOpts := base.WithResponseHeaders(headers)
	opts = append(opts, headersOpts)
	payload := transformToXMLPayResponse(response)
	res := base.NewResponse(200, payload, opts...)

	c.rp.Reply(writer, res)

}
