package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main()  {
	defer func() {
		systemExit()
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
			systemExit()
		}
	}

}

