package moneylover

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

// TokenExpired returns true if the given JWT token is expired based on its `exp` claim.
func TokenExpired(token string) (bool, error) {
	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		return false, errors.New("invalid token")
	}
	payload, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return false, err
	}
	var claims struct {
		Exp int64 `json:"exp"`
	}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return false, err
	}
	if claims.Exp == 0 {
		return false, errors.New("exp not found")
	}
	return time.Now().Unix() > claims.Exp, nil
}
