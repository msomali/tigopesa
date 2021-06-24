package internal

import (
	"bytes"
	"encoding/xml"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

const (
	xmlStr =
		`<COMMAND>
  <NAME>Tigo Pesa Developer</NAME>
  <AGE>20</AGE>
  <WISE>true</WISE>
</COMMAND>`

	jsonStr=`{"name":"Tigo Pesa Developer","age":20,"wise":true}`
)

type (
	request struct {
		XMLName xml.Name `json:"-" xml:"COMMAND"`
		Text    string   `json:"-" xml:",chardata"`
		Name    string   `json:"name" xml:"NAME"`
		Age     int      `json:"age" xml:"AGE"`
		Wise    bool     `json:"wise" xml:"WISE"`
	}
)

func TestMarshalPayload(t *testing.T) {
	data := url.Values{}
	data.Set("username", "user")
	data.Set("password", "password")
	data.Set("grant_type", "password")

	reader := strings.NewReader(data.Encode())
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(reader)

	req := request{
		Name: "Tigo Pesa Developer",
		Age:  20,
		Wise: true,
	}
	type args struct {
		payloadType PayloadType
		payload     interface{}
	}
	tests := []struct {
		name       string
		args       args
		wantBuffer *bytes.Buffer
		wantErr    bool
	}{
		{
			name: "test marshalling json",
			args: args{
				payloadType: JsonPayload,
				payload:     req,
			},
			wantBuffer: bytes.NewBufferString(jsonStr),
			wantErr:    false,
		},
		{
			name: "test marshalling xml",
			args: args{
				payloadType: XmlPayload,
				payload:     req,
			},
			wantBuffer: bytes.NewBufferString(xmlStr),
			wantErr:    false,
		},
		{
			name:       "marshal wrong form type",
			args:       args{
				payloadType: FormPayload,
				payload:     req,
			},
			wantBuffer: nil,
			wantErr:    true,
		},
		{
			name:       "marshal form",
			args:       args{
				payloadType: FormPayload,
				payload:     data,
			},
			wantBuffer: buf,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBuffer, err := MarshalPayload(tt.args.payloadType, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalPayload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotBuffer, tt.wantBuffer) {
				t.Errorf("MarshalPayload() gotBuffer = %v, want %v", gotBuffer, tt.wantBuffer)
			}
		})
	}
}
