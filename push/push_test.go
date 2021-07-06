package push

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/techcraftt/tigosdk/pkg/tigo"
	"os"
	"testing"
	"time"
)

const (
	pushUsername   = "TIGO_PUSH_USERNAME"
	pushPassword   = "TIGO_PUSH_PASSWORD"
	pushPasswordGrantType = "TIGO_PUSH_PASSWORD_GRANT_TYPE"
	pushBaseUrl    = "TIGO_PUSH_BASE_URL"
	pushTokenUrl   = "TIGO_PUSH_TOKEN_URL"
	pushMSISDN     = "TIGO_PUSH_BILLER_MSISDN"
	pushBillerCode = "TIGO_PUSH_BILLER_CODE"
	pushPayURL     = "TIGO_PUSH_PAY_URL"
	pushReverseURL = "TIGO_PUSH_REVERSE_URL"
	pushHealthURL = "TIGO_PUSH_HEALTH_URL"
)

var PushConfig *Config

func TestMain(m *testing.M) {
	err := godotenv.Load("push.env")
	if err != nil {
		return
	}

	PushConfig = &Config{
		Username:              os.Getenv(pushUsername),
		Password:              os.Getenv(pushPassword),
		PasswordGrantType:     os.Getenv(pushPasswordGrantType),
		ApiBaseURL:            os.Getenv(pushBaseUrl),
		GetTokenURL:           os.Getenv(pushTokenUrl),
		BillerMSISDN:          os.Getenv(pushMSISDN),
		BillerCode:            os.Getenv(pushBillerCode),
		PushPayURL:            os.Getenv(pushPayURL),
		ReverseTransactionURL: os.Getenv(pushReverseURL),
		HealthCheckURL:        os.Getenv(pushHealthURL),
	}
	os.Exit(m.Run())
}

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

		{
			name: "test if token is returned",
			fields: fields{
				Config:     PushConfig,
				BaseClient: tigo.NewBaseClient(),
			},
			args: args{
				ctx: context.TODO(),
			},
			want:    "",
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
			if got != client.token {
				t.Errorf("Token() got = %v, want %v", got, tt.want)
			}

			if got == ""{
				t.Errorf("did not expect empty token to be returned\n")
			}

			t.Logf("\ntoken: %s\nclient.token %s\n", got, client.token)
		})
	}
}
