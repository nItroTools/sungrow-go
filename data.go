package main

import "github.com/nItroTools/sungrow-go/ws"

var (
	pvKeys = ws.Keys{
		"I18N_COMMON_TOTAL_DCPOWER":                   "sunPower",
		"I18N_COMMON_FEED_NETWORK_TOTAL_ACTIVE_POWER": "netFeedIn",
		"I18N_CONFIG_KEY_4060":                        "netPower",
		"I18N_COMMON_LOAD_TOTAL_ACTIVE_POWER":         "totalConsumption",
		"I18N_COMMON_AIR_TEM_INSIDE_MACHINE":          "inverterTemp",
	}

	batteryKeys = ws.Keys{
		"I18N_COMMON_BATTERY_SOC":         "batteryLevel",
		"I18N_CONFIG_KEY_3907":            "batteryCharge",
		"I18N_CONFIG_KEY_3921":            "batteryDischarge",
		"I18N_COMMON_BATTARY_HEALTH":      "batteryHealth",
		"I18N_COMMON_BATTERY_TEMPERATURE": "batteryTemp",
	}
)
