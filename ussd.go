package tigo

import "context"

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
