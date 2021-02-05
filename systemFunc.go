package main

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"syscall"
	"time"
)

type SystemConf struct {
	PassMd5                []byte
	PlaidClientID          []byte
	PlaidSandboxSecret     []byte
	PlaidDevelopmentSecret []byte
}


func startDaemon(systemPass string,ENV string,hostAndPort string,serverCrt string,serverPrivateKey string){

	if len(systemPass)<8{
		systemPass=systemPass+systemPass[:8-len(systemPass)]

	}

	systemConf,err:=viewSystemConf()
	if err!=nil{
		fmt.Println(err)
		panic(err)
	}
	if len(systemConf.PassMd5) ==0{
		panic("system password is nil, please check")
	}
	passMd5:=md5.Sum([]byte(systemPass))
	if !bytes.Equal(passMd5[:],systemConf.PassMd5) {
		panic("system password not match, start fail")
	}
	closeDatabase()

	var logFile *os.File
	if  _,err:=os.Stat("daemon_mode_child_process.log");os.IsNotExist(err){
		logFile,err = os.Create("daemon_mode_child_process.log")
	}else{
		logFile,err= os.OpenFile("daemon_mode_child_process.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	}


	for true{
		childCmd=exec.Cmd{
			Path:         os.Args[0],
			Args:         []string{
				os.Args[0],

				"-p",
				args.systemPass,
				"-r",
				"-e",
				args.env,
				"-crt",
				args.crtFilePath,
				"-key",
				args.privateKeyPath,
			},
			Env:          nil,
			Dir:          "",
			Stdin:        nil,
			Stdout:       logFile,
			Stderr:       logFile,
			ExtraFiles:   nil,
			SysProcAttr:  nil,
			Process:      nil,
			ProcessState: nil,
		}
		errs:=childCmd.Start()
		if errs!=nil{
			fmt.Println(err.Error())
			break
		}
		fmt.Println("current pid is",os.Getpid())
		fmt.Println("child pid is running as ",childCmd.Process.Pid)

		if errChild:=childCmd.Wait();errChild==nil{
			fmt.Println("child process has exit normally,now exit daemon process")
			break
		}


	}
	systemExit()


}
func systemStart(systemPass string,ENV string,hostAndPort string,serverCrt string,serverPrivateKey string) {
	defer func() {
		if err:=recover();err!=nil{
			fmt.Println(err)

		}
	}()

	if len(systemPass)<8{
		systemPass=systemPass+systemPass[:8-len(systemPass)]

	}

	systemConf,err:=viewSystemConf()
	if err!=nil{
		fmt.Println(err)
		panic(err)
	}
	if len(systemConf.PassMd5) ==0{
		panic("system password is nil, please check")
	}
	passMd5:=md5.Sum([]byte(systemPass))
	if !bytes.Equal(passMd5[:],systemConf.PassMd5) {
		panic("system password not match, start fail")
	}

	passHash:=md5.Sum([]byte(systemPass))
	key:= append([]byte(systemPass),passHash[len([]byte(systemPass))-8:]...)

	plaidClientId:=string(AesDecryptCBC([]byte(systemConf.PlaidClientID),key))
	plaidDevelopmentSecret:=string(AesDecryptCBC([]byte(systemConf.PlaidDevelopmentSecret),key))
	plaidSandboxSecret:=string(AesDecryptCBC([]byte(systemConf.PlaidSandboxSecret),key))


	//setup jwt
	jwtReady(systemPass)

	//setup plaid according according to env
	if ENV==ENV_DEVELOPMENT {
		plaidReady(ENV_DEVELOPMENT, plaidClientId, plaidDevelopmentSecret)
	}else if ENV==ENV_SANDBOX{
		plaidReady(ENV_SANDBOX,plaidClientId,plaidSandboxSecret)
	}else{
		panic("got unknown environment of plaid, plaid setup fail\n")
	}

	webServerStart(hostAndPort,serverCrt,serverPrivateKey)




}
func systemExit(){
	if childCmd.Process==nil{

		fmt.Println("no child process to terminate")

	}else{
		switch runtime.GOOS {
		case "windows":
			fmt.Println("[warning] current os is windows, please check zombie process by yourself")
		case "linux":
			for true{

				if childCmd.ProcessState!=nil&&childCmd.ProcessState.Exited(){
					fmt.Println("[PID]:",childCmd.ProcessState.Pid(),"child process exited")
					break
				}else{
					if err:=childCmd.Process.Signal(syscall.SIGTERM);err!=nil{
						fmt.Println("terminate child process fail，error："+err.Error())
						break
					}else{
						fmt.Println("sent SIGTERM to Process-PID:",childCmd.Process.Pid,",now checking")
					}
				}

				time.Sleep(time.Second)
			}

		}

	}

	fmt.Println("system: shutting down")
	webServerShutdown()
	closeDatabase()
	fmt.Println("system: good byte!")

	os.Exit(0)
}


func systemSetup(systemPass string,plaidClientID string,plaidSandboxSecret string,plaidDevelopmentSecret string) error{

	if len(systemPass)<8{
		systemPass=systemPass+systemPass[:8-len(systemPass)]
	}
	fmt.Println(systemPass)
	if len(systemPass)>24{
		return errors.New("systemPass too long")
	}

	//passMd5 length is 16 byte, compose part of pass's md5 and pass itself to generate the key of 24bytes length
	passMd5 :=md5.Sum([]byte(systemPass))
	key:= append([]byte(systemPass), passMd5[len([]byte(systemPass))-8:]...)

	var systemConf SystemConf
	systemConf.PassMd5=passMd5[:]
	systemConf.PlaidClientID= AesEncryptCBC([]byte(plaidClientID), key)
	systemConf.PlaidDevelopmentSecret=AesEncryptCBC([]byte(plaidDevelopmentSecret),key)
	systemConf.PlaidSandboxSecret=AesEncryptCBC([]byte(plaidSandboxSecret),key)

	err:=updateSystemConf(systemConf)

	return err

}

func systemCheck(systemPass string)  (string, string,string,string){
	systemConf,err:=viewSystemConf()
	if err!=nil{
		fmt.Println(err)
		panic(err)
	}
	if len(systemPass)<8{
		systemPass=systemPass+systemPass[:8-len(systemPass)]
	}

	if len(systemPass)>24{
		return "system password style not right","","",""
	}
	passMd5:=md5.Sum([]byte(systemPass))
	if !bytes.Equal(passMd5[:],systemConf.PassMd5) {
		return "system password not match, start fail","","",""
	}

	passHash:=md5.Sum([]byte(systemPass))
	key:= append([]byte(systemPass),passHash[len([]byte(systemPass))-8:]...)

	plaidClientId:=string(AesDecryptCBC([]byte(systemConf.PlaidClientID),key))
	plaidDevelopmentSecret:=string(AesDecryptCBC([]byte(systemConf.PlaidDevelopmentSecret),key))
	plaidSandboxSecret:=string(AesDecryptCBC([]byte(systemConf.PlaidSandboxSecret),key))

	return "",plaidClientId,plaidDevelopmentSecret,plaidSandboxSecret

}

func systemConfShow() string{
	systemConf,err:=viewSystemConf()
	if err!=nil{
		fmt.Println(err)
		return err.Error()
	}
	jsonByteConf,_:=json.Marshal(systemConf)
	return string(jsonByteConf)
}



//only for debug
/*func systemTry(systemPass string)  (string, string,string,string){
	systemConf,err:=viewSystemConf()
	if err!=nil{
		fmt.Println(err)
		panic(err)
	}
	if len(systemPass)<8{
		systemPass=systemPass+systemPass[:8-len(systemPass)]

	}

	if len(systemPass)>24{
		return "system password style not right","","",""
	}


	passHash:=md5.Sum([]byte(systemPass))
	key:= append([]byte(systemPass),passHash[len([]byte(systemPass))-8:]...)

	plaidClientId:=string(AesDecryptCBC([]byte(systemConf.PlaidClientID),key))
	plaidDevelopmentSecret:=string(AesDecryptCBC([]byte(systemConf.PlaidDevelopmentSecret),key))
	plaidSandboxSecret:=string(AesDecryptCBC([]byte(systemConf.PlaidSandboxSecret),key))

	return "",plaidClientId,plaidDevelopmentSecret,plaidSandboxSecret

}*/

