package moneylover

import "time"

// DTO structs for Money Lover API responses.

// UserInfo describes the authenticated user.
type UserInfo struct {
	ID               string                 `json:"_id"`
	Email            string                 `json:"email"`
	IconPackage      []string               `json:"icon_package"`
	LimitDevice      int                    `json:"limitDevice"`
	Tags             []string               `json:"tags"`
	ClientSetting    map[string]interface{} `json:"client_setting"`
	Purchased        bool                   `json:"purchased"`
	SubscribeProduct string                 `json:"subscribeProduct"`
	DeviceID         string                 `json:"deviceId"`
}

// WalletUser describes a user that has access to a wallet.
type WalletUser struct {
	ID    string `json:"_id"`
	Email string `json:"email"`
}

// Wallet represents a Money Lover wallet as returned by the API.
type Wallet struct {
	ID                      string              `json:"_id"`
	Name                    string              `json:"name"`
	CurrencyID              int                 `json:"currency_id"`
	Owner                   string              `json:"owner"`
	SortIndex               int                 `json:"sortIndex"`
	TransactionNotification bool                `json:"transaction_notification"`
	Archived                bool                `json:"archived"`
	AccountType             int                 `json:"account_type"`
	ExcludeTotal            bool                `json:"exclude_total"`
	Icon                    string              `json:"icon"`
	ListUser                []WalletUser        `json:"listUser"`
	CreatedAt               string              `json:"createdAt"`
	UpdateAt                string              `json:"updateAt"`
	IsDelete                bool                `json:"isDelete"`
	Balance                 []map[string]string `json:"balance"`
}

// Category represents a transaction category.
type CategoryParent struct {
	ID       string `json:"_id"`
	Name     string `json:"name"`
	Icon     string `json:"icon"`
	Type     int    `json:"type"`
	Metadata string `json:"metadata"`
}

type Category struct {
	ID       string          `json:"_id"`
	Name     string          `json:"name"`
	Icon     string          `json:"icon"`
	Account  string          `json:"account"`
	Type     int             `json:"type"`
	Metadata string          `json:"metadata"`
	Group    int             `json:"group"`
	Parent   *CategoryParent `json:"parent,omitempty"`
}

// Transaction describes a wallet transaction.
type AccountInfo struct {
	ID          string `json:"_id"`
	Name        string `json:"name"`
	CurrencyID  int    `json:"currency_id"`
	AccountType int    `json:"account_type"`
	Icon        string `json:"icon"`
}

type Transaction struct {
	ID            string      `json:"_id"`
	Note          string      `json:"note"`
	Account       AccountInfo `json:"account"`
	Category      Category    `json:"category"`
	Amount        float64     `json:"amount"`
	DisplayDate   string      `json:"displayDate"`
	Remind        int         `json:"remind"`
	Address       string      `json:"address"`
	Longtitude    float64     `json:"longtitude"`
	Latitude      float64     `json:"latitude"`
	With          []string    `json:"with"`
	Campaign      []string    `json:"campaign"`
	LastEditBy    WalletUser  `json:"lastEditBy"`
	ExcludeReport bool        `json:"exclude_report"`
	Images        []string    `json:"images"`
	CreatedAt     string      `json:"createdAt"`
}

type DateRange struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

// TransactionsResponse is returned by GetTransactions.
type TransactionsResponse struct {
	Daterange    DateRange     `json:"daterange"`
	Transactions []Transaction `json:"transactions"`
}

// AddTransactionResponse represents the data returned when creating a transaction.
type AddTransactionResponse struct {
	ID          string   `json:"_id"`
	With        []string `json:"with"`
	Account     string   `json:"account"`
	Category    string   `json:"category"`
	Amount      float64  `json:"amount"`
	Note        string   `json:"note"`
	DisplayDate string   `json:"displayDate"`
	TokenDevice string   `json:"tokenDevice"`
}

// TransactionParams represents parameters used when creating income or expense.
type TransactionParams struct {
	WalletID   string    // wallet/account ID
	CategoryID string    // category ID
	Amount     string    // amount as a string to match API expectations
	Note       string    // optional note
	Date       time.Time // transaction date
}
