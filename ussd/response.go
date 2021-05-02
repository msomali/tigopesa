package ussd

import "encoding/xml"

type SubscriberNameResponse struct {
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



type W2AResponse struct {
	XMLName          xml.Name `xml:"COMMAND"`
	Text             string   `xml:",chardata"`
	Type             string   `xml:"TYPE"`
	TxnID             string   `xml:"TXNID"`
	RefID            string   `xml:"REFID"`
	Result           string   `xml:"RESULT"`
	ErrorCode        string   `xml:"ERRORCODE"`
	ErrorDescription string   `xml:"ERRORDESCRIPTION"`
	Msisdn           string   `xml:"MSISDN"`
	Flag            string   `xml:"FLAG"`
	Content          string   `xml:"CONTENT"`
}

type A2WResponse struct {
	XMLName     xml.Name `xml:"COMMAND"`
	Text        string   `xml:",chardata"`
	Type        string   `xml:"TYPE"`
	ReferenceID string   `xml:"REFERENCEID"`
	TxnID                string   `xml:"TXNID"`
	TxnStatus string `xml:"TXNSTATUS"`
	Message string `xml:"MESSAGE"`
}
