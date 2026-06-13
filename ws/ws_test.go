package ws

import (
	"io"
	"net"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"golang.org/x/net/websocket"
)

var (
	testIp   = net.ParseIP("192.168.2.100")
	testPort = 443
	testPath = "/ws/home/overview"
)

func TestNewWS(t *testing.T) {
	tests := []struct {
		name    string
		params  ConnectionParams
		want    *WS
		wantErr error
	}{
		{
			name:    "1",
			params:  ConnectionParams{Ip: testIp, Port: testPort, Path: testPath},
			want:    &WS{ip: testIp, port: testPort, path: testPath},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewWS(tt.params)

			if reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("NewWS() = %v, want %v", reflect.TypeOf(got), reflect.TypeOf(tt.want))
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWS() = %v, want %v", got, tt.want)
			}
		})
	}
}

// realRow is a single telemetry row served by the fake inverter.
type realRow struct{ name, value, unit string }

// serverConfig configures the fake inverter's responses.
type serverConfig struct {
	connectMsg string // result_msg for the "connect" service; defaults to "success"
	loginMsg   string // result_msg for the "login" service; defaults to "success"
	rows       []realRow
}

// newWSServer starts an httptest TLS server that speaks the Sungrow WebSocket
// protocol: it answers connect/login/real/real_battery requests synchronously
// on a single long-lived connection, mirroring the real dongle.
func newWSServer(t *testing.T, cfg serverConfig) *httptest.Server {
	t.Helper()
	if cfg.connectMsg == "" {
		cfg.connectMsg = "success"
	}
	if cfg.loginMsg == "" {
		cfg.loginMsg = "success"
	}

	handler := websocket.Handler(func(c *websocket.Conn) {
		for {
			var req map[string]interface{}
			if err := websocket.JSON.Receive(c, &req); err != nil {
				return // connection closed by client
			}
			switch req["service"] {
			case "connect":
				_ = websocket.JSON.Send(c, map[string]interface{}{
					"result_code": 1,
					"result_msg":  cfg.connectMsg,
					"result_data": map[string]interface{}{"token": "connect-token", "uid": 42},
				})
			case "login":
				_ = websocket.JSON.Send(c, map[string]interface{}{
					"result_code": 1,
					"result_msg":  cfg.loginMsg,
					"result_data": map[string]interface{}{"token": "login-token", "uid": 42},
				})
			case "real", "real_battery":
				list := make([]map[string]interface{}, 0, len(cfg.rows))
				for _, r := range cfg.rows {
					list = append(list, map[string]interface{}{
						"data_name":  r.name,
						"data_value": r.value,
						"data_unit":  r.unit,
					})
				}
				_ = websocket.JSON.Send(c, map[string]interface{}{
					"result_code": 1,
					"result_msg":  "success",
					"result_data": map[string]interface{}{"list": list, "count": len(list)},
				})
			}
		}
	})

	return httptest.NewTLSServer(handler)
}

// newTestWS builds a WS pointing at the given test server.
func newTestWS(t *testing.T, server *httptest.Server, user, password string) *WS {
	t.Helper()
	u, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("parse server url: %v", err)
	}
	host, portStr, err := net.SplitHostPort(u.Host)
	if err != nil {
		t.Fatalf("split host/port: %v", err)
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		t.Fatalf("parse port: %v", err)
	}
	return NewWS(ConnectionParams{
		Ip:       net.ParseIP(host),
		Port:     port,
		Path:     testPath,
		User:     user,
		Password: password,
	})
}

// captureStdout runs fn while capturing everything written to os.Stdout.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stdout = w
	fn()
	_ = w.Close()
	os.Stdout = old
	out, _ := io.ReadAll(r)
	return string(out)
}

func TestConnectAnonymous(t *testing.T) {
	srv := newWSServer(t, serverConfig{})
	defer srv.Close()

	ws := newTestWS(t, srv, "", "")
	if err := ws.Connect(); err != nil {
		t.Fatalf("Connect() error = %v", err)
	}
	defer ws.Close()

	// No credentials -> token stays the one from the connect handshake.
	if ws.token != "connect-token" {
		t.Errorf("token = %q, want %q", ws.token, "connect-token")
	}
	if ws.uid != "42" {
		t.Errorf("uid = %q, want %q", ws.uid, "42")
	}
}

func TestConnectWithLogin(t *testing.T) {
	srv := newWSServer(t, serverConfig{})
	defer srv.Close()

	ws := newTestWS(t, srv, "admin", "secret")
	if err := ws.Connect(); err != nil {
		t.Fatalf("Connect() error = %v", err)
	}
	defer ws.Close()

	// Credentials present -> login replaces the token.
	if ws.token != "login-token" {
		t.Errorf("token = %q, want %q", ws.token, "login-token")
	}
}

func TestConnectHandshakeFails(t *testing.T) {
	srv := newWSServer(t, serverConfig{connectMsg: "error"})
	defer srv.Close()

	ws := newTestWS(t, srv, "", "")
	err := ws.Connect()
	if err == nil {
		t.Fatal("Connect() error = nil, want connection failure")
	}
	if !strings.Contains(err.Error(), "connection request failed") {
		t.Errorf("Connect() error = %q, want it to mention connection failure", err)
	}
}

func TestConnectLoginFails(t *testing.T) {
	srv := newWSServer(t, serverConfig{loginMsg: "wrong password"})
	defer srv.Close()

	ws := newTestWS(t, srv, "admin", "wrong")
	err := ws.Connect()
	if err == nil {
		t.Fatal("Connect() error = nil, want login failure")
	}
	if !strings.Contains(err.Error(), "login request failed") {
		t.Errorf("Connect() error = %q, want it to mention login failure", err)
	}
}

func TestConnectDialError(t *testing.T) {
	srv := newWSServer(t, serverConfig{})
	ws := newTestWS(t, srv, "", "")
	srv.Close() // free the port so the dial is refused

	if err := ws.Connect(); err == nil {
		t.Fatal("Connect() error = nil, want dial error")
	}
}

func TestPvOutput(t *testing.T) {
	srv := newWSServer(t, serverConfig{rows: []realRow{
		{"I18N_KNOWN", "1.5", "kW"},
		{"I18N_UNKNOWN", "9.9", "x"}, // not in keyList -> dropped
		{"I18N_NAN", "abc", "%"},     // unparseable -> 0.000
	}})
	defer srv.Close()

	ws := newTestWS(t, srv, "", "")
	if err := ws.Connect(); err != nil {
		t.Fatalf("Connect() error = %v", err)
	}
	defer ws.Close()

	keys := Keys{"I18N_KNOWN": "known", "I18N_NAN": "nan"}
	out := captureStdout(t, func() {
		if err := ws.Pv(keys, ","); err != nil {
			t.Fatalf("Pv() error = %v", err)
		}
	})

	want := "known,1.500,kW\nnan,0.000,%\n"
	if out != want {
		t.Errorf("Pv() output = %q, want %q", out, want)
	}
}

func TestBatteryOutputSeparator(t *testing.T) {
	srv := newWSServer(t, serverConfig{rows: []realRow{
		{"I18N_COMMON_BATTERY_SOC", "87", "%"},
	}})
	defer srv.Close()

	ws := newTestWS(t, srv, "", "")
	if err := ws.Connect(); err != nil {
		t.Fatalf("Connect() error = %v", err)
	}
	defer ws.Close()

	keys := Keys{"I18N_COMMON_BATTERY_SOC": "batteryLevel"}
	out := captureStdout(t, func() {
		if err := ws.Battery(keys, ";"); err != nil {
			t.Fatalf("Battery() error = %v", err)
		}
	})

	want := "batteryLevel;87.000;%\n"
	if out != want {
		t.Errorf("Battery() output = %q, want %q", out, want)
	}
}

func TestClose(t *testing.T) {
	// Close on a fresh WS (nil conn) must not panic.
	ws := NewWS(ConnectionParams{Ip: testIp, Port: testPort, Path: testPath})
	ws.Close()

	// Close after a real connection should be safe and idempotent.
	srv := newWSServer(t, serverConfig{})
	defer srv.Close()
	connected := newTestWS(t, srv, "", "")
	if err := connected.Connect(); err != nil {
		t.Fatalf("Connect() error = %v", err)
	}
	connected.Close()
	connected.Close() // idempotent
}
