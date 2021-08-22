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
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/techcraftlabs/tigopesa/internal/term"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

const (
	jsonContentTypeString   = "application/json"
	xmlContentTypeString    = "text/xml"
	appXMLContentTypeString = "application/xml"
)

type (
	BaseClient struct {
		HTTP      *http.Client
		Logger    io.Writer // for logging purposes
		DebugMode bool
	}
)

func NewBaseClient(opts ...ClientOption) *BaseClient {

	cl := &http.Client{
		Timeout: 60 * time.Second,
	}
	client := &BaseClient{
		HTTP:      cl,
		Logger:    term.Stderr,
		DebugMode: false,
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

func (client *BaseClient) logPayload(t PayloadType, prefix string, payload interface{}) {
	buf, _ := MarshalPayload(t, payload)
	_, _ = client.Logger.Write([]byte(fmt.Sprintf("%s: %s\n\n", prefix, buf.String())))
}

// log is called to print the details of http.Request sent from Tigo during
// callback, namecheck or ussd payment. It is used for debugging purposes.
func (client *BaseClient) log(name string, request *http.Request) {
	if request != nil {
		reqDump, _ := httputil.DumpRequest(request, true)
		_, err := fmt.Fprintf(client.Logger, "%s REQUEST: %s\n", name, reqDump)
		if err != nil {
			log.Printf("error while logging %s request: %v\n",
				strings.ToLower(name), err)

			return
		}

		return
	}
}

// logOut is like log except this is for outgoing client requests:
// http.Request that is supposed to be sent to tigo
func (client *BaseClient) logOut(name string, request *http.Request, response *http.Response) {
	if request != nil {
		reqDump, _ := httputil.DumpRequestOut(request, true)
		_, err := fmt.Fprintf(client.Logger, "%s REQUEST: %s\n", name, reqDump)
		if err != nil {
			log.Printf("error while logging %s request: %v\n",
				strings.ToLower(name), err)
		}
	}

	if response != nil {
		respDump, _ := httputil.DumpResponse(response, true)
		_, err := fmt.Fprintf(client.Logger, "%s RESPONSE: %s\n", name, respDump)
		if err != nil {
			log.Printf("error while logging %s response: %v\n",
				strings.ToLower(name), err)
		}
	}

	return
}

func (client *BaseClient) Send(ctx context.Context, rn RequestName, request *Request, v interface{}) error {
	var req *http.Request
	var res *http.Response

	var reqBodyBytes []byte
	var resBodyBytes []byte
	defer func(debug bool) {
		if debug {
			req.Body = io.NopCloser(bytes.NewBuffer(reqBodyBytes))
			name := strings.ToUpper(rn.String())
			if res == nil {
				client.logOut(name, req, nil)
				return
			}
			res.Body = io.NopCloser(bytes.NewBuffer(resBodyBytes))
			client.logOut(name, req, res)

		}
	}(client.DebugMode)

	//creates http request with context
	req, err := request.newRequestWithContext()

	if err != nil {
		return err
	}

	if req.Body != nil {
		reqBodyBytes, _ = io.ReadAll(req.Body)
	}

	if v == nil {
		return errors.New("v interface can not be empty")
	}

	req.Body = io.NopCloser(bytes.NewBuffer(reqBodyBytes))
	res, err = client.HTTP.Do(req)

	if err != nil {
		return err
	}

	if res.Body != nil {
		resBodyBytes, _ = io.ReadAll(res.Body)
	}

	contentType := res.Header.Get("Content-Type")
	if strings.Contains(contentType, jsonContentTypeString) {
		if err := json.NewDecoder(bytes.NewBuffer(resBodyBytes)).Decode(v); err != nil {
			if err != io.EOF {
				return err
			}
		}
	}

	if strings.Contains(contentType, xmlContentTypeString) ||
		strings.Contains(contentType, appXMLContentTypeString) {
		if err := xml.NewDecoder(bytes.NewBuffer(resBodyBytes)).Decode(v); err != nil {
			if err != io.EOF {
				return err
			}
		}
	}
	return nil
}
