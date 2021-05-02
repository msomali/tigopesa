package ussd

import (
	"context"
	"encoding/xml"
)

//type TxnStatus map[string]string
//
//var (
//	Statuses = TxnStatus{
//		"200":"Success",
//		"0":"Success",
//		"410": "unable to complete transaction, amount more than maximum limit",
//	}
//)

type NameCheckErrorCode string

const (
	SUCCESS        NameCheckErrorCode = "error000"
	NOT_REGISTERED NameCheckErrorCode = "error010"
	SUSPENDED      NameCheckErrorCode = "error020"
	INVALID_FORMAT NameCheckErrorCode = "error030"
)

type W2AErrorCode string

const (
	SUCCESSFUL_TRANSACTION      W2AErrorCode = "error000"
	SERVICE_NOT_AVAILABLE       W2AErrorCode = "error001"
	INVALID_CUSTOMER_REF_NUMBER W2AErrorCode = "error010"
	CUSTOMER_REF_NUM_LOCKED     W2AErrorCode = "error011"
	INVALID_AMOUNT              W2AErrorCode = "error012"
	AMOUNT_INSUFFICIENT         W2AErrorCode = "error013"
	AMOUNT_TOO_HIGH             W2AErrorCode = "error014"
	AMOUNT_TOO_LOW              W2AErrorCode = "error015"
	INVALID_PAYMENT             W2AErrorCode = "error016"
	GENERAL_ERROR               W2AErrorCode = "error100"
	RETRY_CONDITION_NO_RESPONSE W2AErrorCode = "error111"
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

// WalletToAccountReq Input to the Partner Application from the Payment Gateway – SYNC_BILL
// PAY_API Request
// TYPE The Request Type of Transaction, the value will be constant in all request
// it will be SYNC_BILLPAY_REQUEST
// TXID Alphanumeric this filed will be Tigo Transaction ID, unique for all Transactions
// MSISDN Payer MSISDN it will be without country code
// AMOUNT amount in the transaction
// COMPANYNAME biller code or Business Number
// CUSTOMERREFERENCEID BillPay reference number normally generated/shared by partner application to validate payments
// SENDERNAME Name of Sender or Payer (50)
type WalletToAccountReq struct {
	XMLName             xml.Name `xml:"COMMAND"`
	Text                string   `xml:",chardata"`
	TYPE                string   `xml:"TYPE"`
	TXID                string   `xml:"TXID"`
	MSISDN              string   `xml:"MSISDN"`
	AMOUNT              string   `xml:"AMOUNT"`
	COMPANYNAME         string   `xml:"COMPANYNAME"`
	CUSTOMERREFERENCEID string   `xml:"CUSTOMERREFERENCEID"`
	SENDERNAME          string   `xml:"SENDERNAME"`
}

type WalletToAccountResp struct {
	XMLName          xml.Name `xml:"COMMAND"`
	Text             string   `xml:",chardata"`
	TYPE             string   `xml:"TYPE"`
	TXID             string   `xml:"TXID"`
	REFID            string   `xml:"REFID"`
	RESULT           string   `xml:"RESULT"`
	ERRORCODE        string   `xml:"ERRORCODE"`
	ERRORDESCRIPTION string   `xml:"ERRORDESCRIPTION"`
	MSISDN           string   `xml:"MSISDN"`
	FLAG             string   `xml:"FLAG"`
	CONTENT          string   `xml:"CONTENT"`
}

type AccountToWalletReq struct {
	XMLName     xml.Name `xml:"COMMAND"`
	Text        string   `xml:",chardata"`
	TYPE        string   `xml:"TYPE"`
	REFERENCEID string   `xml:"REFERENCEID"`
	MSISDN      string   `xml:"MSISDN"`
	PIN         string   `xml:"PIN"`
	MSISDN1     string   `xml:"MSISDN1"`
	AMOUNT      string   `xml:"AMOUNT"`
	SENDERNAME  string   `xml:"SENDERNAME"`
	LANGUAGE1   string   `xml:"LANGUAGE1"`
	BRANDID     string   `xml:"BRAND_ID"`
}

type AccountToWalletResp struct {
	XMLName     xml.Name `xml:"COMMAND"`
	Text        string   `xml:",chardata"`
	TYPE        string   `xml:"TYPE"`
	REFERENCEID string   `xml:"REFERENCEID"`
	TXNID                string   `xml:"TXNID"`
	TXNSTATUS string `xml:"TXNSTATUS"`
	MESSAGE string `xml:"MESSAGE"`
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

	WalletToAccount(ctx context.Context, req WalletToAccountReq) (resp WalletToAccountResp, err error)

	AccountToWallet(ctx context.Context, req AccountToWalletReq) (resp AccountToWalletResp, err error)
}
