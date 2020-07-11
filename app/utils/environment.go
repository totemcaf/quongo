package utils

import (
	"os"
	"strconv"
	"strings"
)

// GetIntEnv ...
func GetIntEnv(name string, defValue int) (int, error) {
	v := strings.TrimSpace(os.Getenv(name))

	if v == "" {
		return defValue, nil
	}

	return strconv.Atoi(v)
}

// GetStrEnv ...
func GetStrEnv(name string, defValue string) string {
	v := strings.TrimSpace(os.Getenv(name))

	if v == "" {
		return defValue
	}

	return v
}
