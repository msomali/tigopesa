package aw

import "encoding/xml"

type (
	Config struct {

	}

	DisburseRequest struct {
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

	DisburseResponse struct {
		XMLName     xml.Name `xml:"COMMAND" json:"-"`
		Text        string   `xml:",chardata" json:"-"`
		Type        string   `xml:"TYPE" json:"type"`
		ReferenceID string   `xml:"REFERENCEID" json:"reference_id"`
		TxnID       string   `xml:"TXNID" json:"txnid"`
		TxnStatus   string   `xml:"TXNSTATUS" json:"txn_status"`
		Message     string   `xml:"MESSAGE" json:"message"`
	}

)

