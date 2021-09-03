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

package disburse

import (
	"context"
	"encoding/xml"
	"github.com/techcraftlabs/tigopesa/internal"
	"net/http"
)

const (
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

	requestType    = "REQMFCI"
	senderLanguage = "EN"
)

type (
	Service interface {
		Disburse(ctx context.Context, request Request) (Response, error)
	}

	Config struct {
		AccountName   string
		AccountMSISDN string
		BrandID       string
		PIN           string
		RequestURL    string
	}

	Request struct {
		ReferenceID string  `json:"reference"`
		MSISDN      string  `json:"msisdn"`
		Amount      float64 `json:"amount"`
	}

	disburseRequest struct {
		XMLName     xml.Name `xml:"COMMAND"`
		Text        string   `xml:",chardata"`
		Type        string   `xml:"TYPE"`
		ReferenceID string   `xml:"REFERENCEID"`
		Msisdn      string   `xml:"MSISDN"`
		PIN         string   `xml:"PIN"`
		Msisdn1     string   `xml:"MSISDN1"`
		Amount      float64  `xml:"AMOUNT"`
		SenderName  string   `xml:"SENDERNAME"`
		Language1   string   `xml:"LANGUAGE1"`
		BrandID     string   `xml:"BRAND_ID"`
	}

	Response struct {
		XMLName     xml.Name `xml:"COMMAND" json:"-"`
		Text        string   `xml:",chardata" json:"-"`
		Type        string   `xml:"TYPE" json:"type"`
		ReferenceID string   `xml:"REFERENCEID" json:"reference,omitempty"`
		TxnID       string   `xml:"TXNID" json:"id,omitempty"`
		TxnStatus   string   `xml:"TXNSTATUS" json:"status,omitempty"`
		Message     string   `xml:"MESSAGE" json:"message,omitempty"`
	}

	Client struct {
		*Config
		base *internal.BaseClient
	}
)

func NewClient(config *Config, opts ...ClientOption) *Client {
	client := &Client{
		Config:     config,
		base: internal.NewBaseClient(),
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

func (client *Client) Disburse(ctx context.Context, request Request) (response Response, err error) {
	var reqOpts []internal.RequestOption
	ctxOpt := internal.WithRequestContext(ctx)
	headers := map[string]string{
		"Content-Type": "application/xml",
	}
	headersOpt := internal.WithRequestHeaders(headers)
	reqOpts = append(reqOpts, ctxOpt, headersOpt)
	r := disburseRequest{
		Type:        requestType,
		ReferenceID: request.ReferenceID,
		Msisdn:      client.Config.AccountMSISDN,
		PIN:         client.Config.PIN,
		Msisdn1:     request.MSISDN,
		Amount:      request.Amount,
		SenderName:  client.Config.AccountName,
		Language1:   senderLanguage,
		BrandID:     client.Config.BrandID,
	}

	req := internal.NewRequest(ctx, http.MethodPost, client.RequestURL, internal.XmlPayload, r, reqOpts...)

	err = client.base.Send(ctx, internal.DisburseRequest, req, &response)

	if err != nil {
		return
	}

	return
}

