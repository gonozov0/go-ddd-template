package envutils

import (
	"os"
	"strconv"
)

func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return fallback
}

func GetEnvFloat(key string, fallback float64) float64 {
	if strValue, exists := os.LookupEnv(key); exists {
		if floatValue, err := strconv.ParseFloat(strValue, 64); err == nil {
			return floatValue
		}
	}

	return fallback
}

func GetEnvBool(key string, fallback bool) bool {
	if strValue, exists := os.LookupEnv(key); exists {
		if boolValue, err := strconv.ParseBool(strValue); err == nil {
			return boolValue
		}
	}

	return fallback
}
