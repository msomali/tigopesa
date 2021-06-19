package wa

import (
	"context"
	"encoding/xml"
	"net/http"
)

type (

	Service interface {
		NameQueryHandler(writer http.ResponseWriter, request *http.Request)

		PaymentHandler(writer http.ResponseWriter, request *http.Request)
	}

	PayHandler interface {
		Do(ctx context.Context, request PayRequest) (PayResponse, error)
	}

	PayHandlerFunc func(ctx context.Context, request PayRequest) (PayResponse, error)

	NameHandler interface {
		Do(ctx context.Context, request NameRequest) (NameResponse, error)
	}

	NameHandlerFunc func(ctx context.Context, request NameRequest) (NameResponse, error)

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
)

const (

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
	_ PayHandler  = (*PayHandlerFunc)(nil)
	_ NameHandler = (*NameHandlerFunc)(nil)
)

func (handler PayHandlerFunc) Do(ctx context.Context, request PayRequest) (PayResponse, error) {
	return handler(ctx, request)
}

func (handler NameHandlerFunc) Do(ctx context.Context, request NameRequest) (NameResponse, error) {
	return handler(ctx, request)
}
