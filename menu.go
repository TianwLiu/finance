package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
)

func menuShow(){
	fmt.Println("system config:\n-->:",systemConfShow())

	fmt.Println("web user info:")
	for _,webUser :=range viewAllWebUser(){
		fmt.Println("-->:",webUser)
	}
	fmt.Println("web group info")
	for _,webGroup:= range viewAllWebGroup(){
		fmt.Println("-->:",webGroup)
	}
}
func menuCheck(){
	if args.systemPass!=""{
		fmt.Println(systemCheck(args.systemPass))
	}else{
		fmt.Println("wrong args, please check")
		flag.Usage()

	}
}

func menuSetUp(){
	var input string
	for {

		fmt.Println("please input: show(to show system config), \n system, webuser, webgroup to setup,\n del to delete webuser or webgroup\n input quit to exit")
		fmt.Scanln(&input)
		switch input {
		case "show":
			menuShow()
		case "system":
			var password,plaidClientId,plaidClientDevelopmentSecret,plaidClientSandboxSecret string
			fmt.Println("pleas input your system password, max length 24 byte:")
			fmt.Scanln(&password)
			fmt.Println("pleas input your plaid client ID:")
			fmt.Scanln(&plaidClientId)
			fmt.Println("pleas input your plaid Development Secret")
			fmt.Scanln(&plaidClientDevelopmentSecret)
			fmt.Println("pleas input your plaid SandBox Secret")
			fmt.Scanln(&plaidClientSandboxSecret)
			err:=systemSetup(password,plaidClientId,plaidClientSandboxSecret,plaidClientDevelopmentSecret)
			if err!=nil{
				fmt.Println(err.Error())
			}else{
				fmt.Println("setup successfully")
			}
		case"webuser":
			var webUser WebUser
			var credentials []Credential
			var unencryptedWebUserPass string
			fmt.Println("pleas input webUserID:")
			fmt.Scanln(&webUser.UserId)
			webUser.GroupId=webUser.UserId
			fmt.Println("pleas input webUserPass:")
			fmt.Scanln(&unencryptedWebUserPass)

			fmt.Scanln(webUser.Credentials)
			md5Pass:=md5.Sum([]byte(unencryptedWebUserPass))
			webUser.UserPass=hex.EncodeToString(md5Pass[:])
			for {
				fmt.Println("you can input infinite credentials\n,please input institution name or  input quit to exit")
				fmt.Scanln(&input)
				if input=="quit"{
					break
				}
				credential:=Credential{InstitutionName: input}
				fmt.Println("please input the accessToken for the institution")
				fmt.Scanln(&credential.AccessToken)
				credentials=append(credentials, credential)
			}
			webUser.Credentials=credentials
			fmt.Println("current webuser is:\n",webUser)
			fmt.Println("input ok to write to database(overwrite if exist), input others to give up")
			fmt.Scanln(&input)
			if input=="ok"{
				err:=updateWebUser(webUser)
				if err!=nil{
					fmt.Println(err.Error())
				}else{
					fmt.Println("setup successfully")
				}
			}else{
				fmt.Println("you gived up")
			}

		case "webgroup":
			var webGroup WebGroup

			fmt.Println("pleas input groupID:")
			fmt.Scanln(&webGroup.GroupId)

			for {
				fmt.Println("you can input infinite members Ids \n,please input institution name or  input quit to exit")
				fmt.Scanln(&input)
				if input=="quit"{
					break
				}
				webGroup.MemberUserIds=append(webGroup.MemberUserIds, input)

			}
			fmt.Println("current webuser is:\n",webGroup)
			fmt.Println("input ok to write to database(overwrite if exist), input others to give up")
			fmt.Scanln(&input)
			if input=="ok"{
				err:=updateWebGroup(webGroup)
				if err!=nil{
					fmt.Println(err.Error())
				}else{
					fmt.Println("setup successfully")
				}
			}else{
				fmt.Println("you gived up")
			}
		case "del":
			fmt.Println("please input webuser or webgroup to enter del procedure ")
			fmt.Scanln(&input)
			switch input {
			case "webuser":
				fmt.Println("please input the user id you want to del ")
				var id string
				fmt.Scanln(&id)
				fmt.Println("you're going to del\n-->:",viewWebUser(id),"\n input ok to del,others to give up")
				fmt.Scanln(&input)
				if input=="ok"{
					err:=delWebUser(id)
					if err!=nil{
						fmt.Println(err.Error())
					}else{
						fmt.Println("del successfully")
					}
				}else{
					fmt.Println("you gived up")
				}
			case"webgroup":
				fmt.Println("please input the group id you want to del ")
				var id string
				fmt.Scanln(&id)
				fmt.Println("you're going to del\n-->:",viewWebGroup(id),"\n input ok to del,others to give up")
				fmt.Scanln(&input)
				if input=="ok"{
					err:=delWebGroup(id)
					if err!=nil{
						fmt.Println(err.Error())
					}else{
						fmt.Println("del successfully")
					}
				}else{
					fmt.Println("you gived up")
				}

			}
		case "quit":
			fmt.Println("good bye")
			os.Exit(0)
		}
	}


}