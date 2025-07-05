package moneylover

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

// roundTripFunc allows mocking http.RoundTripper
type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func newResponse(body string) *http.Response {
	return &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		Header:     make(http.Header),
	}
}

func TestNewClient(t *testing.T) {
	c := NewClient("token123")
	if c.Token != "token123" {
		t.Errorf("expected token to be set")
	}
}

func TestGetToken(t *testing.T) {
	orig := http.DefaultClient.Transport
	defer func() { http.DefaultClient.Transport = orig }()

	call := 0
	http.DefaultClient.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		call++
		switch call {
		case 1:
			if r.URL.String() != "https://web.moneylover.me/api/user/login-url" {
				t.Fatalf("unexpected url %s", r.URL)
			}
			return newResponse(`{"data":{"request_token":"req","login_url":"https://ml?client=cli"}}`), nil
		case 2:
			if r.URL.String() != "https://oauth.moneylover.me/token" {
				t.Fatalf("unexpected url %s", r.URL)
			}
			if r.Header.Get("Authorization") != "Bearer req" {
				t.Fatalf("missing auth header")
			}
			if r.Header.Get("Client") != "cli" {
				t.Fatalf("missing client header")
			}
			return newResponse(`{"access_token":"tok"}`), nil
		default:
			t.Fatalf("unexpected request %d", call)
			return nil, nil
		}
	})

	tok, err := GetToken("email", "pass")
	if err != nil {
		t.Fatalf("GetToken error: %v", err)
	}
	if tok != "tok" {
		t.Errorf("expected tok, got %s", tok)
	}
}

func TestGetUserInfo(t *testing.T) {
	orig := http.DefaultClient.Transport
	defer func() { http.DefaultClient.Transport = orig }()

	http.DefaultClient.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.String() != "https://web.moneylover.me/api/user/info" {
			t.Fatalf("unexpected url %s", r.URL)
		}
		if r.Header.Get("Authorization") != "AuthJWT mytoken" {
			t.Fatalf("wrong auth header")
		}
		return newResponse(`{"error":0,"data":{"_id":"uid"}}`), nil
	})

	c := NewClient("mytoken")
	info, err := c.GetUserInfo()
	if err != nil {
		t.Fatalf("GetUserInfo error: %v", err)
	}
	if info.ID != "uid" {
		t.Errorf("expected uid, got %s", info.ID)
	}
}

func TestAPIRequestError(t *testing.T) {
	orig := http.DefaultClient.Transport
	defer func() { http.DefaultClient.Transport = orig }()

	http.DefaultClient.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return newResponse(`{"error":1,"msg":"bad"}`), nil
	})

	c := NewClient("t")
	err := c.apiRequest("/something", nil, nil, nil)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestAddTransaction(t *testing.T) {
	orig := http.DefaultClient.Transport
	defer func() { http.DefaultClient.Transport = orig }()

	http.DefaultClient.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.String() != "https://web.moneylover.me/api/transaction/add" {
			t.Fatalf("unexpected url %s", r.URL)
		}
		return newResponse(`{"error":0,"data":{"_id":"tx1","account":"w1","category":"c1","amount":100,"note":"","displayDate":"2020-01-01"}}`), nil
	})

	c := NewClient("tok")
	p := TransactionParams{WalletID: "w1", CategoryID: "c1", Amount: "100", Date: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)}
	res, err := c.AddTransaction(p)
	if err != nil {
		t.Fatalf("AddTransaction error: %v", err)
	}
	if res.ID != "tx1" {
		t.Errorf("expected tx1, got %s", res.ID)
	}
}
