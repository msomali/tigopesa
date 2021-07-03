package tigo_test

import (
	"context"
	"io"
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/techcraftt/tigosdk/aw"
	"github.com/techcraftt/tigosdk/internal"
	"github.com/techcraftt/tigosdk/pkg/tigo"
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
		t       internal.PayloadType
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
				t:      internal.JsonPayload,
				prefix: "disburse request",
				payload: aw.DisburseRequest{
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
				t:       internal.XmlPayload,
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
				t:      internal.XmlPayload,
				prefix: "xml payload",
				payload: aw.DisburseResponse{
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
			client := &tigo.BaseClient{
				HttpClient: tt.fields.HttpClient,
				Ctx:        tt.fields.Ctx,
				Timeout:    tt.fields.Timeout,
				Logger:     tt.fields.Logger,
				DebugMode:  tt.fields.DebugMode,
			}

			client.LogPayload(tt.args.t, tt.args.prefix, tt.args.payload)
		})
	}
}

func TestNewBaseClient(t *testing.T) {
	type args struct {
		opts []tigo.ClientOption
	}
	tests := []struct {
		name string
		args args
		want *tigo.BaseClient
	}{
		{
			name: "Normal Base Client",
			args: args{
				opts: []tigo.ClientOption{
					tigo.WithHTTPClient(http.DefaultClient),
					tigo.WithLogger(os.Stderr),
					tigo.WithTimeout(10 * time.Second),
					tigo.WithDebugMode(true),
					tigo.WithContext(context.TODO()),
				},
			},
			want: &tigo.BaseClient{
				HttpClient: http.DefaultClient,
				Ctx:        context.TODO(),
				Timeout:    10 * time.Second,
				Logger:     os.Stderr,
				DebugMode:  true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tigo.NewBaseClient(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBaseClient() = %v, want %v", got, tt.want)
			}
		})
	}
}
