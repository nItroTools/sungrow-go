package ws

import (
	"net"
	"reflect"
	"testing"
)

var (
	testIp   = net.ParseIP("192.168.2.100")
	testPort = 8082
	testPath = "/ws/home/overview"
)

func TestNewWS(t *testing.T) {
	type fields struct {
		ip   net.IP
		port int
		path string
	}

	tests := []struct {
		name    string
		fields  fields
		want    *WS
		wantErr error
	}{
		{
			name:    "1",
			fields:  fields{testIp, testPort, testPath},
			want:    &WS{ip: testIp, port: testPort, path: testPath},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewWS(tt.fields.ip, tt.fields.port, tt.fields.path)

			if reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("NewWS() = %v, want %v", reflect.TypeOf(got), reflect.TypeOf(tt.want))
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWS() = %v, want %v", got, tt.want)
			}
		})
	}
}
