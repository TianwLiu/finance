package main

import (
	"github.com/plaid/plaid-go/plaid"
)


/*functions within this file should only deal with task about cashflow, not balance.
cashflow definition: a flow of cash in a period of time, like transactions of a period
so it's highly possible that these functions only deal with transactions.
*/

type FinanceIndicator struct {

	R_finance_accumulation float64 `json:"R_finance_accumulation"`
	R_current_interest     float64 `json:R_current_interest`

	R_finance_health    float64 `json:"R_finance_health"`
	R_finance_indepence float64 `json:"R_finance_indepence"`

}

// warning: Income not include fluctuation of Assest value
// Income is based on per month
/*
generally, Work_income means job salary need to devote a lot of time and strength,
Passive_income means on the contrary
*/

type IncomeSheet struct{
	Work_income    float64 `json:"Work_income"`
	Passive_income float64 `json:"Passive_income"`
}

//class holding all data of one person
type Sheets struct {
	Income       float64          `json:"Income"`
	Expense      float64          `json:"Expense"`
	Assest       float64          `json:"Assest"`
	Liability    float64          `json:"Liability"`
	Income_sheet IncomeSheet      `json:"Income_sheet"`
	Indicators   FinanceIndicator `json:"Indicators"`
}

//build Sheets object
func NewSheets(transactions []plaid.Transaction,accounts []plaid.Account) Sheets {


	//get accountId of checking account, credit account
	//fmt.Println(transactionResp.Accounts)
	var accountIdListofCredit,accountIdListofDepository []string
	//accountIdListofDepository:=[]string
	for _,account:=range accounts {
		if account.Type=="credit" {
			accountIdListofCredit = append(accountIdListofCredit,account.AccountID )
		}else if account.Type=="depository"{
			accountIdListofDepository = append(accountIdListofDepository, account.AccountID)
		}
	}

	//split original transactions list into two: depository , credit
	var transactionsOfDepository,transactionsOfCredit []plaid.Transaction

	for _,transaction :=range transactions{
		var allocated bool
		for _,accountIdofDepository:= range accountIdListofDepository {
			if transaction.AccountID==accountIdofDepository {
				transactionsOfDepository = append(transactionsOfDepository, transaction)
				allocated=true
				break
			}
		}

		//continue if transaction have been added, cause one transactions only have one account type
		if allocated{ continue }

		for _,accountIdofCredit:=range accountIdListofCredit {
			if transaction.AccountID==accountIdofCredit {
				transactionsOfCredit = append(transactionsOfCredit, transaction)
				break
			}
		}
		//if allocated{ continue } not need to check allocated as the end of loop

	}

	//deleting internal transactions between different depository accounts for accurate cashFlow result
	length:= len(transactionsOfDepository)


	for i:=0;i<length-1;i++{
		if transactionsOfDepository[i].Code=="drop" {continue}
		for j:=i+1;j<length;j++{
			if transactionsOfDepository[j].Code=="drop" {continue}

			if (transactionsOfDepository[i].Amount+transactionsOfDepository[j].Amount)==0 &&
				(transactionsOfDepository[i].Category[0]=="Transfer")&&
				(transactionsOfDepository[j].Category[0]=="Transfer")&&
				((transactionsOfDepository[i].Category[1]=="Debit"&&transactionsOfDepository[j].Category[1]=="Credit") ||
						(transactionsOfDepository[j].Category[1]=="Debit"&&transactionsOfDepository[i].Category[1]=="Credit")){
				transactionsOfDepository[i].Code="drop"
				transactionsOfDepository[j].Code="drop"
				break
			}
		}
	}

	var incomeAmount, expenseAmount float64
	for _,transaction:=range transactionsOfDepository{

		if transaction.Code!="drop"{
			if transaction.Amount>0 {
				expenseAmount +=transaction.Amount
			}else {
				incomeAmount -=transaction.Amount
			}
		}

	}

	for _,transaction:=range transactionsOfCredit{
		expenseAmount +=transaction.Amount
	}

	sheets:= Sheets{
		Income:       incomeAmount,
		Expense:      expenseAmount,
		Assest:       0,
		Liability:    0,
		Income_sheet: IncomeSheet{},

	}
	sheets.Indicators =sheets.getIndicator()

	return sheets
}

//get all Indicators used to show in front end
func (sheets *Sheets)getIndicator() FinanceIndicator {
	return FinanceIndicator{
		sheets.Income -sheets.Liability -sheets.Expense,
		sheets.Income_sheet.Passive_income,
		sheets.Income /(sheets.Liability +sheets.Expense),
		sheets.Income_sheet.Passive_income /(sheets.Liability +sheets.Expense),
		}
}



