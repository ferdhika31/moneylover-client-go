package moneylover

// Login authenticates with email and password and stores the JWT token.
func Login(email, password string) (*Client, error) {
	token, err := GetToken(email, password)
	if err != nil {
		return nil, err
	}
	if err := SaveToken(token); err != nil {
		return nil, err
	}
	return NewClient(token), nil
}

// Income creates an income transaction.
func (c *Client) Income(p TransactionParams) (*AddTransactionResponse, error) {
	return c.AddTransaction(p)
}

// Expense creates an expense transaction.
func (c *Client) Expense(p TransactionParams) (*AddTransactionResponse, error) {
	return c.AddTransaction(p)
}
