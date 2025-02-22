package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PostgreHost     string
	PostgrePort     uint
	PostgreUser     string
	PostgrePassword string
	PostgreDbName   string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	config := &Config{
		PostgreHost:     os.Getenv("PostgreHost"),
		PostgrePort:     convertStringToUInt(os.Getenv("PostgrePort")),
		PostgreUser:     os.Getenv("PostgreUser"),
		PostgrePassword: os.Getenv("PostgrePassword"),
		PostgreDbName:   os.Getenv("PostgreDbName"),
	}

	return config, nil
}

func convertStringToUInt(value string) uint {
	num, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		panic(err)
	}

	uintNum := uint(num)
	return uintNum
}
