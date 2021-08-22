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
