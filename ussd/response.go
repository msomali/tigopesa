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



type A2WResponse struct {
	XMLName     xml.Name `xml:"COMMAND"`
	Text        string   `xml:",chardata"`
	TYPE        string   `xml:"TYPE"`
	REFERENCEID string   `xml:"REFERENCEID"`
	TXNID                string   `xml:"TXNID"`
	TXNSTATUS string `xml:"TXNSTATUS"`
	MESSAGE string `xml:"MESSAGE"`
}
