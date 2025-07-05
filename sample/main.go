package main

import (
	"fmt"
	"log"
	"time"

	moneylover "github.com/ferdhika31/moneylover-client-go"
)

func main() {
	// try to load a saved session
	token, err := moneylover.LoadToken()
	var client *moneylover.Client
	if err == nil && token != "" {
		client = moneylover.NewClient(token)
		if _, err = client.GetUserInfo(); err == nil {
			fmt.Println("loaded existing session")
		} else {
			fmt.Println("stored session expired, logging in again")
			client = nil
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
