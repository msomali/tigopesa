package internal_test

import (
	"context"
	"github.com/techcraftlabs/tigopesa/internal"
	"io"
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/techcraftlabs/tigopesa/disburse"
)

func TestBaseClient_LogPayload(t *testing.T) {
	type fields struct {
		HttpClient *http.Client
		Ctx        context.Context
		Timeout    time.Duration
		Logger     io.Writer
		DebugMode  bool
	}
	type args struct {
		t       PayloadType
		prefix  string
		payload interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Logging Normal Json Payload",
			fields: fields{
				HttpClient: http.DefaultClient,
				Ctx:        context.Background(),
				Timeout:    60 * time.Second,
				Logger:     os.Stderr,
				DebugMode:  true,
			},
			args: args{
				t:      JsonPayload,
				prefix: "disburse request",
				payload: disburse.Request{
					ReferenceID: "3673E67DDGVHSWHBJS89W89W",
					MSISDN:      "0765xxxxxxxx",
					Amount:      90000,
				},
			},
		},
		{
			name: "Logging Nil Payload",
			fields: fields{
				HttpClient: http.DefaultClient,
				Ctx:        context.Background(),
				Timeout:    60 * time.Second,
				Logger:     os.Stderr,
				DebugMode:  true,
			},
			args: args{
				t:       XmlPayload,
				prefix:  "nil payload",
				payload: nil,
			},
		},
		{
			name: "Normal XML Payload",
			fields: fields{
				HttpClient: http.DefaultClient,
				Ctx:        context.Background(),
				Timeout:    60 * time.Second,
				Logger:     os.Stdout,
				DebugMode:  true,
			},
			args: args{
				t:      XmlPayload,
				prefix: "xml payload",
				payload: disburse.Response{
					Type:        "RMFCI",
					ReferenceID: "ABAYU89282892JH2B",
					TxnID:       "NAJAIUIOWOWOWKS",
					TxnStatus:   "status of the transaction",
					Message:     "message",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &internal.BaseClient{
				Http:      tt.fields.HttpClient,
				Ctx:       tt.fields.Ctx,
				Timeout:   tt.fields.Timeout,
				Logger:    tt.fields.Logger,
				DebugMode: tt.fields.DebugMode,
			}

			client.logPayload(tt.args.t, tt.args.prefix, tt.args.payload)
		})
	}
}

func TestNewBaseClient(t *testing.T) {
	type args struct {
		opts []internal.ClientOption
	}
	tests := []struct {
		name string
		args args
		want *internal.BaseClient
	}{
		{
			name: "Normal Base Client",
			args: args{
				opts: []internal.ClientOption{
					internal.WithHTTPClient(http.DefaultClient),
					internal.WithLogger(os.Stderr),
					internal.WithTimeout(10 * time.Second),
					internal.WithDebugMode(true),
					internal.WithContext(context.TODO()),
				},
			},
			want: &internal.BaseClient{
				Http:      http.DefaultClient,
				Ctx:       context.TODO(),
				Timeout:   10 * time.Second,
				Logger:    os.Stderr,
				DebugMode: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := internal.NewBaseClient(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBaseClient() = %v, want %v", got, tt.want)
			}
		})
	}
}
