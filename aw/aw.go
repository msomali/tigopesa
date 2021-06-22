package aw

import (
	"context"
	"encoding/xml"
	"github.com/techcraftt/tigosdk/internal"
	"github.com/techcraftt/tigosdk/pkg/tigo"
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
)

var (
	_ DisburseHandler = (*DisburseHandlerFunc)(nil)
	_ DisburseHandler = (*Client)(nil)
)

type (
	DisburseHandler interface {
		Disburse(ctx context.Context, request DisburseRequest) (DisburseResponse, error)
	}

	DisburseHandlerFunc func(ctx context.Context, request DisburseRequest) (DisburseResponse, error)

	Config struct {
		AccountName   string
		AccountMSISDN string
		BrandID       string
		PIN           string
		RequestURL    string
	}

	DisburseRequest struct {
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

	DisburseResponse struct {
		XMLName     xml.Name `xml:"COMMAND" json:"-"`
		Text        string   `xml:",chardata" json:"-"`
		Type        string   `xml:"TYPE" json:"type"`
		ReferenceID string   `xml:"REFERENCEID" json:"reference_id"`
		TxnID       string   `xml:"TXNID" json:"txnid"`
		TxnStatus   string   `xml:"TXNSTATUS" json:"txn_status"`
		Message     string   `xml:"MESSAGE" json:"message"`
	}

	Client struct {
		*Config
		*tigo.BaseClient
	}
)

func (client *Client) Disburse(ctx context.Context, request DisburseRequest) (response DisburseResponse, err error) {

	var reqOpts []tigo.RequestOption
	ctxOpt := tigo.WithRequestContext(ctx)
	headers := map[string]string{
		"Content-Type": "application/xml",
	}
	headersOpt := tigo.WithRequestHeaders(headers)
	reqOpts = append(reqOpts, ctxOpt, headersOpt)

	req := tigo.NewRequest(http.MethodPost, client.RequestURL, internal.XmlPayload, request, reqOpts...)

	err = client.Send(ctx, req, &response)

	if err != nil {
		return
	}

	return
}

func (handler DisburseHandlerFunc) Disburse(ctx context.Context, request DisburseRequest) (DisburseResponse, error) {
	return handler(ctx, request)
}
