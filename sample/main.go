package main

import (
	"fmt"
	"log"
	"time"

	moneylover "github.com/ferdhika31/moneylover-client-go"
)

func main() {
	// try to load a saved session for a specific email
	token, err := moneylover.LoadTokenForUser("email@example.com")
	var client *moneylover.Client
	if err == nil && token != "" {
		if expired, _ := moneylover.TokenExpired(token); !expired {
			client = moneylover.NewClient(token)
			if _, err = client.GetUserInfo(); err == nil {
				fmt.Println("loaded existing session")
			} else {
				fmt.Println("stored session invalid, logging in again")
				client = nil
			}
		} else {
			fmt.Println("stored session expired, logging in again")
		}
	}
	if client == nil {
		client, err = moneylover.Login("email@example.com", "password")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("logged in and session saved")
	}

	// list wallets
	wallets, err := client.GetWallets()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("wallets:", wallets)

	// create an expense example
	txParams := moneylover.TransactionParams{
		WalletID:   "<walletID>",
		CategoryID: "<categoryID>",
		Amount:     "1000",
		Note:       "Cilok",
		Date:       time.Now(),
	}
	// Expense returns *moneylover.AddTransactionResponse
	tx, err := client.Expense(txParams)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("added transaction:", tx)
}
