package ussd

import "encoding/xml"

type SubscriberNameRequest struct {
	XMLName             xml.Name `xml:"COMMAND"`
	Text                string   `xml:",chardata"`
	Type                string   `xml:"TYPE"`
	Msisdn              string   `xml:"MSISDN"`
	CompanyName         string   `xml:"COMPANYNAME"`
	CustomerReferenceID string   `xml:"CUSTOMERREFERENCEID"`
}

// W2ARequest Input to the Partner Application from the Payment Gateway – SYNC_BILL
// PAY_API Request
// TYPE The Request Type of Transaction, the value will be constant in all request
// it will be SYNC_BILLPAY_REQUEST
// TXID Alphanumeric this filed will be Tigo Transaction ID, unique for all Transactions
// MSISDN Payer MSISDN it will be without country code
// AMOUNT amount in the transaction
// COMPANYNAME biller code or Business Number
// CUSTOMERREFERENCEID BillPay reference number normally generated/shared by partner application to validate payments
// SENDERNAME Name of Sender or Payer (50)
type W2ARequest struct {
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

type A2WRequest struct {
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
