package tigo

import (
	"context"
	"encoding/xml"
)
// QuerySubscriberReq a request body sent to the gateway to get user details
// it is in form of below xml
// <?xml version="1.0" encoding="UTF-8"?>
// <COMMAND><TYPE>QuerySubscriberReq</TYPE>
//     <REFERENCEID>Techcradt466559202007011</REFERENCEID>
//     <MSISDN>The MSISDN of partner disbursement</MSISDN>
//     <PIN>PIN of partner disbursement </PIN>
//     <MSISDN1>Subscriber account</MSISDN1>
//     <AMOUNT>Numeric only</AMOUNT>
//     <LANGUAGE1>en</LANGUAGE1>
//     <BRAND_ID></BRAND_ID>
// </COMMAND>
type QuerySubscriberReq struct {
	XMLName     xml.Name `xml:"COMMAND"`
	Text        string   `xml:",chardata"`
	TYPE        string   `xml:"TYPE"`
	REFERENCEID string   `xml:"REFERENCEID"`
	MSISDN      string   `xml:"MSISDN"`
	PIN         string   `xml:"PIN"`
	MSISDN1     string   `xml:"MSISDN1"`
	AMOUNT      string   `xml:"AMOUNT"`
	LANGUAGE1   string   `xml:"LANGUAGE1"`
	BRANDID     string   `xml:"BRAND_ID"`
}

// XMLer take a go struct and generate XML that can be sent over HTTP as a response
// also take the XML request and turn it into go struct
type XMLer interface {
	MakeXML() (err error)
	FromXML() (interface{}, error)
}

type NameReq struct {
	Type        string
	Txid        string
	Msisdn      string
	Amount      float64
	CompanyName string
	CustRefId   string
	SenderName  string
}

type NameResp struct {
}

type USSDService interface {
	Name(ctx context.Context, req NameReq) (resp NameResp, err error)
}
