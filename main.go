package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main()  {
	defer func() {
		startExit()
	}()

	flag.Parse()
	switch {
	case args.help:
		flag.Usage()
		return
	case args.show:
			menuShow()
		return
	case args.check:
			menuCheck()
		return
	case args.setup:
			menuSetUp()
		return
	case args.start:
		if args.systemPass!=""&&args.env!=""{
			go func() {
				systemStart(args.systemPass,args.env,args.hostAndPort,args.crtFilePath,args.privateKeyPath)
			}()

			listenSignal()
		}else{
			fmt.Println("wrong args, please check")
			flag.Usage()
			return
		}
	default:
		flag.Usage()
	}

}

func listenSignal(){
	c:=make(chan os.Signal,4)
	signal.Notify(c,syscall.SIGINT,syscall.SIGTERM,syscall.SIGQUIT,syscall.SIGHUP)

	for sig:=range c {
		switch sig {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP:
			fmt.Printf("get Signal as %v", sig)
			fmt.Println(", system enter exit procedure")
			startExit()
		}
	}

}

func startExit(){
	fmt.Println("system: shutting down")


	if server==nil{
		fmt.Println("-->web server not running, no need to close")
	}else{
		fmt.Println("-->closing web server")
		if err:=server.Shutdown(context.Background());err!=nil{
			fmt.Printf("-->web server shutting down%v",err)
		}
		fmt.Println("-->webserver closed")
	}


	if db==nil{
		fmt.Println("-->database not running, no need to close")
	}else{
		fmt.Println("-->closing database bolt")
		if err:=db.Close();err!=nil{
			fmt.Println(err.Error())
		}else{
			fmt.Println("-->database blot closed")
		}
	}



	fmt.Println("system: good byte!")

	os.Exit(0)
}