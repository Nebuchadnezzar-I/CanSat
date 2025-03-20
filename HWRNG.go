package main

import (
	"fmt"
	"os"
)

func readHWRNG(n int) ([]byte, error) {
	file, err := os.Open("/dev/hwrng")
	if err != nil {
		return nil, fmt.Errorf("failed to open hwrng: %w", err)
	}
	defer file.Close()

	randBytes := make([]byte, n)
	_, err = file.Read(randBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to read from hwrng: %w", err)
	}
	return randBytes, nil
}
