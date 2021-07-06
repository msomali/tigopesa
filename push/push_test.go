package push

import (
	"context"
	"github.com/techcraftt/tigosdk/pkg/tigo"
	"testing"
	"time"
)

func TestClient_Token(t *testing.T) {
	type fields struct {
		Config          *Config
		BaseClient      *tigo.BaseClient
		CallbackHandler CallbackHandler
		token           string
		tokenExpires    time.Time
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "test token generation",
			fields:  fields{
				Config:          &Config{
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
				},
				BaseClient:      tigo.NewBaseClient(),
				CallbackHandler: nil,
				token:           "",
				tokenExpires:    time.Time{},
			},
			args:    args{
				ctx: context.TODO(),
			},
			want:    "this is token",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{
				Config:          tt.fields.Config,
				BaseClient:      tt.fields.BaseClient,
				CallbackHandler: tt.fields.CallbackHandler,
				token:           tt.fields.token,
				tokenExpires:    tt.fields.tokenExpires,
			}
			got, err := client.Token(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Token() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != client.token{
				t.Errorf("Token() got = %v, want %v", got, tt.want)
			}
		})
	}
}
