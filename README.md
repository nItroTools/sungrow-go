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
time,2023-08-25T16:34:00+02:00,RFC3339
sunPowerDay,3.500,kWh
sunPowerTotal,97499.900,kWh
inverterTemp,43.600,℃
consumptionRate,57.143,%
netFeedIn,1.500,kW
netPower,0.000,kW
netFeedInDay,0.100,kWh
netFeedInTotal,55694.800,kWh
netPowerDay,0.700,kWh
netPowerTotal,15586.600,kWh
consumption,0.500,kW
sunPower,4.000,kW
batteryCharge,2.500,kW
batteryDischarge,0.000,kW
batteryTemp,23.000,℃
batteryLevel,43.100,%
batteryHealth,99.000,%
```

## Supported inverters
Tested Sungrow inverters with WiNet-S dongle:
- SH10RT