package main

import (
	"fmt"
	"github.com/d2r2/go-i2c"
	"time"
)

type SHT40 struct {
	i2c *i2c.I2C
}

func NewSHT40(bus int) (*SHT40, error) {
	i2cMux, err := i2c.NewI2C(0x70, bus)
	if err != nil {
		return nil, err
	}
	defer i2cMux.Close()
	_, err = i2cMux.WriteBytes([]byte{0x01})
	if err != nil {
		return nil, fmt.Errorf("failed to select I2C channel for SHT40: %v", err)
	}

	i2c, err := i2c.NewI2C(0x44, bus)
	if err != nil {
		return nil, err
	}
	return &SHT40{i2c: i2c}, nil
}

func (s *SHT40) ReadSensorData() (float64, float64, error) {
	_, err := s.i2c.WriteBytes([]byte{0xFD})
	if err != nil {
		return 0, 0, fmt.Errorf("failed to send measurement command: %v", err)
	}

	time.Sleep(10 * time.Millisecond)

	data := make([]byte, 6)
	_, err = s.i2c.ReadBytes(data)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to read data: %v", err)
	}

	if !checkCRC(data[:3]) || !checkCRC(data[3:]) {
		return 0, 0, fmt.Errorf("CRC check failed")
	}

	rawTemperature := uint16(data[0])<<8 | uint16(data[1])
	rawHumidity := uint16(data[3])<<8 | uint16(data[4])

	temperature := -45.0 + 175.0*(float64(rawTemperature)/65535.0)
	humidity := 100.0 * (float64(rawHumidity) / 65535.0)

	return temperature, humidity, nil
}

func (s *SHT40) Close() {
	s.i2c.Close()
}

func checkCRC(data []byte) bool {
	polynomial := byte(0x31)
	crc := byte(0xFF)
	for _, b := range data[:2] {
		crc ^= b
		for i := 0; i < 8; i++ {
			if crc&0x80 != 0 {
				crc = (crc << 1) ^ polynomial
			} else {
				crc <<= 1
			}
		}
	}
	return crc == data[2]
}
