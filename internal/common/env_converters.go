package common

import (
	"strconv"
	"time"
)

func IntEnv(getEnv func(string) string, key string, defaultValue int) int {
	value := getEnv(key)

	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return intValue
}

func DurationEnv(getEnv func(string) string, key string, defaultValue time.Duration) time.Duration {
	value := getEnv(key)
	if value == "" {
		return defaultValue
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}

	return duration
}

func StringEnv(getEnv func(string) string, key string, defaultValue string) string {
	value := getEnv(key)
	if value == "" {
		return defaultValue
	}

	return value
}

func BoolEnv(getEnv func(string) string, key string, defaultValue bool) bool {
	value := getEnv(key)
	if value == "" || (value != "true" && value != "false") {
		return defaultValue
	}

	return value == "true"
}
