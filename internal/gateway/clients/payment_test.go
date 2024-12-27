package clients

import (
	"bytes"
	"io"
	"net/http"
	"testing"
)

type httpClientStub struct {
	err        error
	statusCode int
}

func (h *httpClientStub) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{StatusCode: h.statusCode, Body: io.NopCloser(bytes.NewBufferString("test"))}, h.err
}

func TestPaymentClient_Cancel(t *testing.T) {
	type fields struct {
		client  httpClient
		baseUrl string
	}
	type args struct {
		uuid string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				client: &httpClientStub{
					err:        nil,
					statusCode: 204,
				},
				baseUrl: "http://localhost/payments/cancel",
			},
			args: args{
				uuid: "test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PaymentClient{
				client:  tt.fields.client,
				baseUrl: tt.fields.baseUrl,
			}
			if err := p.Cancel(tt.args.uuid); (err != nil) != tt.wantErr {
				t.Errorf("Cancel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
