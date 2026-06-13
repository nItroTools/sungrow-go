package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/nItroTools/sungrow-go/ws"
)

type inverter struct {
	ws        *ws.WS
	ipS       string
	ip        net.IP
	port      int
	user      string
	password  string
	path      string
	data      string
	separator string
	types     []string
}

var inv *inverter

func main() {
	inv = &inverter{}

	// Flags
	flags()

	// Connect to inverter
	cP := ws.ConnectionParams{
		Ip:       inv.ip,
		Port:     inv.port,
		User:     inv.user,
		Password: inv.password,
		Path:     inv.path,
	}
	inv.ws = ws.NewWS(cP)
	if err := inv.ws.Connect(); err != nil {
		log.Fatalln(err)
	}
	defer inv.ws.Close()

	// Output timestamp row
	fmt.Printf("%s%s%s%s%s\n", "time", inv.separator, time.Now().Format(time.RFC3339), inv.separator, "RFC3339")

	// Fetch values from inverter
	for _, t := range inv.types {
		switch t {
		case "pv":
			_ = inv.ws.Pv(pvKeys, inv.separator)
			break
		case "battery":
			_ = inv.ws.Battery(batteryKeys, inv.separator)
			break
		}
	}
}

// flags define, parse and validate command-line flags from os.Args[1:]
func flags() {
	ipS := flag.String("ip", "", "Required: IP address of the Sungrow inverter")
	port := flag.Int("port", 443, "Secure WebSocket port of the Sungrow inverter")
	user := flag.String("user", "", "Required: username for the Sungrow inverter web ui login, e.g. admin")
	password := flag.String("password", "", "Required: password for the Sungrow inverter web ui login")
	path := flag.String("path", "/ws/home/overview", "Server path from where data is requested")
	data := flag.String("data", "pv,battery", "Select the data to be requested comma separated.\nPossible values are \"pv\" and \"battery\"")
	separator := flag.String("separator", ",", "Output data separator")
	flag.Parse()

	inv.ipS = *ipS
	inv.port = *port
	inv.user = *user
	inv.password = *password
	inv.path = *path
	inv.data = *data
	inv.separator = *separator

	// Validate flags
	flagsValidate()
}

// flagsValidate validates all flags
func flagsValidate() {
	if inv.ip = net.ParseIP(inv.ipS); inv.ip == nil {
		log.Fatalln("Required parameter 'ip' not set or invalid ip address!\n'sungrow-go -help' lists available parameters.")
	}

	if strings.TrimSpace(inv.user) == "" {
		log.Fatalln("Required parameter 'user' not set!\n'sungrow-go -help' lists available parameters.")
	}

	if strings.TrimSpace(inv.password) == "" {
		log.Fatalln("Required parameter 'password' not set!\n'sungrow-go -help' lists available parameters.")
	}

	inv.types = strings.Split(inv.data, ",")
	if len(inv.types) < 1 {
		log.Fatalln("Required parameter 'data' not set or invalid value!\n'sungrow-go -help' lists available parameters and values.")
	}
	for _, t := range inv.types {
		switch t {
		case "pv":
		case "battery":
			break
		default:
			log.Fatalf("Invalid value \"%s\" for parameter 'data'!\n'sungrow-go -help' lists available parameters and values.\n", t)
		}
	}
}
