package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
)

func menuShow(){
	fmt.Println("System config:\n-->:",systemConfShow())

	fmt.Println("Web user info:")
	for _,webUser :=range viewAllWebUser(){
		fmt.Println("-->:",webUser)
	}
	fmt.Println("Web group info")
	for _,webGroup:= range viewAllWebGroup(){
		fmt.Println("-->:",webGroup)
	}
}
func menuCheck(){
	if args.systemPass!=""{
		fmt.Println("Check successful")
		fmt.Println(systemCheck(args.systemPass))
	}else{
		fmt.Println("Wrong , please check")
		flag.Usage()

	}
}

func menuSetUp(){
	var input string
	for {

		fmt.Println("Please input: show(to show system config), \n system, webuser, webgroup to setup,\n del to delete webuser or webgroup\n input quit to exit")
		fmt.Scanln(&input)
		switch input {
		case "show":
			menuShow()
		case "system":
			var password,plaidClientId,plaidClientDevelopmentSecret,plaidClientSandboxSecret string
			fmt.Println("Please input your system password, max length 24 byte:")
			fmt.Scanln(&password)
			fmt.Println("Please input your plaid client ID:")
			fmt.Scanln(&plaidClientId)
			fmt.Println("Please input your plaid Development Secret")
			fmt.Scanln(&plaidClientDevelopmentSecret)
			fmt.Println("Please input your plaid SandBox Secret")
			fmt.Scanln(&plaidClientSandboxSecret)
			err:=systemSetup(password,plaidClientId,plaidClientSandboxSecret,plaidClientDevelopmentSecret)
			if err!=nil{
				fmt.Println(err.Error())
			}else{
				fmt.Println("Setup successfully")
			}
		case"webuser":
			var webUser WebUser
			var credentials []Credential
			var unencryptedWebUserPass string
			fmt.Println("Please input webUserID:")
			fmt.Scanln(&webUser.UserId)
			webUser.GroupId=webUser.UserId
			fmt.Println("Please input webUserPass:")
			fmt.Scanln(&unencryptedWebUserPass)

			fmt.Scanln(webUser.Credentials)
			md5Pass:=md5.Sum([]byte(unencryptedWebUserPass))
			webUser.UserPass=hex.EncodeToString(md5Pass[:])
			for {
				fmt.Println("You can input infinite credentials\n,please input institution name or  input quit to exit")
				fmt.Scanln(&input)
				if input=="quit"{
					break
				}
				credential:=Credential{InstitutionName: input}
				fmt.Println("Please input the accessToken for the institution")
				fmt.Scanln(&credential.AccessToken)
				credentials=append(credentials, credential)
			}
			webUser.Credentials=credentials
			fmt.Println("Current webuser is:\n",webUser)
			fmt.Println("Input ok to write to database(overwrite if exist), input others to give up")
			fmt.Scanln(&input)
			if input=="ok"{
				err:=updateWebUser(webUser)
				if err!=nil{
					fmt.Println(err.Error())
				}else{
					fmt.Println("Setup successfully")
				}
			}else{
				fmt.Println("You gave up")
			}

		case "webgroup":
			var webGroup WebGroup

			fmt.Println("Please input groupID:")
			fmt.Scanln(&webGroup.GroupId)

			for {
				fmt.Println("You can input infinite members Ids \n,please input institution name or  input quit to exit")
				fmt.Scanln(&input)
				if input=="quit"{
					break
				}
				webGroup.MemberUserIds=append(webGroup.MemberUserIds, input)

			}
			fmt.Println("Current webuser is:\n",webGroup)
			fmt.Println("Input ok to write to database(overwrite if exist), input others to give up")
			fmt.Scanln(&input)
			if input=="ok"{
				err:=updateWebGroup(webGroup)
				if err!=nil{
					fmt.Println(err.Error())
				}else{
					fmt.Println("Setup successfully")
				}
			}else{
				fmt.Println("You gived up")
			}
		case "del":
			fmt.Println("Please input webuser or webgroup to enter del procedure ")
			fmt.Scanln(&input)
			switch input {
			case "webuser":
				fmt.Println("Please input the user id you want to del ")
				var id string
				fmt.Scanln(&id)
				fmt.Println("You're going to del\n-->:",viewWebUser(id),"\n input ok to del,others to give up")
				fmt.Scanln(&input)
				if input=="ok"{
					err:=delWebUser(id)
					if err!=nil{
						fmt.Println(err.Error())
					}else{
						fmt.Println("Del successfully")
					}
				}else{
					fmt.Println("You gave up")
				}
			case"webgroup":
				fmt.Println("Please input the group id you want to del ")
				var id string
				fmt.Scanln(&id)
				fmt.Println("You're going to del\n-->:",viewWebGroup(id),"\n input ok to del,others to give up")
				fmt.Scanln(&input)
				if input=="ok"{
					err:=delWebGroup(id)
					if err!=nil{
						fmt.Println(err.Error())
					}else{
						fmt.Println("Del successfully")
					}
				}else{
					fmt.Println("You gave up")
				}

			}
		case "quit":
			fmt.Println("Good bye")
			os.Exit(0)
		}
	}


}