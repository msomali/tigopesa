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
	"github.com/techcraftlabs/tigopesa/internal"
	"net/http"
	"time"
)

const (
	SyncLookupResponse  = "SYNC_LOOKUP_RESPONSE"
	SyncBillPayResponse = "SYNC_BILLPAY_RESPONSE"

	//NameQuery Error Codes

	ErrNameNotRegistered = "error010"
	ErrNameInvalidFormat = "error030"
	ErrNameUserSuspended = "error030"
	NoNamecheckErr       = "error000"

	//WalletToAccount error codes

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
		HandleNameQuery(writer http.ResponseWriter, request *http.Request)

		HandlePayment(writer http.ResponseWriter, request *http.Request)
	}

	PaymentHandler interface {
		PaymentRequest(ctx context.Context, request PayRequest) (PayResponse, error)
	}

	PaymentHandleFunc func(ctx context.Context, request PayRequest) (PayResponse, error)

	NameQueryHandler interface {
		NameQuery(ctx context.Context, request NameRequest) (NameResponse, error)
	}

	NameQueryFunc func(ctx context.Context, request NameRequest) (NameResponse, error)

	NameRequest struct {
		XMLName             xml.Name `xml:"COMMAND"`
		Text                string   `xml:",chardata"`
		Type                string   `xml:"TYPE"`
		Msisdn              string   `xml:"MSISDN"`
		CompanyName         string   `xml:"COMPANYNAME"`
		CustomerReferenceID string   `xml:"CUSTOMERREFERENCEID"`
	}

	NameResponse struct {
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
		base *internal.BaseClient
		*Config
		PaymentHandler   PaymentHandler
		NameQueryHandler NameQueryHandler
	}
)

func NewClient(config *Config, handler PaymentHandler, queryHandler NameQueryHandler, opts ...ClientOption) *Client {
	client := &Client{

		Config:           config,
		PaymentHandler:   handler,
		NameQueryHandler: queryHandler,
		base:             internal.NewBaseClient(),
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

func (handler PaymentHandleFunc) PaymentRequest(ctx context.Context, request PayRequest) (PayResponse, error) {
	return handler(ctx, request)
}

func (handler NameQueryFunc) NameQuery(ctx context.Context, request NameRequest) (NameResponse, error) {
	return handler(ctx, request)
}

func (client *Client) HandleNameQuery(writer http.ResponseWriter, request *http.Request) {

	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()
	var req NameRequest

	err := client.base.Receive(request, internal.NameQueryRequest, internal.XmlPayload, &req)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

	response, err := client.NameQueryHandler.NameQuery(ctx, req)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

	resp := internal.NewResponse(200, response, internal.XmlPayload)
	client.base.Reply("name query response", resp, writer)

}

func (client *Client) HandlePayment(writer http.ResponseWriter, request *http.Request) {

	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()

	var req PayRequest

	err := client.base.Receive(request, internal.PaymentRequest, internal.XmlPayload, &req)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

	response, err := client.PaymentHandler.PaymentRequest(ctx, req)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

	resp := internal.NewResponse(200, response, internal.XmlPayload)

	client.base.Reply("ussd payment response", resp, writer)

}
