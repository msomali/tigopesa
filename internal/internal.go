package internal

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/url"
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
func MarshalPayload(payloadType PayloadType, payload interface{}) (buffer *bytes.Buffer, err error) {

	switch payloadType {
	case JsonPayload:
		buf, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}

		buffer = bytes.NewBuffer(buf)

		return buffer, nil
	case XmlPayload:
		buf, err := xml.MarshalIndent(payload, "", "  ")
		if err != nil {
			return nil, err
		}
		buffer = bytes.NewBuffer(buf)
		return buffer, nil


	case FormPayload:

		form, ok := payload.(url.Values)
		if !ok{
			err := fmt.Errorf("can not marshal the payload: invalid form has been submitted")
			return nil, err
		}

		buffer = bytes.NewBufferString(form.Encode())

		return buffer, nil

	default:
		err := fmt.Errorf("can not marshal the payload: invalid payload type")
		return nil, err
	}

}





