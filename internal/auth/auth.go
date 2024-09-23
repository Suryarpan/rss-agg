package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetApiKey(headers http.Header) (string, error) {
	val := headers.Get("X-Auth")
	if val == "" {
		return "", errors.New("no authentication present")
	}

	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return "", errors.New("auth header is malformed")
	}

	if vals[0] != "ApiKey" {
		return "", errors.New("auth header is malformed")
	}
	return vals[1], nil
}
