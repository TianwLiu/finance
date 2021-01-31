package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
)
var db *bolt.DB
const webGroupsBucketName="webGroups"
const webUsersBucketName="webUsers"
const SYSTEM_BUCKET="system"
const SYSTEM_CONF="systemConf"

func init()  {
	var err error
	db,err=bolt.Open("finance.db",0666,nil)
	if err!=nil{
		panic(err)
	}
	if err=backUpDatabase("finance_last_start.db.bak");err!=nil{
		fmt.Println("database back up of [start stage] failed",err.Error())
	}

}

func closeDatabase(){
	if db==nil{
		fmt.Println("-->database not running, no need to close")
	}else{

		fmt.Println("----->database backing up")
		if err:=backUpDatabase("finance_last_close.db.bak");err!=nil{
			fmt.Println(err.Error())
		}else{
			fmt.Println("-->database back up finish")
		}
		fmt.Println("-->database closing")
		if err:=db.Close();err!=nil{
			fmt.Println(err.Error())
		}else{
			fmt.Println("-->database closed")
		}
	}
}

func backUpDatabase(bakFilePath string )error{

	return db.View(func(tx *bolt.Tx) error {
		return tx.CopyFile(bakFilePath,0666)

	})

}

func updateSystemConf(systemConf SystemConf)  error {

	systemConfBuff,_:=json.Marshal(systemConf)
	err:=db.Update(func(tx *bolt.Tx) error {
		systemBucket,err:=tx.CreateBucketIfNotExists([]byte(SYSTEM_BUCKET))
		if err!=nil{
			return err
		}
		return systemBucket.Put([]byte(SYSTEM_CONF),systemConfBuff)
	})
	return err
}

func viewSystemConf()  (SystemConf,error){
	var systemConf SystemConf
	err:=db.View(func(tx *bolt.Tx) error {
		systemBucket:=tx.Bucket([]byte(SYSTEM_BUCKET))
		if systemBucket==nil{
			return errors.New("no system information found, please set system first.")
		}
		systemConfByte:=systemBucket.Get([]byte(SYSTEM_CONF))
		if systemConfByte==nil{
			return errors.New("no system information found, please set system first.")
		}
		return json.Unmarshal(systemConfByte,&systemConf)


	})
	return systemConf,err
}


//recommend to use this api only in user.go, not call this directly
func delWebUser(webUserId string) error {



	err:=db.Update(func(tx *bolt.Tx) error {
		webUsersBucket :=tx.Bucket([]byte("webUsers"))
		if webUsersBucket ==nil{
			return errors.New("warning:trying to del non-exist user")
		}

		err:=webUsersBucket.Delete([]byte(webUserId))
		return err
	})

	return err
}
func updateWebUser(webUser WebUser)  error{

	webUserKeyBuff:=[]byte(webUser.UserId)
	webUserBuff,_:=json.Marshal(webUser)
	err:=db.Update(func(tx *bolt.Tx) error {
		webUsersBucket,err:=tx.CreateBucketIfNotExists([]byte("webUsers"))
		if err!=nil{
			return err
		}
		err=webUsersBucket.Put(webUserKeyBuff,webUserBuff)
		return err
	})
	return err
}
func viewWebUser(webUserID string) WebUser {

	var webUserResult WebUser


	//defer db.Close()
	db.View(func(tx *bolt.Tx) error {
		webUsersBucket :=tx.Bucket([]byte("webUsers"))
		if webUsersBucket ==nil{
			return nil
		}

		webUserBuff:= webUsersBucket.Get([]byte(webUserID))

		if json.Unmarshal(webUserBuff,&webUserResult)!=nil{
			return nil
		}
		return nil
	})

	return webUserResult
}
func viewAllWebUser()  []WebUser{

	var webUserList []WebUser


	db.View(func(tx *bolt.Tx) error {
		webUsersBucket :=tx.Bucket([]byte("webUsers"))
		if webUsersBucket ==nil{
			return errors.New("webUsers Bucket not exist")
		}

		cursor:= webUsersBucket.Cursor()
		for webUserKeyBuff,webUserBuff:=cursor.First();webUserKeyBuff!=nil;webUserKeyBuff,webUserBuff=cursor.Next(){
			var webUser WebUser
			if json.Unmarshal(webUserBuff,&webUser)!=nil{
				return nil
			}
			webUserList = append(webUserList,webUser )
		}

		return nil
	})

	return webUserList
}

//recommend to use this api only in group.go, not call this directly
func delWebGroup(webGroupId string) error {


	err:=db.Update(func(tx *bolt.Tx) error {
		webGroupsBucket,err:=tx.CreateBucketIfNotExists([]byte("webGroups"))
		if err!=nil{
			return err
		}
		err= webGroupsBucket.Delete([]byte(webGroupId))
		return err
	})
	return err
}
func updateWebGroup(webGroup WebGroup) error{

	webGroupKeyBuff :=[]byte(webGroup.GroupId)
	webGroupBuff,_:=json.Marshal(webGroup)
	err:=db.Update(func(tx *bolt.Tx) error {
		webGroupsBucket,err:=tx.CreateBucketIfNotExists([]byte("webGroups"))
		if err!=nil{
			return err
		}
		err= webGroupsBucket.Put(webGroupKeyBuff, webGroupBuff)
		return err
	})
	return err
}
func viewWebGroup(webGroupId string) WebGroup{
	var webGroupResult WebGroup

	db.View(func(tx *bolt.Tx) error {
		webGroupsBucket :=tx.Bucket([]byte("webGroups"))
		if webGroupsBucket ==nil{
			return nil
		}

		webGroupBuff := webGroupsBucket.Get([]byte(webGroupId))

		if json.Unmarshal(webGroupBuff,&webGroupResult)!=nil{
			return nil
		}
		return nil
	})

	return webGroupResult
}
func viewAllWebGroup() []WebGroup {
	var webGroupList []WebGroup


	db.View(func(tx *bolt.Tx) error {
		webGroupsBucket :=tx.Bucket([]byte("webGroups"))
		if webGroupsBucket ==nil{
			return nil
		}

		cursor:= webGroupsBucket.Cursor()
		for webGroupKeyBuff, webGroupBuff :=cursor.First(); webGroupKeyBuff !=nil; webGroupKeyBuff, webGroupBuff =cursor.Next(){
			var webGroup WebGroup
			if json.Unmarshal(webGroupBuff,&webGroup)!=nil{
				return nil
			}
			webGroupList = append(webGroupList, webGroup)
		}

		return nil
	})

	return webGroupList
}