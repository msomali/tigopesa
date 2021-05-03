package w2a

import (
	"context"
	"encoding/xml"
)

//type  ErrorCodeKey string
//
//const (
//	SUCCESSFUL_TRANSACTION       ErrorCodeKey = "error000"
//	SERVICE_NOT_AVAILABLE        ErrorCodeKey = "error001"
//	INVALID_CUSTOMER_REF_NUMBER  ErrorCodeKey = "error010"
//	CUSTOMER_REF_NUM_LOCKED      ErrorCodeKey = "error011"
//	INVALID_AMOUNT               ErrorCodeKey = "error012"
//	AMOUNT_INSUFFICIENT          ErrorCodeKey = "error013"
//	AMOUNT_TOO_HIGH              ErrorCodeKey = "error014"
//	AMOUNT_TOO_LOW               ErrorCodeKey = "error015"
//	INVALID_PAYMENT              ErrorCodeKey = "error016"
//	GENERAL_ERROR                ErrorCodeKey = "error100"
//	RETRY_CONDITION_NO_RESPONSE  ErrorCodeKey = "error111"
//)
//
//var ErrorCode = map[ErrorCodeKey]string{
//	SUCCESSFUL_TRANSACTION:"Successful Transaction",
//	SERVICE_NOT_AVAILABLE: "Service Not Available",
//
//}

type XMLRequest struct {
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

type Request struct {
	TxnID               string
	Msisdn              string
	Amount              float64
	CompanyName         string
	CustomerReferenceID string
	SenderName          string
}

type XMLResponse struct {
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

type Response struct {
	TxnID            string
	RefID            string
	Result           string
	ErrorCode        string
	ErrorDescription string
	Msisdn           string
	Flag             string
	Content          string
}

type HandleFunc func(context.Context, Request) (Response, error)

func (f HandleFunc) Handle(ctx context.Context, request Request) (response Response, err error) {
	return f(ctx, request)
}

type Handler interface {
	Handle(ctx context.Context, request Request) (response Response, err error)
}

type Service interface {
	WalletToAccount(ctx context.Context, req XMLRequest)(resp XMLResponse, err error)
}
