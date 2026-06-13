package main

import (
	"os"
	"os/exec"
	"reflect"
	"testing"
)

// fatalCases maps a case name to an inverter state that must make
// flagsValidate call log.Fatalln (i.e. exit the process).
var fatalCases = map[string]inverter{
	"missingIp":       {ipS: "", user: "admin", password: "x", data: "pv"},
	"invalidIp":       {ipS: "not-an-ip", user: "admin", password: "x", data: "pv"},
	"missingUser":     {ipS: "192.168.2.100", user: "", password: "x", data: "pv"},
	"missingPassword": {ipS: "192.168.2.100", user: "admin", password: "", data: "pv"},
	"invalidData":     {ipS: "192.168.2.100", user: "admin", password: "x", data: "pv,foo"},
	"emptyData":       {ipS: "192.168.2.100", user: "admin", password: "x", data: ""},
}

func TestFlagsValidateSuccess(t *testing.T) {
	inv = &inverter{ipS: "192.168.2.100", user: "admin", password: "secret", data: "pv,battery"}

	flagsValidate()

	if inv.ip == nil {
		t.Fatal("flagsValidate() did not parse a valid ip")
	}
	want := []string{"pv", "battery"}
	if !reflect.DeepEqual(inv.types, want) {
		t.Errorf("inv.types = %v, want %v", inv.types, want)
	}
}

// TestFlagsValidateFatal verifies that invalid flags terminate the process.
// flagsValidate uses log.Fatalln (os.Exit), so each case runs in a subprocess
// and the parent asserts the non-zero exit.
func TestFlagsValidateFatal(t *testing.T) {
	if name := os.Getenv("FATAL_CASE"); name != "" {
		c := fatalCases[name]
		inv = &c
		flagsValidate()
		return // reached only if flagsValidate failed to exit
	}

	for name := range fatalCases {
		t.Run(name, func(t *testing.T) {
			cmd := exec.Command(os.Args[0], "-test.run=TestFlagsValidateFatal")
			cmd.Env = append(os.Environ(), "FATAL_CASE="+name)
			err := cmd.Run()

			ee, ok := err.(*exec.ExitError)
			if !ok || ee.Success() {
				t.Fatalf("case %q: expected non-zero exit, got err = %v", name, err)
			}
		})
	}
}
