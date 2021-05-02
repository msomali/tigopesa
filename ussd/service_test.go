package ussd

import (
	"encoding/xml"
	"fmt"
	"testing"
)

func TestUnMarshalReq(t *testing.T) {
	req :=
		`<?xml version="1.0" encoding="UTF-8"?>
<COMMAND>
<TYPE>SYNC_LOOKUP_REQUEST</TYPE>
<MSISDN>User-msisdn</MSISDN>
<COMPANYNAME> copny_ID </COMPANYNAME>
<CUSTOMERREFFERENCEID>Partner_refernce_number</CUSTOMERREFFERENCEID>
</COMMAND>
`
	var queryReq SubscriberNameRequest

	queryResp := SubscriberNameResponse{
		Type:      "yehes",
		Result:    "hshsnsn",
		ErrorCode: "msnsmsm",
		ErrorDesc: "msmsm",
		Msisdn:    "msmsms",
		Flag:      "msmsms",
		Content:   "msmsmsm",
	}

	t.Run("testing request unmarshalling", func(t *testing.T) {
		err := xml.Unmarshal([]byte(req), &queryReq)

		if err != nil {
			t.Error(err)
		}

		t.Logf("type: %s, msisdn: %s, comapny: %s", queryReq.Type, queryReq.Msisdn, queryReq.CompanyName)

	})

	t.Run("response marshalling", func(t *testing.T) {

		xmlstring, err := xml.MarshalIndent(queryResp, "", "    ")
		if err != nil {
			t.Error(err)
		}
		xmlstring = []byte(xml.Header + string(xmlstring))
		fmt.Printf("%s\n", xmlstring)
	})

}
