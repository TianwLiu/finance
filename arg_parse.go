package main

import "flag"

type Args struct{
	help bool
	start bool
	setup bool
	check bool
	show bool

	systemPass string
	env string
	hostAndPort string
	crtFilePath string
	privateKeyPath string
}
var args Args

func init()  {

	flag.BoolVar(&args.help,"h",false,"show this help menu")
	flag.BoolVar(&args.start,"r",false,"run system and webServer")
	flag.BoolVar(&args.setup,"s",false,"setup all system conf, web user or web group")
	flag.BoolVar(&args.check,"v",false,"verify system")
	flag.BoolVar(&args.show,"l",false,"list database content")

	flag.StringVar(&args.hostAndPort,"addr",":10016","set the listen `address` for the server\n " +
		"ex: [127.0.0.1:8080], [:8080], [www.website.name:8080], default:[:10016]")
	flag.StringVar(&args.systemPass,"p","","system `password` ")
	flag.StringVar(&args.env,"e","","set plaid `environment value`: "+ENV_SANDBOX+" or "+ENV_DEVELOPMENT)
	flag.StringVar(&args.crtFilePath,"crt","","the path of server crt file")
	flag.StringVar(&args.privateKeyPath,"key","","the path of server private key")
}