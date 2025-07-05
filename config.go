package moneylover

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// configFile returns the path to the configuration file used to store the token.
func configFile() string {
	dir, _ := os.UserHomeDir()
	return filepath.Join(dir, ".moneylover-client")
}

// SaveToken stores the JWT token in the config file.
func SaveToken(token string) error {
	return SaveTokenForUser("jwtToken", token)
}

// SaveTokenForUser stores the JWT token for the given email.
func SaveTokenForUser(email, token string) error {
	cfg, err := readTokenMap()
	if err != nil {
		return err
	}
	cfg[email] = token
	return writeTokenMap(cfg)
}

// LoadToken reads the stored JWT token from the config file.
func LoadToken() (string, error) {
	return LoadTokenForUser("jwtToken")
}

// LoadTokenForUser reads the stored JWT token for the given email.
func LoadTokenForUser(email string) (string, error) {
	cfg, err := readTokenMap()
	if err != nil {
		return "", err
	}
	tok, ok := cfg[email]
	if !ok {
		return "", errors.New("token not found")
	}
	return tok, nil
}

// ClearToken removes the config file.
func ClearToken() error {
	return ClearTokenForUser("jwtToken")
}

func readTokenMap() (map[string]string, error) {
	f, err := os.Open(configFile())
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]string{}, nil
		}
		return nil, err
	}
	defer f.Close()
	var cfg map[string]string
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func writeTokenMap(cfg map[string]string) error {
	f, err := os.Create(configFile())
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(cfg)
}

// ClearTokenForUser removes the stored token for the given email.
func ClearTokenForUser(email string) error {
	cfg, err := readTokenMap()
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	delete(cfg, email)
	if len(cfg) == 0 {
		return os.Remove(configFile())
	}
	return writeTokenMap(cfg)
}
