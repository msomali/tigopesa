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
	"context"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

type (
	TestPayload struct {
		Name string `json:"name,omitempty"`
		Age  int    `json:"age,omitempty"`
	}
)

func TestBaseClient_Send(t *testing.T) {
	var reqOpts []RequestOption
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	headersOpt := WithRequestHeaders(headers)
	reqOpts = append(reqOpts, headersOpt)
	reqPayload := TestPayload{
		Name: "tigopesa",
		Age:  10,
	}
	request := NewRequest(context.TODO(), http.MethodPost, "https://utopolo.com/piusalfred", JsonPayload, reqPayload, reqOpts...)
	response := TestPayload{}
	type fields struct {
		Http      *http.Client
		Ctx       context.Context
		Timeout   time.Duration
		Logger    io.Writer
		DebugMode bool
	}
	type args struct {
		ctx     context.Context
		rn      RequestName
		request *Request
		v       interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "sending to non existent url",
			fields: fields{
				Http:      http.DefaultClient,
				Ctx:       context.TODO(),
				Timeout:   60 * time.Second,
				Logger:    os.Stderr,
				DebugMode: true,
			},
			args: args{
				ctx:     context.TODO(),
				rn:      PaymentRequest,
				request: request,
				v:       &response,
			},
			wantErr: true,
		},
		{
			name: "sending to wrong address",
			fields: fields{
				Http:      http.DefaultClient,
				Ctx:       context.TODO(),
				Timeout:   60 * time.Second,
				Logger:    os.Stderr,
				DebugMode: true,
			},
			args: args{
				ctx:     context.TODO(),
				rn:      PaymentRequest,
				request: NewRequest(context.TODO(), http.MethodPost, "https://github.com/piusalfred", JsonPayload, reqPayload, reqOpts...),
				v:       &response,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &BaseClient{
				HTTP:      tt.fields.Http,
				Ctx:       tt.fields.Ctx,
				Timeout:   tt.fields.Timeout,
				Logger:    tt.fields.Logger,
				DebugMode: tt.fields.DebugMode,
			}
			if err := client.Send(tt.args.ctx, tt.args.rn, tt.args.request, tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
