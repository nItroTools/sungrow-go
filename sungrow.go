package main

import (
	"flag"
	"github.com/nItroTools/sungrow-go/ws"
	"log"
	"net"
	"strings"
)

type inverter struct {
	ws        *ws.WS
	ipS       string
	ip        net.IP
	port      int
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
	inv.ws = ws.NewWS(inv.ip, inv.port, inv.path)
	if err := inv.ws.Connect(); err != nil {
		log.Fatalln(err)
	}
	defer inv.ws.Close()

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

// flags defines, parses and validates command-line flags from os.Args[1:]
func flags() {
	ipS := flag.String("ip", "", "IP address of the Sungrow inverter")
	port := flag.Int("port", 8082, "WebSocket port of the Sungrow inverter")
	path := flag.String("path", "/ws/home/overview", "Server path from where data is requested")
	data := flag.String("data", "pv,battery", "Select the data to be requested comma separated.\nPossible values are \"pv\" and \"battery\"")
	separator := flag.String("separator", ",", "Output data separator")
	flag.Parse()

	inv.ipS = *ipS
	inv.port = *port
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
