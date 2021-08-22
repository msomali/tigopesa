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

package push_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/techcraftlabs/tigopesa/push"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var _ push.CallbackHandler = (*testHandler)(nil)

type (
	testHandler int
)

func (t testHandler) Respond(ctx context.Context, request push.CallbackRequest) (push.CallbackResponse, error) {
	_, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	return push.CallbackResponse{
		ResponseCode:        push.SuccessCode,
		ResponseDescription: request.Description,
		ResponseStatus:      request.Status,
		ReferenceID:         request.ReferenceID,
	}, nil
}

func TestHealthCheckHandler(t *testing.T) {
	conf := &push.Config{
		Username:              "",
		Password:              "",
		PasswordGrantType:     "",
		ApiBaseURL:            "",
		GetTokenURL:           "",
		BillerMSISDN:          "",
		BillerCode:            "",
		PushPayURL:            "",
		ReverseTransactionURL: "",
		HealthCheckURL:        "",
	}
	cHandler := testHandler(1)
	client := push.NewClient(conf, cHandler, push.WithDebugMode(true))
	reqPayload := push.CallbackRequest{
		Status:           true,
		Description:      "this is test",
		MFSTransactionID: "TWTSVBSVBSGFSYA",
		ReferenceID:      "WWTYTYW6W67WTW",
		Amount:           "10000",
	}

	buf, _ := json.Marshal(reqPayload)
	req, err := http.NewRequestWithContext(context.TODO(), http.MethodPost, "/tigopesa/callback", bytes.NewBuffer(buf))
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(client.Callback)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"ResponseCode":"BILLER-30-0000-S","ResponseDescription":"this is test","ResponseStatus":true,"ReferenceID":"WWTYTYW6W67WTW"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
