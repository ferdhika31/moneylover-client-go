package moneylover

import (
	"encoding/json"
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
	cfg := map[string]string{"jwtToken": token}
	f, err := os.Create(configFile())
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(cfg)
}

// LoadToken reads the stored JWT token from the config file.
func LoadToken() (string, error) {
	f, err := os.Open(configFile())
	if err != nil {
		return "", err
	}
	defer f.Close()
	var cfg map[string]string
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return "", err
	}
	return cfg["jwtToken"], nil
}

// ClearToken removes the config file.
func ClearToken() error {
	return os.Remove(configFile())
}
