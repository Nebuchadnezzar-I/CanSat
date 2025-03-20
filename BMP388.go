package main

import (
	"fmt"
	"github.com/d2r2/go-i2c"
	"time"
)

type BMP388 struct {
	i2c *i2c.I2C
}

func NewBMP388(bus int) (*BMP388, error) {
	i2cMux, err := i2c.NewI2C(0x70, bus)
	if err != nil {
		return nil, err
	}
	defer i2cMux.Close()
	_, err = i2cMux.WriteBytes([]byte{0x02})
	if err != nil {
		return nil, fmt.Errorf("failed to select I2C channel for BMP388: %v", err)
	}

	i2c, err := i2c.NewI2C(0x77, bus)
	if err != nil {
		return nil, err
	}
	return &BMP388{i2c: i2c}, nil
}

func (bmp *BMP388) Close() {
	bmp.i2c.Close()
}

func (bmp *BMP388) ReadPressure() (float64, error) {
	_, err := bmp.i2c.WriteBytes([]byte{0x1B, 0x13})
	if err != nil {
		return 0, fmt.Errorf("failed to configure BMP388: %v", err)
	}

	time.Sleep(10 * time.Millisecond)

	data := make([]byte, 3)
	for i := 0; i < 3; i++ {
		_, err = bmp.i2c.ReadBytes(data[i : i+1])
		if err != nil {
			return 0, fmt.Errorf("failed to read BMP388 pressure data: %v", err)
		}
	}

	rawPressure := (uint32(data[0]) << 16) | (uint32(data[1]) << 8) | uint32(data[2])
	rawPressure >>= 4 // BMP388 pressure data is 20-bit

	pressure := float64(rawPressure) * 1.25

	return pressure, nil
}
