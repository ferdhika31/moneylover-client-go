package moneylover

import (
	"net/http"
	"testing"
	"time"
)

func TestLoginAndIncomeExpense(t *testing.T) {
	orig := http.DefaultClient.Transport
	defer func() { http.DefaultClient.Transport = orig }()

	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	call := 0
	http.DefaultClient.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		call++
		switch call {
		case 1:
			return newResponse(`{"data":{"request_token":"req","login_url":"https://ml?client=cli"}}`), nil
		case 2:
			return newResponse(`{"access_token":"tok"}`), nil
		case 3:
			if r.URL.String() != "https://web.moneylover.me/api/transaction/add" {
				t.Fatalf("unexpected url %s", r.URL)
			}
			return newResponse(`{"error":0,"data":{"_id":"tx1"}}`), nil
		case 4:
			return newResponse(`{"error":0,"data":{"_id":"tx2"}}`), nil
		default:
			t.Fatalf("unexpected call %d", call)
			return nil, nil
		}
	})

	client, err := Login("e@mail", "pass")
	if err != nil {
		t.Fatalf("login error: %v", err)
	}
	if client.Token != "tok" {
		t.Fatalf("unexpected token %s", client.Token)
	}

	tok, err := LoadTokenForUser("e@mail")
	if err != nil || tok != "tok" {
		t.Fatalf("token not saved")
	}

	p := TransactionParams{WalletID: "w", CategoryID: "c", Amount: "1", Date: time.Now()}
	res, err := client.Income(p)
	if err != nil || res.ID != "tx1" {
		t.Fatalf("income failed")
	}
	res2, err := client.Expense(p)
	if err != nil || res2.ID != "tx2" {
		t.Fatalf("expense failed")
	}
}
