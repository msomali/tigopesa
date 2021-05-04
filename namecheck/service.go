package namecheck

import (
	"context"
	"encoding/xml"
)

type Request struct {
	Msisdn              string   `xml:"MSISDN"`
	CompanyName         string   `xml:"COMPANYNAME"`
	CustomerReferenceID string   `xml:"CUSTOMERREFERENCEID"`
}

type Response struct {
	Result    string   `xml:"RESULT"`
	ErrorCode string   `xml:"ERRORCODE"`
	ErrorDesc string   `xml:"ERRORDESC"`
	Msisdn    string   `xml:"MSISDN"`
	Flag      string   `xml:"FLAG"`
	Content   string   `xml:"CONTENT"`
}

type XMLRequest struct {
	XMLName             xml.Name `xml:"COMMAND"`
	Text                string   `xml:",chardata"`
	Type                string   `xml:"TYPE"`
	Msisdn              string   `xml:"MSISDN"`
	CompanyName         string   `xml:"COMPANYNAME"`
	CustomerReferenceID string   `xml:"CUSTOMERREFERENCEID"`
}

type XMLResponse struct {
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

type HandleFunc func(context.Context, Request) (Response, error)

func (f HandleFunc) Handle(ctx context.Context, request Request) (response Response, err error) {
	return f(ctx, request)
}

type Handler interface {
	Handle(ctx context.Context, request Request) (response Response, err error)
}

type Service interface {

	// QuerySubscriberName is API to handle TigoPesa system Query
	// of the Customer’s Name from Partner (Third-party) – Synchronous protocol
	// req is the input from TigoPesa
	QuerySubscriberName(ctx context.Context, req XMLRequest) (resp XMLResponse, err error)

}
