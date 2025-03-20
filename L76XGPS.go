package main

import (
	"os/exec"
	"strings"
)

type L76XGPS struct{}

func NewL76XGPS() *L76XGPS {
	return &L76XGPS{}
}

func (gps *L76XGPS) ReadGPS() (string, error) {
	cmd := exec.Command("timeout", "0.2", "cat", "/dev/serial0")
	output, err := cmd.Output()
	if err != nil {
		return "No GPS", nil
	}

	gpsData := strings.TrimSpace(string(output))
	if gpsData == "" {
		return "No GPS", nil
	}

	return gpsData, nil
}
