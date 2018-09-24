package utils

import (
	"github.com/ant0ine/go-json-rest/rest"
	"strconv"
)

func IntParam(r *rest.Request, key string, defaultValue int) (int, error) {
	valueStr := r.FormValue(key)
	if valueStr == "" {
		return defaultValue, nil
	}

	return strconv.Atoi(valueStr)
}

func StrParam(r *rest.Request, key string, defaultValue string) (string, error) {
	valueStr := r.FormValue(key)
	if valueStr == "" {
		return defaultValue, nil
	}

	return valueStr, nil
}
