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
