package wa

import (
	"context"
	"encoding/xml"
	"github.com/techcraftt/tigosdk/internal"
	"github.com/techcraftt/tigosdk/pkg/tigo"
	"net/http"
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
		HandlePaymentRequest(ctx context.Context, request PayRequest) (PayResponse, error)
	}

	PaymentHandleFunc func(ctx context.Context, request PayRequest) (PayResponse, error)

	NameQueryHandler interface {
		HandleSubscriberNameQuery(ctx context.Context, request NameRequest) (NameResponse, error)
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
		*tigo.BaseClient
		*Config
		PaymentHandler   PaymentHandler
		NameQueryHandler NameQueryHandler
	}
)

func (handler PaymentHandleFunc) HandlePaymentRequest(ctx context.Context, request PayRequest) (PayResponse, error) {
	return handler(ctx, request)
}

func (handler NameQueryFunc) HandleSubscriberNameQuery(ctx context.Context, request NameRequest) (NameResponse, error) {
	return handler(ctx, request)
}

func (client *Client) HandleNameQuery(writer http.ResponseWriter, request *http.Request) {

	if client.DebugMode {
		client.Log("NAME QUERY", request)
	}
	ctx, cancel := context.WithTimeout(client.Ctx, client.Timeout)
	defer cancel()
	var req NameRequest

	err := tigo.Receive(request, internal.XmlPayload, &req)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

	response, err := client.NameQueryHandler.HandleSubscriberNameQuery(ctx, req)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

	resp := tigo.NewResponse(200, response, internal.XmlPayload)
	tigo.Reply(resp, writer)

}

func (client *Client) HandlePayment(writer http.ResponseWriter, request *http.Request) {

	if client.DebugMode {
		client.Log("USSD WA PAYMENT", request)
	}
	ctx, cancel := context.WithTimeout(client.Ctx, client.Timeout)
	defer cancel()

	var req PayRequest

	err := tigo.Receive(request, internal.XmlPayload, &req)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

	response, err := client.PaymentHandler.HandlePaymentRequest(ctx, req)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

	resp := tigo.NewResponse(200, response, internal.XmlPayload)

	tigo.Reply(resp, writer)

}
