package main

import (
	"github.com/plaid/plaid-go/plaid"
	"sort"
	"time"
)

func sortTransactions(transactions []plaid.Transaction) {
	const shorterm = "2006-01-02"
	sort.Slice(transactions, func(i, j int) bool {
		time1,_:=time.Parse(shorterm,transactions[i].Date)
		time2,_:=time.Parse(shorterm,transactions[j].Date)
		if time1.After(time2) {
			return true
		}
		return false
	})

}

func splitTransactions(transactions []plaid.Transaction) [][]plaid.Transaction{
	const shorterm = "2006-01-02"
	var splitedTransactions [][]plaid.Transaction
	var transactionsUnit []plaid.Transaction
	timeStart:=time.Now()
	for _,transaction := range transactions{

		date,_:=time.Parse(shorterm,transaction.Date)
		if timeStart.Sub(date)>=time.Hour*24*30{
			splitedTransactions=append(splitedTransactions, transactionsUnit)
			transactionsUnit =nil
			timeStart=date
		}
		transactionsUnit =append(transactionsUnit, transaction)

	}
	if transactionsUnit!=nil{
		splitedTransactions=append(splitedTransactions, transactionsUnit)
	}


	return splitedTransactions
}