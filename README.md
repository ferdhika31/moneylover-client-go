# Money Lover Go Library

A minimal Go client for the [Money Lover](https://moneylover.me/) API. **This library is unofficial and not affiliated with Money Lover.**

## Installation

```bash
go get github.com/ferdhika31/moneylover-client-go
```

## Usage

Authenticate with your Money Lover credentials and use the returned `Client` to access the API. API methods return typed DTO structs. The sample in the `sample` directory demonstrates the basic flow.

```go
package main

import (
    "fmt"
    "log"
    "time"

    ml "github.com/ferdhika31/moneylover-client-go"
)

func main() {
    client, err := ml.Login("tamvan@dika.web.id", "password")
    if err != nil {
        log.Fatal(err)
    }

    txParams := ml.TransactionParams{
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
    fmt.Println(tx)
}
```

Run the example with:

```bash
cd sample
go run .
```

### Session handling

`Login` saves the JWT token to `~/.moneylover-client` keyed by your email
address. Load it on the next run with `LoadTokenForUser("email")` and only log
in again when the saved session is invalid:

```go
token, err := ml.LoadTokenForUser("tamvan@dika.web.id")
var client *ml.Client
if err == nil && token != "" {
    client = ml.NewClient(token)
    if _, err := client.GetUserInfo(); err != nil {
        client, err = ml.Login("tamvan@dika.web.id", "password")
    }
} else {
    client, err = ml.Login("tamvan@dika.web.id", "password")
}
```

Use `SaveTokenForUser(email, token)` to store sessions for additional accounts
and `LoadTokenForUser(email)` to restore them later. Call `TokenExpired(token)`
to check whether a stored JWT is still valid before hitting the API.

### ID placeholders

The sample JSON below uses placeholder IDs so it's easier to read:

- `<walletID>` – your wallet or account identifier
- `<categoryID>` – a category identifier
- `<userID>` – your Money Lover user ID
- `<transactionID>` – a transaction identifier

### Sample wallet response

A successful call to `GetWallets` returns data like:

```json
{
    "error": 0,
    "msg": "get_list_account_success",
    "action": "wallet_list",
    "data": [
        {
            "_id": "<walletID>",
            "name": "BRI",
            "currency_id": 44,
            "owner": "<userID>",
            "sortIndex": 19,
            "transaction_notification": false,
            "archived": true,
            "account_type": 0,
            "exclude_total": true,
            "icon": "icon_56",
            "listUser": [
                {
                    "_id": "<userID>",
                    "email": "tamvan@dika.web.id"
                }
            ],
            "createdAt": "2016-08-13T05:54:20.031Z",
            "updateAt": "2024-01-04T19:18:40.558Z",
            "isDelete": false,
            "balance": [
                {
                    "IDR": "14000.00"
                }
            ]
        }
    ]
}
```

### Sample user info response

When authenticated, `GetUserInfo` responds with data similar to:

```json
{
    "error": 0,
    "msg": "get_info_success",
    "action": "user_info",
    "data": {
        "_id": "<userID>",
        "email": "tamvan@dika.web.id",
        "icon_package": [
            "icon_ml_lunar_new_year",
            "worldcup_icon_pack_2014",
            "summer_icon_pack_2014"
        ],
        "limitDevice": 5,
        "tags": [
            "country:id",
            "device:android",
            "utm_source:google_play",
            "utm_medium:organic",
            "check_linked",
            "device:web",
            "buy_pre_fail_1:3",
            "goal_wallet",
            "check_store",
            "check_store_lw",
            "Premium_store_ver1",
            "purchase:premium",
            "device:ios",
            "credit_wallet",
            "discount:60",
            "campaign_discount:c74e5e9194ed45cf84cca36e14e92593",
            "gc_migrate_process:finish",
            "ques1_1",
            "money_insider:year:27052024:premium_sub_insider_year_1"
        ],
        "client_setting": {
            "is_done_tooltip_txn_ai_add_trans": false,
            "is_done_tooltip_txn_ai_home": false,
            "is_done_intro_money_insider_detail": true,
            "is_done_intro_money_insider_home": true,
            "ob_step_create_wallet": false,
            "lw__banner_status2": true,
            "om": 0,
            "sb": 1,
            "er": true,
            "fmy": 0,
            "fdw": 2,
            "pl": 1,
            "av": 1703,
            "sl": false,
            "dr": 20,
            "fd": 28,
            "df": 0,
            "l": "id",
            "ps": 1,
            "ds": 1,
            "sa": false,
            "sc": true,
            "sd": false,
            "show_advance_add_transaction": true,
            "future_period": 1,
            "nps__last_ask": 1751065732,
            "su": true,
            "ls": 1751065731568,
            "main_currency": "IDR",
            "ol": false,
            "ob_step_closed": false,
            "is_show_step_view_for_user": false,
            "ob_step_add_transaction": true,
            "ob_budget_suggest_show": true,
            "ob_step_add_budget": true,
            "ob_step_get_it": false,
            "is_done_intro_home": true,
            "update_password_fb": "false"
        },
        "purchased": true,
        "subscribeProduct": "premium_sub_year_1",
        "deviceId": "415c05be-b992-477a-afac-1db54c8d3c41"
    }
}
```

If the token is invalid, the response looks like:

```json
{
    "error": 1,
    "msg": "user_unauthenticated",
    "action": "user_info"
}
```

### Sample categories response

A successful call to `GetCategories` returns data like:

```json
{
    "error": 0,
    "msg": "get_cate_success",
    "action": "category_list",
    "data": [
        {
            "_id": "<categoryID>",
            "name": "Tagihan & utilitas",
            "icon": "icon_135",
            "account": "<walletID>",
            "type": 2,
            "metadata": "utilities0",
            "group": 0
        },
        {
            "_id": "<categoryID>",
            "name": "Transportasi",
            "icon": "icon/icon_ml_lunar_new_year/icon_ml_lunar_new_year_2",
            "account": "<walletID>",
            "type": 2,
            "metadata": "transport0",
            "group": 0
        },
        {
            "_id": "<categoryID>",
            "name": "Belanja",
            "icon": "ic_category_shopping",
            "account": "<walletID>",
            "type": 2,
            "metadata": "shopping0",
            "group": 0
        }
    ],
    "walletId": "<walletID>"
}
```

The error format looks like:

```json
{
    "error": 1,
    "msg": "web_get_category_error",
    "action": "category_list"
}
```

### Sample transactions response

A call to `GetTransactions` returns transaction data and the requested date range:

```json
{
    "error": 0,
    "msg": "get_transaction_success",
    "action": "transaction_list",
    "data": {
        "daterange": {
            "startDate": "2016-08-13",
            "endDate": "2025-07-05"
        },
        "transactions": [
            {
                "_id": "<transactionID>",
                "note": "",
                "account": {
                    "_id": "<walletID>",
                    "name": "Keuangan Rumah Tangga",
                    "currency_id": 44,
                    "account_type": 0,
                    "icon": "icon"
                },
                "category": {
                    "_id": "<categoryID>",
                    "name": "Keluarga",
                    "icon": "icon/icon_ml_lunar_new_year/icon_ml_lunar_new_year_4",
                    "account": "<walletID>",
                    "type": 2,
                    "metadata": "",
                    "parent": {
                        "_id": "<categoryID>",
                        "name": "Hadiah & Donasi",
                        "icon": "ic_category_donations",
                        "type": 2,
                        "metadata": "gifts_donations0"
                    }
                },
                "amount": 600000,
                "displayDate": "2025-05-29T00:00:00.000Z",
                "remind": 0,
                "address": "",
                "longtitude": 0,
                "latitude": 0,
                "with": [
                    "Ayah"
                ],
                "campaign": [],
                "lastEditBy": {
                    "_id": "<userID>",
                    "email": "tamvan@dika.web.id"
                },
                "exclude_report": false,
                "images": [],
                "createdAt": "2025-05-29T11:17:45.474Z"
            }
        ]
    }
}
```

If no data is found, the response contains an empty array:

```json
{
    "error": 0,
    "msg": "cashbook_no_data",
    "action": "transaction_list",
    "data": {
        "daterange": {
            "startDate": "2025-07-05",
            "endDate": "2025-07-05"
        },
        "transactions": []
    }
}
```

An invalid request yields:

```json
{
    "error": 1,
    "msg": "sync_error_have_not_permission",
    "action": "transaction_list"
}
```

### Sample add transaction response

A successful call to add a transaction returns the created data:

```json
{
    "error": 0,
    "msg": "transaction_add_success",
    "action": "transaction_create",
    "data": {
        "_id": "<transactionID>",
        "with": [
            "Ama"
        ],
        "account": "<walletID>",
        "category": "<categoryID>",
        "amount": 400000,
        "note": "Ama gendit",
        "displayDate": "2025-07-05",
        "tokenDevice": "web"
    }
}
```

The error format may contain messages like `sync_error_data_invalid_format` or `sync_error_have_not_permission`:

```json
{
    "error": 1,
    "msg": "sync_error_data_invalid_format",
    "action": "transaction_create"
}
```

```json
{
    "error": 1,
    "msg": "sync_error_have_not_permission",
    "action": "transaction_create"
}
```

## Contributing

- Fork the repository.
- Create a new branch for your feature or bugfix.
- Commit your changes and push them to your fork.
- Create a Pull Request describing your changes.

We appreciate all contributions, whether they are documentation improvements, bug fixes, or new features!
