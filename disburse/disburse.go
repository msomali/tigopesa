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

var (
	_ Service = (*HandlerFunc)(nil)
	_ Service = (*Client)(nil)
)

type (
	Service interface {
		Disburse(ctx context.Context, referenceId, msisdn string, amount float64) (Response, error)
	}

	ClientX internal.BaseClient

	HandlerFunc func(ctx context.Context, referenceId, msisdn string, amount float64) (Response, error)

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
		*internal.BaseClient
	}
)

func NewClient(config *Config, opts ...ClientOption) *Client {

	client := &Client{
		Config: config,
	}
	client.Logger = defaultWriter
	client.Ctx = defaultCtx
	client.DebugMode = false
	client.Timeout = defaultTimeout
	client.Http = defaultHttpClient

	for _, opt := range opts {
		opt(client)
	}

	return client
}

func (client *Client) Disburse(ctx context.Context, referenceId, msisdn string, amount float64) (response Response, err error) {

	var reqOpts []internal.RequestOption
	ctxOpt := internal.WithRequestContext(ctx)
	headers := map[string]string{
		"Content-Type": "application/xml",
	}
	headersOpt := internal.WithRequestHeaders(headers)
	reqOpts = append(reqOpts, ctxOpt, headersOpt)

	request := disburseRequest{
		Type:        requestType,
		ReferenceID: referenceId,
		Msisdn:      client.Config.AccountMSISDN,
		PIN:         client.Config.PIN,
		Msisdn1:     msisdn,
		Amount:      amount,
		SenderName:  client.Config.AccountName,
		Language1:   senderLanguage,
		BrandID:     client.Config.BrandID,
	}

	req := internal.NewRequest(ctx, http.MethodPost, client.RequestURL, internal.XmlPayload, request, reqOpts...)

	err = client.Send(ctx, internal.DISBURSE_REQUEST, req, &response)

	if err != nil {
		return
	}

	return
}

func (handler HandlerFunc) Disburse(ctx context.Context, referenceId, msisdn string, amount float64) (Response, error) {
	return handler(ctx, referenceId, msisdn, amount)
}
