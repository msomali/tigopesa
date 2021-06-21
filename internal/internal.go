package internal

import (
	"encoding/json"
	"encoding/xml"
)

const (
	JsonPayload PayloadType = iota
	XmlPayload
	FormPayload
)

type (
	PayloadType int
)

// MarshalPayload returns the JSON/XML encoding of payload.
func MarshalPayload(payloadType PayloadType, payload interface{}) (buf []byte, err error) {
	switch payloadType {
	case JsonPayload:
		buf, err = json.Marshal(payload)
		if err != nil {
			return nil, err
		}
	case XmlPayload:
		buf, err = xml.MarshalIndent(payload, "", "  ")
		if err != nil {
			return nil, err
		}
	}
	return
}



