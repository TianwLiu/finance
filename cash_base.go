package main

type CashBase struct {
	DepositoryBalance float64 `json:"depository_balance"`
	CreditLiability	float64 `json:"credit_liability"`

}

func NewCashBase(accounts []AccountPlus)  CashBase{

	var depositoryBalance,creditLiablity float64
	for _,account:=range accounts {
		if account.Type=="credit" {
			creditLiablity+=account.Balances.Current
		}else if account.Type=="depository"{
			depositoryBalance+=account.Balances.Current
		}
	}

	return CashBase{
		DepositoryBalance: depositoryBalance,
		CreditLiability:   creditLiablity,
	}

}

