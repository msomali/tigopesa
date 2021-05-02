package ussd

import (
	"context"
	"encoding/xml"
)

type NameCheckErrorCode string

const (
	SUCCESS        NameCheckErrorCode = "error000"
	NOT_REGISTERED NameCheckErrorCode = "error010"
	SUSPENDED      NameCheckErrorCode = "error020"
	INVALID_FORMAT NameCheckErrorCode = "error030"
)

type QuerySubscriberNameReq struct {
	XMLName             xml.Name `xml:"COMMAND"`
	Text                string   `xml:",chardata"`
	Type                string   `xml:"TYPE"`
	Msisdn              string   `xml:"MSISDN"`
	CompanyName         string   `xml:"COMPANYNAME"`
	CustomerReferenceID string   `xml:"CUSTOMERREFERENCEID"`
}

type QuerySubscriberNameResp struct {
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

type WalletToAccountReq struct {

}

type WalletToAccountResp struct {

}

type AccountToWalletReq struct {

}

type AccountToWalletResp struct {

}

type Service interface {

	// QuerySubscriberName is API to handle TigoPesa system Query
	// of the Customer’s Name from Partner (Third-party) – Synchronous protocol
	// req is the input from TigoPesa
	// <? xml version="1.0" encoding="UTF-8"?>
	//<COMMAND>
	//<TYPE>SYNC_LOOKUP_REQUEST</TYPE>
	//<MSISDN>User-msisdn<MSISDN>
	//<COMPANYNAME> copny_ID </COMPANYNAME>
	//<CUSTOMERREFERNCEID>Partner_refernce_number</CUSTOMERREFFERENCEID>
	//</COMMAND>
	QuerySubscriberName(ctx context.Context, req QuerySubscriberNameReq) (resp QuerySubscriberNameResp, err error)

	WalletToAccount(ctx context.Context, req WalletToAccountReq)(resp WalletToAccountResp, err error)

	AccountToWallet(ctx context.Context, req AccountToWalletReq)(resp AccountToWalletResp, err error)
}
