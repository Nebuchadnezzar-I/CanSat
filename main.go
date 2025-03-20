package main

import (
	"fmt"
	"log"
	"time"

	"github.com/d2r2/go-logger"
	"golang.org/x/crypto/chacha20poly1305"
)

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func main() {
	logger.ChangePackageLogLevel("i2c", logger.FatalLevel)

	for {
		startTime := time.Now()

		// SHT40
		sensor, err := NewSHT40(1)
		if err != nil {
			fmt.Println("Failed to initialize SHT40:", err)
			continue
		}
		defer sensor.Close()

		temp, hum, err := sensor.ReadSensorData()
		if err != nil {
			fmt.Println("Error reading sensor data:", err)
			continue
		}

		// ADS1115
		adc, err := NewADS1115(1)
		if err != nil {
			fmt.Println("Failed to initialize ADS1115:", err)
			continue
		}
		defer adc.Close()

		voltage, err := adc.ReadVoltage()
		if err != nil {
			fmt.Println("Error reading voltage:", err)
			continue
		}

		// BMP388 Sensor
		bmp388, err := NewBMP388(1)
		if err != nil {
			fmt.Println("Failed to initialize BMP388:", err)
			continue
		}
		defer bmp388.Close()

		pressure, err := bmp388.ReadPressure()
		if err != nil {
			fmt.Println("Error reading pressure:", err)
			continue
		}

		// L76XGPS
		gps := NewL76XGPS()
		gpsData, err := gps.ReadGPS()
		if err != nil {
			fmt.Println("Error reading GPS data:", err)
			continue
		}

		message := fmt.Sprintf(
			"Temperature: %.2f°C\nHumidity: %.2f%%RH\nVoltage: %.3fV\nPressure: %.2f Pa\nGPS Data: %s\n",
			temp, hum, voltage, pressure, gpsData,
		)

		// Encryption
		key := make([]byte, chacha20poly1305.KeySize)
		nonce := make([]byte, chacha20poly1305.NonceSizeX)

		key, err = readHWRNG(chacha20poly1305.KeySize) // Key from hwrng
		if err != nil {
			log.Fatal("Failed to read key from hwrng:", err)
		}

		nonce, err = readHWRNG(chacha20poly1305.NonceSizeX) // Nonce from hwrng
		if err != nil {
			log.Fatal("Failed to read nonce from hwrng:", err)
		}

		messageBytes := []byte(message)

		// Encrypt
		ciphertext, err := Encrypt(key, nonce, messageBytes)
		if err != nil {
			log.Fatal("Encryption failed:", err)
		}

		// Decrypt
		plaintext, err := Decrypt(key, nonce, ciphertext)
		if err != nil {
			log.Fatal("Decryption failed:", err)
		}

		duration := time.Since(startTime)

		clearScreen()
		fmt.Println("Encrypted Message:\n", string(ciphertext))
		fmt.Println("Decrypted Message:\n" + string(plaintext))
		fmt.Printf("\nCycle completed in: %v\n", duration)

		time.Sleep(1 * time.Second)
	}
}
