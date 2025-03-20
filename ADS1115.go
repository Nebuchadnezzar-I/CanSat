package main

import (
	"fmt"
	"github.com/d2r2/go-i2c"
	"time"
)

type ADS1115 struct {
	i2c *i2c.I2C
}

func NewADS1115(bus int) (*ADS1115, error) {
	i2cMux, err := i2c.NewI2C(0x70, bus)
	if err != nil {
		return nil, err
	}
	defer i2cMux.Close()
	_, err = i2cMux.WriteBytes([]byte{0x04})
	if err != nil {
		return nil, fmt.Errorf("failed to select I2C channel for ADS1115: %v", err)
	}

	i2c, err := i2c.NewI2C(0x48, bus)
	if err != nil {
		return nil, err
	}
	return &ADS1115{i2c: i2c}, nil
}

func (adc *ADS1115) Close() {
	adc.i2c.Close()
}

func (adc *ADS1115) ReadVoltage() (float64, error) {
	config := []byte{0x84, 0x83}
	_, err := adc.i2c.WriteBytes([]byte{0x01, config[0], config[1]})
	if err != nil {
		return 0, fmt.Errorf("failed to configure ADS1115: %v", err)
	}

	time.Sleep(10 * time.Millisecond)

	data := make([]byte, 2)
	_, err = adc.i2c.ReadBytes(data)
	if err != nil {
		return 0, fmt.Errorf("failed to read ADC value: %v", err)
	}

	rawValue := int16(uint16(data[0])<<8 | uint16(data[1]))

	voltage := float64(rawValue) * (4.096 / 32768.0)

	if voltage < 0 {
		voltage = 0
	}

	return voltage, nil
}
