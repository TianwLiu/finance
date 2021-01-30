package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/plaid/plaid-go/plaid"
	"sort"
	"time"
)

//serve balances with user_id
//return date type refer to data_type.go file
func getBalances(webClient *gin.Context) {
	userId:=webClient.Keys["id"].(string)
	var realTime bool
	if webClient.Query("real_time")=="true"{
		realTime=true
	}else if webClient.Query("real_time")=="false" {
		realTime=false
	}


	accessTokens,err:= getUserAccessTokens(userId)
	if err!=nil {
		webClient.JSON(200,err.Error())
		return
	}

	// get balances for all accounts
	//balanceResp, err := plaidClient.GetBalances(accessTokens[0])
	var accounts []plaid.Account
	var errs []error
	for _,accessToken :=range accessTokens {

		if realTime {
			balanceResp, err := plaidClient.GetBalances(accessToken)
			accounts = append(accounts, balanceResp.Accounts...)
			errs = append(errs, err)
		}else{
			balanceResp, err := plaidClient.GetAccounts(accessToken)
			accounts = append(accounts, balanceResp.Accounts...)
			errs = append(errs, err)
		}



	}

	if accounts!=nil {
		webClient.JSON(200, UserAccounts{userId,NewCashBase(accounts), accounts})
	}
}

//serve group balances with user_id
//return date type refer to data_type.go file
func getGroupBalances(webClient *gin.Context) {

	groupId :=webClient.Keys["id"].(string)
	var realTime bool
	if webClient.Query("real_time")=="true"{
		realTime=true
	}else if webClient.Query("real_time")=="false" {
		realTime=false
	}

	userIds,err:=getGroupMemberID(groupId)
	if err!=nil{
		webClient.JSON(200,err.Error())
		return
	}

	//webClient.JSON(200,userIds)
	//var accessTokens []string

	var userAccountsList []UserAccounts
	var groupAccounts []plaid.Account
	for _, userId :=range userIds {


		//code dealing with one userID
		accessTokens, _ := getUserAccessTokens(userId)
		//accessTokens = append(accessTokens, accessTokensSub...)
		// get balances for all accounts
		//balanceResp, err := plaidClient.GetBalances(accessTokens[0])
		var accounts []plaid.Account
		var errs []error


		for _,accessToken :=range accessTokens {
			if realTime {
				balanceResp, err := plaidClient.GetBalances(accessToken)
				accounts = append(accounts, balanceResp.Accounts...)
				errs = append(errs, err)
			}else{
				balanceResp, err := plaidClient.GetAccounts(accessToken)
				accounts = append(accounts, balanceResp.Accounts...)
				errs = append(errs, err)
			}
		}
		groupAccounts = append(groupAccounts, accounts...)

		if accounts !=nil {
			userAccountsList = append(userAccountsList, UserAccounts{userId,NewCashBase(accounts) , accounts})
		}
	}

	webClient.JSON(200,GroupAccounts{
		GroupId:          groupId,
		CashBase:         NewCashBase(groupAccounts),
		UserAccountsList: userAccountsList,
	})

}


//serve transactions with user_id ,now it will return 3 month transactions
//return date type refer to data_type.go file
func getTransactions(webClient *gin.Context) {
	userId :=webClient.Keys["id"].(string)
	accessTokens,err:= getUserAccessTokens(userId)
	if err!=nil {
		webClient.JSON(200,err.Error())
		return
	}

	//refresh plaid transactions' record, current client didn't buy refresh product.
	/*
	responseRefresh, err := plaidClient.RefreshTransactions(accessTokens);
	if err!=nil  {
		fmt.Println(err)
		return
	}else{
		fmt.Println(responseRefresh.RequestID)
	}
	*/

	today:=time.Now().Format("2006-01-02")
	dayThreeMonthBefore:=time.Now().AddDate(0,-3,0).Format("2006-01-02")
	var accounts []plaid.Account
	var transactions []plaid.Transaction
	for _,accessToken := range accessTokens {
		response,_:=plaidClient.GetTransactions(accessToken,dayThreeMonthBefore,today)
		accounts = append(accounts, response.Accounts...)
		transactions = append(transactions, response.Transactions...)
	}

	const shorterm = "2006-01-02"
	sort.Slice(transactions, func(i, j int) bool {
		time1,_:=time.Parse(shorterm,transactions[i].Date)
		time2,_:=time.Parse(shorterm,transactions[j].Date)
		if time1.After(time2) {
			return true
		}
		return false
	})

	//response.RequestID="sd"
	//cashflow:=NewSheets(response)
	if err!=nil{
		webClient.JSON(200,err)
	}else{
		webClient.JSON(200,UserTransactions{
			CashFlow:     NewSheets(transactions,accounts),
			Transactions:transactions,
		})
	}
}

//serve group transactions with user_id ,now it will return 3 month transactions
//return date type refer to data_type.go file
func getGroupTransactions(webClient *gin.Context )  {
	groupId :=webClient.Keys["id"].(string)
	userIds,err:=getGroupMemberID(groupId)
	if err!=nil{
		webClient.JSON(200,err.Error())
		return
	}

	//webClient.JSON(200,userIds)
	var accessTokens []string
	for _, userId :=range userIds {
		accessTokensSub, _ := getUserAccessTokens(userId)
		accessTokens = append(accessTokens, accessTokensSub...)
	}



	//refresh plaid transactions' record, current client didn't buy refresh product.
	/*
		responseRefresh, err := plaidClient.RefreshTransactions(accessTokens);
		if err!=nil  {
			fmt.Println(err)
			return
		}else{
			fmt.Println(responseRefresh.RequestID)
		}
	*/
	if accessTokens==nil{
		webClient.JSON(200,gin.H{
			"flag":false,
			"error":"the group has no any accessTokens",
		})
		return
	}
	today:=time.Now().Format("2006-01-02")
	dayThreeMonthBefore:=time.Now().AddDate(0,-3,0).Format("2006-01-02")
	var accounts []plaid.Account
	var transactions []plaid.Transaction
	for _,accessToken := range accessTokens {
		response,_:=plaidClient.GetTransactions(accessToken,dayThreeMonthBefore,today)
		accounts = append(accounts, response.Accounts...)
		transactions = append(transactions, response.Transactions...)
	}

	const shorTerm = "2006-01-02"
	sort.Slice(transactions, func(i, j int) bool {
		time1,_:=time.Parse(shorTerm,transactions[i].Date)
		time2,_:=time.Parse(shorTerm,transactions[j].Date)
		if time1.After(time2) {
			return true
		}
		return false
	})

	//response.RequestID="sd"
	//cashflow:=NewSheets(response)
	if err!=nil{
		webClient.JSON(200,err.Error())
	}else{
		webClient.JSON(200,GroupTransactions{
			CashFlow:    NewSheets(transactions,accounts),
			Transactions: transactions,
		})
	}

}

func postPublicToken(webClient *gin.Context)  {
	userId:=webClient.Keys["id"].(string)
	publicToken:= webClient.PostForm("publicToken")
	institutionName:= webClient.PostForm("institutionName")
	accessToken,err:=getAccessToken(publicToken)
	if err!=nil{
		printlnLog("error:",err.Error(),"| userId:",userId,"| publicToken:",publicToken)
	}
	err=addUserAccessToken(userId,accessToken,institutionName)
	if err!=nil{
		printlnLog("error:",err.Error(),"| userId:",userId,"| accessToken",accessToken)
	}
}


//server user the linkToken to access plaid.com to get publicToken
//and user will give server the publicToken
//server will user publicToken to exchange accessToken
//which can authorize server to access user's online account forever
func getLinkToken(webClient *gin.Context) {
	userId:=webClient.Keys["id"].(string)

	linkTokenResponse,_:= plaidClient.CreateLinkToken(plaid.LinkTokenConfigs{
		User:&plaid.LinkTokenUser{
			ClientUserID:             userId,
			LegalName:                "",
			PhoneNumber:              "",
			EmailAddress:             "",
			PhoneNumberVerifiedTime:  time.Time{},
			EmailAddressVerifiedTime: time.Time{},
		},
		ClientName:            "finance",
		Products:              []string{"transactions"},
		AccessToken:           "",
		CountryCodes:          []string{"US"},
		Webhook:               "finance.trytolog.com",
		AccountFilters:        nil,
		CrossAppItemAdd:       nil,
		PaymentInitiation:     nil,
		Language:              "en",
		LinkCustomizationName: "",
		RedirectUri:           "",
		AndroidPackageName:    "",
	})
	webClient.JSON(200,gin.H{"userId":userId,"linkToken":linkTokenResponse.LinkToken})
}

func register(webClient *gin.Context)  {

	userId:=webClient.PostForm("user_id")
	passWord:=webClient.PostForm("pass_word")

	//user_id already exist
	if checkUser(userId){
		webClient.JSON(200,false)
		return
	}

	err:=addNewUser(WebUser{
		UserId:      userId,
		UserPass:    passWord,
		GroupId:     userId,
		Credentials: nil,
	})

	if err!=nil{
		printlnLog("add new user ----> error:",err.Error())
	}else{
		webClient.JSON(200,true)
	}

	err=addNewGroup(WebGroup{
		GroupId:       userId,
		MemberUserIds: []string{userId},
	})
	if err!=nil{
		printlnLog("add new group --->error:",err.Error())
	}

}
func logIn(webClient *gin.Context){
	userId:=webClient.PostForm("user_id")
	passWord:=webClient.PostForm("pass_word")

	if !verifyPass(userId,passWord){
		webClient.JSON(200,gin.H{"flag":false})
		return
	}
	groupId:=viewWebUser(userId).GroupId
	userIp:=webClient.ClientIP()
	token,_:=GenerateToken(userId,userIp)
	webClient.JSON(200,gin.H{
		"user_id":userId,
		"group_id":groupId,
		"flag":true,
		"token":token,
	})

}

func addMember(webClient *gin.Context)  {
	groupId:=webClient.Keys["id"].(string)

	memberId :=webClient.PostForm("member_id")
	memberPassWord :=webClient.PostForm("member_pass_word")

	if !verifyPass(memberId, memberPassWord){
		webClient.JSON(200,gin.H{
			"flag":false,
			"error":"new member's user id or password is wrong",
		})
		return
	}

	err:=addGroupMember(groupId, memberId)
	if err!=nil{
		webClient.JSON(200,gin.H{
			"flag":false,
			"error":err.Error(),
		})
		return
	}

	webClient.JSON(200,gin.H{
		"flag":true,
	})
	return

}
func delMember(webClient *gin.Context){
	groupId:=webClient.Keys["id"].(string)

	memberId :=webClient.PostForm("member_id")





	err:=delGroupMember(groupId,memberId)
	if err!=nil{
		webClient.JSON(200,gin.H{
			"flag":false,
			"error":err.Error(),
		})
		return
	}

	webClient.JSON(200,gin.H{
		"flag":true,
	})
	return
}





//internal test; called by shell.
func getLinkTokenTest(user_id string)  string{
	linkTokenReponse,_:= plaidClient.CreateLinkToken(plaid.LinkTokenConfigs{
		User:&plaid.LinkTokenUser{
			ClientUserID:             user_id,
			LegalName:                "",
			PhoneNumber:              "",
			EmailAddress:             "",
			PhoneNumberVerifiedTime:  time.Time{},
			EmailAddressVerifiedTime: time.Time{},
		},
		ClientName:            "finance",
		Products:              []string{"transactions"},
		AccessToken:           "",
		CountryCodes:          []string{"US"},
		Webhook:               "finance.trytolog.com",
		AccountFilters:        nil,
		CrossAppItemAdd:       nil,
		PaymentInitiation:     nil,
		Language:              "en",
		LinkCustomizationName: "",
		RedirectUri:           "",
		AndroidPackageName:    "",
	})
	return linkTokenReponse.LinkToken
}


//internal test; called by shell.
func removeItem(accessToken string) {
	response,err := plaidClient.RemoveItem(accessToken)
	if err!=nil {
	} else{
		fmt.Println(response)
	}
}

//internal test; called by shell.
func getIncomeSheet(accessToken string) *IncomeSheet {
	response, _ := plaidClient.GetIncome(accessToken)

	fmt.Println(response.Income)
	response1, err := plaidClient.GetTransactions(accessToken,"2020-12-01","2021-01-20")
	if err!=nil {
		fmt.Println(err)
	}else {
		fmt.Println(response1.Transactions)
	}

	responseAccounts,_:= plaidClient.GetAccounts(accessToken)
	fmt.Println(responseAccounts.Accounts)


	// OR - get auth for selected Accounts

	return &IncomeSheet{
		Work_income:    0,
		Passive_income: 0,
	}

}
