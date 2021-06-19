package aw

import (
	"bytes"
	"context"
	"encoding/xml"
	"github.com/techcraftt/tigosdk/pkg/tigo"
	"io"
	"io/ioutil"
	"net/http"
)

var (
	_ DisburseHandler = (*DisburseHandlerFunc)(nil)
	_ DisburseHandler = (*Client)(nil)
)

type (
	DisburseHandler interface {
		Do(ctx context.Context, request DisburseRequest) (DisburseResponse, error)
	}

	DisburseHandlerFunc func(ctx context.Context, request DisburseRequest) (DisburseResponse, error)

	Config struct {
		AccountName   string
		AccountMSISDN string
		BrandID       string
		PIN           string
		RequestURL    string
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

	Client struct {
		*Config
		*tigo.BaseClient
	}
)

func (client *Client) Do(ctx context.Context, request DisburseRequest) (response DisburseResponse, err error) {
	// Marshal the request body into application/xml
	xmlStr, err := xml.MarshalIndent(request, "", "    ")
	if err != nil {
		return
	}
	xmlStr = []byte(xml.Header + string(xmlStr))

	// Create a HTTP Post Request to be sent to Tigo gateway
	req, err := http.NewRequest(http.MethodPost, client.RequestURL, bytes.NewBuffer(xmlStr)) // URL-encoded payload
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/xml")

	// Create context with a timeout of 1 minute
	ctx, cancel := context.WithTimeout(ctx, client.Timeout)
	defer cancel()

	req = req.WithContext(ctx)

	resp, err := client.HttpClient.Do(req)

	if err != nil {
		return
	}

	//todo: log here:
	go func(debugMode bool) {
		if debugMode{
			client.Log(req,resp)
		}
	}(client.DebugMode)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)
	xmlBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = xml.Unmarshal(xmlBody, &response)
	if err != nil {
		return
	}

	return
}

func (handler DisburseHandlerFunc) Do(ctx context.Context, request DisburseRequest) (DisburseResponse, error) {
	return handler(ctx, request)
}
