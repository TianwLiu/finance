finance

1.0 server's api is defined as following

//register
//"http://localhost:8080/register",
$.ajax(
        "http://localhost:8080/register",
        {
            type:"post",
            data:{"user_id": userId,
                "pass_word":passWord},
            success:function (result){
                window.alert("register success");
            }

        }
    );

//logIn
http://localhost:8080/logIn
$.ajax(
        "http://localhost:8080/logIn",
        {
            type:"post",
            data:{"user_id": userId,
                "pass_word":passWord},
            success:function (result){
                if(result.flag){
                    localStorage.userId=result.user_id
                    localStorage.groupId=result.group_id
                    localStorage.token=result.token
                    console.log(localStorage.token);

                    location.href="./home.html"
                }else{
                    window.alert("user_id or password wrong");
                }

            }

        }
    );



//following api don't need user_id or group_id anymore, but you can use it as a way to decevie the hacker
//return [one] cashFlow and [one] sum of transactions of accesstokens(all online accounts of all banks) within that user_id and type: json
/auth/transactions/get
{
 cash_flow: {
    Income: 602.74,
    Expense: 1228.35,
    Assest: 0,
    Liability: 0,
    Income_sheet: { Work_income: 0, Passive_income: 0 },
    Indicators: {
      R_finance_accumulation: -625.6099999999999,
      R_current_interest: 0,
      R_finance_health: 0.4906907640330525,
      R_finance_indepence: 0
    }
  },
  
 transactions:[
    {
      account_id: '8eK70YNEvVuqeB6aLDRwfjnrKx5nnYIyd81Kx',
      amount: 21.04,
      iso_currency_code: 'USD',
      unofficial_currency_code: '',
      category: [Array],
      category_id: '16001000',
      date: '2021-01-19',
      authorized_date: '',
      location: [Object],
      merchant_name: '',
      name: 'BANK OF AMERICA CREDIT CARD Bill Payment',
      payment_meta: [Object],
      payment_channel: 'other',
      pending: false,
      pending_transaction_id: '',
      account_owner: '',
      transaction_id: 'xbpBqnvAeKcdXYNRokg7fNoMdY70qbUMeqRjK',
      transaction_type: 'special',
      transaction_code: ''
    },
	]
	
}

//return one [cashFlow] for [one] list of transactions and one [sum] of transactions of accesstokens(all online accounts of all banks) within user_ids of given group_id and type: json
//combine all user_id's transactions into one transactions and get the indicators of that transactions
/auth/group/transactions/get
{
 cash_flow: {
    Income: 602.74,
    Expense: 1228.35,
    Assest: 0,
    Liability: 0,
    Income_sheet: { Work_income: 0, Passive_income: 0 },
    Indicators: {
      R_finance_accumulation: -625.6099999999999,
      R_current_interest: 0,
      R_finance_health: 0.4906907640330525,
      R_finance_indepence: 0
    }
  },
  
 transactions:[
    {
      account_id: '8eK70YNEvVuqeB6aLDRwfjnrKx5nnYIyd81Kx',
      amount: 21.04,
      iso_currency_code: 'USD',
      unofficial_currency_code: '',
      category: [Array],
      category_id: '16001000',
      date: '2021-01-19',
      authorized_date: '',
      location: [Object],
      merchant_name: '',
      name: 'BANK OF AMERICA CREDIT CARD Bill Payment',
      payment_meta: [Object],
      payment_channel: 'other',
      pending: false,
      pending_transaction_id: '',
      account_owner: '',
      transaction_id: 'xbpBqnvAeKcdXYNRokg7fNoMdY70qbUMeqRjK',
      transaction_type: 'special',
      transaction_code: ''
    },
	]
	
}


//return real time balances of accesstokens(all online accounts of all banks) within that user_id （non-real-time）
http://localhost:8080/auth/accounts/get?real_time=[false/true]
{
user_id:"user_id"
cash_base: { depository_balance: 2021.93, credit_liability: 37.4 },
accounts:[
	  {
		account_id: 'ZOq3A67oV5f9ObK8Q6jqHy8d9rL884UR0o498',
		balances: {
		  available: 283.51,
		  current: 16.49,
		  limit: 300,
		  iso_currency_code: 'USD',
		  unofficial_currency_code: ''
		},
		mask: '0840',
		name: 'Bank of America Cash Rewards Platinum Plus Mastercard',
		official_name: 'Bank of America Cash Rewards Platinum Plus Mastercard',
		subtype: 'credit card',
		type: 'credit',
		verification_status: ''
	  },
	  {
		account_id: '8eK70YNEvVuqeB6aLDRwfjnrKx5nnYIyd81Kx',
		balances: {
		  available: 2021.93,
		  current: 2021.93,
		  limit: 0,
		  iso_currency_code: 'USD',
		  unofficial_currency_code: ''
		},
		mask: '7744',
		name: 'Adv Plus Banking',
		official_name: 'Adv Plus Banking',
		subtype: 'checking',
		type: 'depository',
		verification_status: ''
	  },...
	]
}

//return real time balances of accesstokens(all online accounts of all banks) within that user_id (non-real-time)
http://localhost:8080/auth/group/accounts/get?real_time=[false/true]
{
group_id:"group_id",
cash_base: { depository_balance: 2021.93, credit_liability: 37.4 },
user_accounts_list:[
	{
	user_id:"user_id"
	cash_base: { depository_balance: 2021.93, credit_liability: 37.4 },
	accounts:[
		  {
			account_id: 'ZOq3A67oV5f9ObK8Q6jqHy8d9rL884UR0o498',
			balances: {
			  available: 283.51,
			  current: 16.49,
			  limit: 300,
			  iso_currency_code: 'USD',
			  unofficial_currency_code: ''
			},
			mask: '0840',
			name: 'Bank of America Cash Rewards Platinum Plus Mastercard',
			official_name: 'Bank of America Cash Rewards Platinum Plus Mastercard',
			subtype: 'credit card',
			type: 'credit',
			verification_status: ''
		  },
		  {
			account_id: '8eK70YNEvVuqeB6aLDRwfjnrKx5nnYIyd81Kx',
			balances: {
			  available: 2021.93,
			  current: 2021.93,
			  limit: 0,
			  iso_currency_code: 'USD',
			  unofficial_currency_code: ''
			},
			mask: '7744',
			name: 'Adv Plus Banking',
			official_name: 'Adv Plus Banking',
			subtype: 'checking',
			type: 'depository',
			verification_status: ''
		  }
		  ,...
		]	
		
	},...
	]

}

//postPublicToken
http://localhost:8080/auth/postPublicToken
$.ajax(
        "http://localhost:8080/postPublicToken",
        {
            type:"post",
            data:{"userId":userId,
                "publicToken":publicToken,
                "institutionName":metadata.institution.name}

        }
    );


//getLinkToken
http://localhost:8080/auth/getLinkToken
$.post(
        "http://localhost:8080/getLinkToken",
            {"user_id": user_id},
            openLinkToken
    );



//addGroupMember
http://localhost:8080/auth/group/addMember
  $.ajax(
        "http://localhost:8080/auth/group/addMember",
        {
            type:"post",
            data:{"group_id":groupId,
                "member_id": userId,
                "member_pass_word":passWord},
            success:function (result){
                if(result.flag){
                    location.reload();
                }else{
                    window.alert(result.error);
                }

            }

        }
    );


//delGroupMember
http://localhost:8080/auth/group/delMember
function delMember(memberId){
    $.ajax(
        "http://localhost:8080/auth/group/delMember",
        {
            type:"post",
            data:{
                "member_id": memberId,
                },
            success:function (result){
                if(result.flag){
                    location.reload();
                }else{
                    window.alert(result.error);
                }

            }

        }
    );
}

