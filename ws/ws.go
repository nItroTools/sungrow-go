package ws

import (
	"crypto/tls"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/google/uuid"
	"golang.org/x/net/websocket"
)

type WS struct {
	ip       net.IP
	port     int
	user     string
	password string
	path     string
	conn     *websocket.Conn
	token    string
	uid      string
}

// NewWS returns a new WS instance.
func NewWS(p ConnectionParams) *WS {
	ws := &WS{
		ip:       p.Ip,
		port:     p.Port,
		user:     p.User,
		password: p.Password,
		path:     p.Path,
	}

	return ws
}

// Connect connects to the inverter using the WebSocket protocol.
func (ws *WS) Connect() (err error) {
	// Connect to WebSocket
	origin := fmt.Sprintf("https://%s", ws.ip.String())
	url := fmt.Sprintf("wss://%s:%d%s", ws.ip.String(), ws.port, ws.path)
	config, err := websocket.NewConfig(url, origin)
	if err != nil {
		return err
	}

	config.TlsConfig = &tls.Config{
		InsecureSkipVerify: true,
	}

	ws.conn, err = websocket.DialConfig(config)
	if err != nil {
		return err
	}

	// Connect to service
	ws.token = uuid.New().String()
	req := RequestConnect{"de_de", ws.token, "connect"}
	if err := websocket.JSON.Send(ws.conn, &req); err != nil {
		return err
	}
	res := ResponseConnect{}
	if err := websocket.JSON.Receive(ws.conn, &res); err != nil {
		return err
	}
	ws.token = res.ResultData.Token
	ws.uid = strconv.Itoa(res.ResultData.Uid)

	if res.ResultMsg != "success" {
		ws.Close()
		return fmt.Errorf("connected but connection request failed")
	}

	// Authenticate
	if ws.user != "" && ws.password != "" {
		reqAuth := RequestAuth{"de_de", ws.token, "login", ws.user, ws.password}
		if err := websocket.JSON.Send(ws.conn, &reqAuth); err != nil {
			return err
		}
		resAuth := ResponseAuth{}
		if err := websocket.JSON.Receive(ws.conn, &resAuth); err != nil {
			return err
		}

		if resAuth.ResultMsg != "success" {
			ws.Close()
			return fmt.Errorf("login request failed. Please check user/password params: %s", resAuth.ResultMsg)
		}

		ws.token = resAuth.ResultData.Token
	}

	return err
}

// Close closes the connection.
func (ws *WS) Close() {
	if ws.conn != nil {
		_ = ws.conn.Close()
	}
}

// Pv fetches pv data from the inverter.
func (ws *WS) Pv(keyList Keys, separator string) (err error) {
	return ws.fetch("real", keyList, separator)
}

// Battery fetches battery data from the inverter.
func (ws *WS) Battery(keyList Keys, separator string) (err error) {
	return ws.fetch("real_battery", keyList, separator)
}

// fetch fetches data from the inverter.
func (ws *WS) fetch(service string, keyList Keys, separator string) (err error) {
	req := RequestReal{"de_de", ws.token, ws.uid, service, time.Now().UnixMilli()}
	if err := websocket.JSON.Send(ws.conn, &req); err != nil {
		return err
	}
	resp := ResponseReal{}
	if err := websocket.JSON.Receive(ws.conn, &resp); err != nil {
		return err
	}

	// Output values
	for _, row := range resp.ResultData.List {
		if _, exists := keyList[row.DataName]; exists {
			val, _ := strconv.ParseFloat(row.DataValue, 64)
			fmt.Printf("%s%s%.3f%s%s\n", keyList[row.DataName], separator, val, separator, row.DataUnit)
		}
	}

	return nil
}
