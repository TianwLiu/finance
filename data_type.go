package main

import "github.com/plaid/plaid-go/plaid"

//date type for local user info storage
type DataRoot struct {
	WebUsers  []WebUser  `json:"web_users"`
	WebGroups []WebGroup `json:"web_groups"`
}






//date type for web api
type UserAccounts struct {
	UserId   string          `json:"user_id"`
	CashBase CashBase		`json:"cash_base"`
	Accounts []plaid.Account `json:"accounts"`
}
type GroupAccounts struct {
	GroupId string `json:"group_id"`
	CashBase CashBase `json:"cash_base"` 
	UserAccountsList	[]UserAccounts `json:"user_accounts_list"`
}

type UserTransactions struct {
	CashFlow Sheets `json:"cash_flow"`
	Transactions []plaid.Transaction `json:"transactions"`
}

type GroupTransactions struct {
	CashFlow Sheets `json:"cash_flow"`
	Transactions []plaid.Transaction `json:"transactions"`
}
