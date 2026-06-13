package main

import "github.com/nItroTools/sungrow-go/ws"

var (
	pvKeys = ws.Keys{
		"I18N_COMMON_TOTAL_DCPOWER":                   "sunPower",
		"I18N_COMMON_PV_DAYILY_ENERGY_GENERATION":     "sunPowerDay",
		"I18N_COMMON_PV_TOTAL_ENERGY_GENERATION":      "sunPowerTotal",
		"I18N_COMMON_FEED_NETWORK_TOTAL_ACTIVE_POWER": "netFeedIn",
		"I18N_COMMON_DAILY_FEED_NETWORK_VOLUME":       "netFeedInDay",
		"I18N_COMMON_TOTAL_FEED_NETWORK_VOLUME":       "netFeedInTotal",
		"I18N_CONFIG_KEY_4060":                        "netPower",
		"I18N_COMMON_ENERGY_GET_FROM_GRID_DAILY":      "netPowerDay",
		"I18N_COMMON_TOTAL_ELECTRIC_GRID_GET_POWER":   "netPowerTotal",
		"I18N_COMMON_LOAD_TOTAL_ACTIVE_POWER":         "consumption",
		"I18N_CONFIG_KEY_1001188":                     "consumptionRate",
		"I18N_COMMON_AIR_TEM_INSIDE_MACHINE":          "inverterTemp",
		"I18N_COMMON_DEVICE_STATUS":                   "deviceStatus",
	}

	batteryKeys = ws.Keys{
		"I18N_COMMON_BATTERY_SOC":         "batteryLevel",
		"I18N_CONFIG_KEY_3907":            "batteryCharge",
		"I18N_CONFIG_KEY_3921":            "batteryDischarge",
		"I18N_COMMON_BATTARY_HEALTH":      "batteryHealth",
		"I18N_COMMON_BATTERY_TEMPERATURE": "batteryTemp",
	}

	// deviceStatusValues maps the i18n status values of I18N_COMMON_DEVICE_STATUS
	// to human readable German labels. Unknown values are passed through verbatim.
	deviceStatusValues = ws.Keys{
		"I18N_COMMON_ON_GRID_OPERATION":  "Netzbetrieb",
		"I18N_COMMON_OFF_GRID_OPERATION": "Notstrombetrieb",
		"I18N_COMMON_STANDBY":            "Standby",
		"I18N_COMMON_STOP":               "Gestoppt",
		"I18N_COMMON_FAULT":              "Fehler",
		"I18N_COMMON_INITIALIZING":       "Initialisierung",
		"I18N_COMMON_DERATING":           "Leistungsbegrenzung",
	}

	// pvValueKeys marks pv data fields whose data_value is itself an i18n string
	// and maps each to its value-translation table.
	pvValueKeys = map[string]ws.Keys{
		"I18N_COMMON_DEVICE_STATUS": deviceStatusValues,
	}
)
