# sungrow-go
GoLang implementation for accessing real-time data from Sungrow inverters with WiNet-S dongle using WebSocket.

## Install
```bash
$ go install .
```

## Usage
List available and required parameters
```bash
$ sungrow-go -help
```
Basic usage with ip address of your inverter (e.g. `192.168.2.100`)
```bash
$ sungrow-go -ip 192.168.2.100
```
Output: `var,value,unit`
```text
inverterTemp,43.600,℃
netFeedIn,1.500,kW
netPower,0.000,kW
totalConsumption,0.368,kW
sunPower,4.068,kW
batteryCharge,2.200,kW
batteryDischarge,0.000,kW
batteryTemp,23.000,℃
batteryLevel,43.100,%
batteryHealth,99.000,%
```

## Supported inverters
Tested Sungrow inverters with WiNet-S dongle:
- SH10RT