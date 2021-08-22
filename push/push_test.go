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

package push

import (
	"context"
	"github.com/techcraftlabs/tigopesa/internal"
	"github.com/techcraftlabs/tigopesa/term"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {

	config := &Config{
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
	var pushOpts []ClientOption
	debugOption := WithDebugMode(true)
	timeOutOption := WithTimeout(60 * time.Second)
	loggerOption := WithLogger(term.Stderr)
	contextOption := WithContext(context.TODO())
	httpOption := WithHTTPClient(http.DefaultClient)
	pushOpts = append(pushOpts, debugOption, loggerOption, timeOutOption, contextOption, httpOption)
	type args struct {
		config  *Config
		handler CallbackHandler
		opts    []ClientOption
	}
	tests := []struct {
		name string
		args args
		want *Client
	}{
		{
			name: "normal creation",
			args: args{
				config:  config,
				handler: nil,
				opts:    pushOpts,
			},
			want: &Client{
				Config: config,
				BaseClient: &internal.BaseClient{
					HTTP:      http.DefaultClient,
					Ctx:       context.TODO(),
					Timeout:   60 * time.Second,
					Logger:    term.Stderr,
					DebugMode: true,
				},
				CallbackHandler: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewClient(tt.args.config, tt.args.handler, tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClient() = %v, want %v", got, tt.want)
			}
		})
	}
}
