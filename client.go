package moneylover

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Client represents a Money Lover client using JWT token authentication.
type Client struct {
	Token string
}

// NewClient creates a new Client with the given JWT token.
func NewClient(token string) *Client {
	return &Client{Token: token}
}

// apiRequest performs a POST request to the Money Lover API and decodes the JSON response into v.
func (c *Client) apiRequest(path string, body io.Reader, headers map[string]string, v interface{}) error {
	req, err := http.NewRequest("POST", "https://web.moneylover.me/api"+path, body)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "AuthJWT "+c.Token)
	req.Header.Set("Cache-Control", "no-cache, max-age=0, no-store, no-transform, must-revalidate")
	for k, val := range headers {
		req.Header.Set(k, val)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	var result struct {
		Error   int             `json:"error"`
		E       int             `json:"e"`
		Msg     string          `json:"msg"`
		Message string          `json:"message"`
		Data    json.RawMessage `json:"data"`
	}
	if err := dec.Decode(&result); err != nil {
		return err
	}

	if result.Error != 0 {
		return fmt.Errorf("error %d: %s", result.Error, result.Msg)
	}
	if result.E != 0 {
		return fmt.Errorf("error %d: %s", result.E, result.Message)
	}
	if v != nil {
		return json.Unmarshal(result.Data, v)
	}
	return nil
}

// GetToken authenticates with email and password and returns a JWT access token.
func GetToken(email, password string) (string, error) {
	loginRes, err := http.Post("https://web.moneylover.me/api/user/login-url", "application/json", nil)
	if err != nil {
		return "", err
	}
	defer loginRes.Body.Close()
	var loginData struct {
		Data struct {
			RequestToken string `json:"request_token"`
			LoginURL     string `json:"login_url"`
		} `json:"data"`
	}
	if err := json.NewDecoder(loginRes.Body).Decode(&loginData); err != nil {
		return "", err
	}

	clientParam := ""
	if u, err := url.Parse(loginData.Data.LoginURL); err == nil {
		clientParam = u.Query().Get("client")
	}
	if clientParam == "" {
		return "", errors.New("unable to parse client from login url")
	}

	form := url.Values{}
	form.Set("email", email)
	form.Set("password", password)
	req, err := http.NewRequest("POST", "https://oauth.moneylover.me/token", strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+loginData.Data.RequestToken)
	req.Header.Set("Client", clientParam)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	var tokenRes struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(res.Body).Decode(&tokenRes); err != nil {
		return "", err
	}
	if tokenRes.AccessToken == "" {
		return "", errors.New("no access token returned")
	}
	return tokenRes.AccessToken, nil
}

// GetUserInfo retrieves the user information for the current client.
func (c *Client) GetUserInfo() (*UserInfo, error) {
	var info UserInfo
	err := c.apiRequest("/user/info", nil, nil, &info)
	return &info, err
}

// GetWallets returns wallet information for the user.
func (c *Client) GetWallets() ([]Wallet, error) {
	var wallets []Wallet
	err := c.apiRequest("/wallet/list", nil, nil, &wallets)
	return wallets, err
}

// GetCategories retrieves categories for a specific wallet.
func (c *Client) GetCategories(walletID string) ([]Category, error) {
	form := url.Values{}
	form.Set("walletId", walletID)
	headers := map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
	var categories []Category
	err := c.apiRequest("/category/list", strings.NewReader(form.Encode()), headers, &categories)
	return categories, err
}

// GetTransactions retrieves transactions for a wallet between two dates.
func (c *Client) GetTransactions(walletID string, startDate, endDate string) (*TransactionsResponse, error) {
	body := map[string]string{
		"startDate": startDate,
		"endDate":   endDate,
		"walletId":  walletID,
	}
	b, _ := json.Marshal(body)
	headers := map[string]string{"Content-Type": "application/json"}
	var data TransactionsResponse
	err := c.apiRequest("/transaction/list", strings.NewReader(string(b)), headers, &data)
	return &data, err
}

// AddTransaction adds a transaction.
func (c *Client) AddTransaction(p TransactionParams) (*AddTransactionResponse, error) {
	body := map[string]interface{}{
		"with":        []interface{}{},
		"account":     p.WalletID,
		"category":    p.CategoryID,
		"amount":      p.Amount,
		"note":        p.Note,
		"displayDate": p.Date.Format("2006-01-02"),
	}
	b, _ := json.Marshal(body)
	headers := map[string]string{"Content-Type": "application/json"}
	var data AddTransactionResponse
	err := c.apiRequest("/transaction/add", strings.NewReader(string(b)), headers, &data)
	return &data, err
}

const (
	CategoryTypeIncome  = 1
	CategoryTypeExpense = 2
)
