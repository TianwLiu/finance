package main

import (
	"log"
	"os"
)



func printlnLog(message ...string)  {
	var logFile *os.File
	if  _,err:=os.Stat("finance.log");os.IsNotExist(err){
		logFile,err = os.Create("finance.log")
	}else{
		logFile,err= os.OpenFile("finance.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	}
	defer logFile.Close()
	logger := log.New(logFile, "", log.Lshortfile|log.Ldate|log.Ltime)
	logger.Println(message)
}