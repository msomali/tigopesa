/*
 * MIT License
 *
 * Copyright (c) 2021 TechCraft Technologies Co. Ltd
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 *
 */

package internal

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net/url"
)

const (
	JsonPayload PayloadType = iota
	XmlPayload
	FormPayload
)

var (
	//ErrInvalidFormPayload is returned when the PayloadType passed is
	// FormPayload but its not of type url.Values
	ErrInvalidFormPayload = errors.New("invalid form submitted: type url.Values is expected")
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
		if !ok {
			return nil, ErrInvalidFormPayload
		}

		buffer = bytes.NewBufferString(form.Encode())

		return buffer, nil

	default:
		err := fmt.Errorf("can not marshal the payload: invalid payload type")
		return nil, err
	}

}
