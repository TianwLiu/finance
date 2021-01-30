package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)





type Credential struct {
	InstitutionName string `json:"institution_name"`
	AccessToken     string `json:"access_token"`
}

type WebUser struct {
	UserId      string       `json:"user_id"`
	UserPass    string       `json:"user_pass"`
	GroupId		string		`json:"group_id"`
	Credentials []Credential `json:"credentials"`
}

var user map[string][]string
var database DataRoot

func writeDataBase(){

	dataBaseRaw,_:=json.Marshal(database)
	fmt.Println(string(dataBaseRaw))
	err := ioutil.WriteFile("database.json", dataBaseRaw, 0644)
	if err!=nil{
		panic(err)
	}

	/*
	var f *os.File
	if  _,err:=os.Stat("database.json");os.IsNotExist(err){
		f,err = os.Create("database.json")//create it casue file not exist
	}else{
		f,err=os.OpenFile("database.json",os.O_APPEND,0666)
	}

	defer f.Close()

	n, err1 := io.WriteString(f, string(dataBaseRaw)) //写入文件(字符串)
	if err1 != nil {
		panic(err1)
	}
	fmt.Printf("写入 %d 个字节n", n)*/
}

func readDataBase(){
	/*dataBaseFile,err:=os.Open("database.json")
	if err!=nil{
		fmt.Println(err.Error())
		return
	}
	ioutil.ReadAll()*/
	dataBaseRaw,err:=ioutil.ReadFile("database.json")
	if err!=nil{
		fmt.Println(err.Error())
		return
	}

	err=json.Unmarshal(dataBaseRaw,&database)
	if err!=nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(database.WebUsers[0].UserPass)
	fmt.Println("sdsd")
}

func init(){
	//readDataBase()
	/*database = DataRoot{
		WebUsers: []WebUser{
			{
				UserId:      "tianwei",
				UserPass:    "passs",
				Credential: nil,
			},
		},
		WebGroups: nil,
	}*/
	user=map[string][]string{
		"one":[]string{"access-development-1aaf1246-aa55-4a00-9eb4-1a8663248289"},
		"hongqi":[]string{"access-development-52860d35-a07e-4f08-b0d5-3241508fb515"},
		"tianwei": []string{"access-development-ddb1ec2d-ab93-4923-86a6-0ee7fa79be44","access-development-1aaf1246-aa55-4a00-9eb4-1a8663248289"},
		"test": []string{"access-development-ddb1ec2d-ab93-4923-86a6-0ee7fa79be44", "tianwei-discover"},
		"guiqin":[]string{"access-development-e1c5cbdd-fae4-450b-a124-38b05fdc7356"},
	}

}

//delUser should only server for root user, not server for web api
func delUser(userId string) error {
	err:=delWebUser(userId)
	return err
}

//return nil if add new user successfully
func addNewUser(webUser WebUser) error{
	if checkUser(webUser.UserId){
		return errors.New("warning: new user trying to overwrite old user,execution cancel")
	}
	err:=updateWebUser(webUser)
	return err
}

//return true if the userId has already exist
func checkUser(userId string)  bool{
	webUser:=viewWebUser(userId)
	if webUser.UserId==userId{
		return true
	}else{
		return false
	}
}


//return true if pass was correct
func verifyPass(userId string,pass string)  bool{

	webUser:=viewWebUser(userId)
	if userId!=webUser.UserId {
		return false
	}
	if pass!= webUser.UserPass {
		return false
	}

	return true
}

//return nil if setUserPass successfully
func setUserPass(userId string,pass string)  error{
	webUser:=viewWebUser(userId)
	webUser.UserPass=pass
	err:=updateWebUser(webUser)
	return err
}

//return nil if addUserAccessToken successfully,
//this function will only add an account credential for the special userId
func addUserAccessToken(userId string,accessToken string,institutionName string) error{
	webUser:=viewWebUser(userId)
	webUser.Credentials = append(webUser.Credentials, Credential{
		InstitutionName: institutionName,
		AccessToken:     accessToken,
	})
	err:=updateWebUser(webUser)
	return err
}

func getUserAccessTokens(userId string)([]string,error){
	webUser:=viewWebUser(userId)
	var userAccessTokens []string
	for _,credential:= range webUser.Credentials {
		userAccessTokens = append(userAccessTokens, credential.AccessToken)
	}
	if userAccessTokens==nil {
		return userAccessTokens,errors.New("can't find the accessToken of specific user")
	}else{
		return userAccessTokens,nil
	}
}

/*func getUserAccessTokens(user_id string) ([]string,error){
	//fmt.Println(user[user_id])

	if user[user_id]==nil {
		return user[user_id],errors.New("can't find the accessToken of specific user")
	}else{
		return user[user_id],nil
	}
}*/
