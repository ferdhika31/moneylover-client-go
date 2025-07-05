package moneylover

import (
	"encoding/base64"
	"encoding/json"
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

func TestSaveToken(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	if err := SaveToken("tok"); err != nil {
		t.Fatalf("SaveToken error: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(dir, ".moneylover-client"))
	if err != nil {
		t.Fatalf("read error: %v", err)
	}
	var m map[string]string
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("json error: %v", err)
	}
	if m["jwtToken"] != "tok" {
		t.Fatalf("unexpected map %v", m)
	}
}

func TestSaveTokenForUser(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	if err := SaveTokenForUser("a", "tok1"); err != nil {
		t.Fatalf("save 1: %v", err)
	}
	if err := SaveTokenForUser("b", "tok2"); err != nil {
		t.Fatalf("save 2: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(dir, ".moneylover-client"))
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	var m map[string]string
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if m["a"] != "tok1" || m["b"] != "tok2" {
		t.Fatalf("unexpected map %v", m)
	}
}

func TestLoadToken(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	os.WriteFile(filepath.Join(dir, ".moneylover-client"), []byte(`{"jwtToken":"t"}`), 0600)
	tok, err := LoadToken()
	if err != nil {
		t.Fatalf("LoadToken error: %v", err)
	}
	if tok != "t" {
		t.Fatalf("unexpected token %s", tok)
	}
}

func TestClearToken(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	os.WriteFile(filepath.Join(dir, ".moneylover-client"), []byte(`{"jwtToken":"t"}`), 0600)
	if err := ClearToken(); err != nil {
		t.Fatalf("ClearToken error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, ".moneylover-client")); !os.IsNotExist(err) {
		t.Fatalf("file still exists")
	}
}
