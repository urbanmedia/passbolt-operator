package util

import (
	"os"
	"strconv"
)

// StringEnvOrDefault returns the value of the environment variable named by the key.
// If the variable is not present in the environment, the default value is returned.
func StringEnvOrDefault(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

// IntEnvOrDefault returns the value of the environment variable named by the key.
// If the variable is not present in the environment, the default value is returned.
func IntEnvOrDefault(key string, defaultValue int) int {
	value := StringEnvOrDefault(key, strconv.Itoa(defaultValue))
	if intValue, err := strconv.Atoi(value); err == nil {
		return intValue
	}
	return defaultValue
}
