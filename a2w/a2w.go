package a2w

import (
	"context"
	"encoding/xml"
)

type Request struct {
	ReferenceID string   `xml:"REFERENCEID"`
	Msisdn      string   `xml:"MSISDN"`
	PIN         string   `xml:"PIN"`
	Msisdn1     string   `xml:"MSISDN1"`
	Amount      float64  `xml:"AMOUNT"`
	SenderName  string   `xml:"SENDERNAME"`
	Language1   string   `xml:"LANGUAGE1"`
	BrandID     string   `xml:"BRAND_ID"`
}

type Response struct {
	ReferenceID string   `xml:"REFERENCEID"`
	TxnID       string   `xml:"TXNID"`
	TxnStatus   string   `xml:"TXNSTATUS"`
	Message     string   `xml:"MESSAGE"`
}

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

type XMLResponse struct {
	XMLName     xml.Name `xml:"COMMAND"`
	Text        string   `xml:",chardata"`
	Type        string   `xml:"TYPE"`
	ReferenceID string   `xml:"REFERENCEID"`
	TxnID       string   `xml:"TXNID"`
	TxnStatus   string   `xml:"TXNSTATUS"`
	Message     string   `xml:"MESSAGE"`
}

type HandleFunc func(context.Context, Request) (Response, error)

func (f HandleFunc) Handle(ctx context.Context, request Request) (response Response, err error) {
	return f(ctx, request)
}

type Handler interface {
	Handle(ctx context.Context, request Request) (response Response, err error)
}

type Service interface {
	AccountToWallet(ctx context.Context, req XMLRequest) (resp XMLResponse, err error)
}
