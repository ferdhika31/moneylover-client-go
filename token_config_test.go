package moneylover

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestTokenExpiredInvalid(t *testing.T) {
	if _, err := TokenExpired("bad"); err == nil {
		t.Fatalf("expected error")
	}
}

func TestTokenExpiredNoExp(t *testing.T) {
	// header.payload.signature with payload having no exp
	payload := base64.RawStdEncoding.EncodeToString([]byte(`{"sub":"1"}`))
	token := "a." + payload + ".b"
	if _, err := TokenExpired(token); err == nil {
		t.Fatalf("expected error")
	}
}

func TestTokenExpiredFalse(t *testing.T) {
	exp := time.Now().Add(1 * time.Hour).Unix()
	payload := base64.RawStdEncoding.EncodeToString([]byte(fmt.Sprintf(`{"exp":%d}`, exp)))
	token := "a." + payload + ".b"
	expired, err := TokenExpired(token)
	if err != nil || expired {
		t.Fatalf("unexpected result %v %v", expired, err)
	}
}

func TestSaveTokenForUserReadError(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	p := filepath.Join(dir, ".moneylover-client")
	os.WriteFile(p, []byte("{"), 0600)
	if err := SaveTokenForUser("e", "tok"); err == nil {
		t.Fatalf("expected error")
	}
}

func TestSaveTokenForUserWriteError(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	p := filepath.Join(dir, ".moneylover-client")
	os.Mkdir(p, 0700)
	if err := SaveTokenForUser("e", "tok"); err == nil {
		t.Fatalf("expected error")
	}
}

func TestLoadTokenForUserNotFound(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	os.WriteFile(filepath.Join(dir, ".moneylover-client"), []byte(`{"other":"tok"}`), 0600)
	if _, err := LoadTokenForUser("e"); err == nil {
		t.Fatalf("expected error")
	}
}

func TestClearTokenForUserReadError(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	p := filepath.Join(dir, ".moneylover-client")
	os.Mkdir(p, 0700)
	if err := ClearTokenForUser("e"); err == nil {
		t.Fatalf("expected error")
	}
}
