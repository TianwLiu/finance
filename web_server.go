package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)
import "github.com/gin-contrib/cors"
var ginEngine *gin.Engine
var server *http.Server
func init()  {
	ginEngine=gin.Default()


	corsFigure:=cors.DefaultConfig()
	corsFigure.AllowAllOrigins=true
	corsFigure.AllowHeaders=append(corsFigure.AllowHeaders, "api_token")
	ginEngine.Use(cors.New(corsFigure))
	bindWebApi()

}


func webServerStart(hostAndPort string,crtFilename string,privateKey string){
	//ginEngine.Run()
	server=&http.Server{
		Addr:              hostAndPort,
		Handler:           ginEngine,
	}

	var errTLS,errHttp error
	go func() {
		time.Sleep(time.Second)
		if errTLS==nil{
			fmt.Println("Webserver[TLS] started finish , serving now at ",server.Addr)
		}else if errHttp==nil{
			fmt.Println("Webserver[Http] started finish , serving now at ",server.Addr)
		}
		return
	}()
	if errTLS=server.ListenAndServeTLS(crtFilename,privateKey);errTLS!=http.ErrServerClosed{
		println("Webserver[TLS] started fail, error:",errTLS.Error())
		println("Trying to start as http")

		if errHttp=server.ListenAndServe();errHttp!=http.ErrServerClosed{
			println("Webserver[Http] started fail, error:",errHttp.Error())
		}
	}


}

func webServerShutdown(){
	if server==nil{
		fmt.Println("-->web server not running, no need to close")
	}else{
		fmt.Println("-->closing web server")
		if err:=server.Shutdown(context.Background());err!=nil{
			fmt.Printf("-->web server shutting down%v",err)
		}
		fmt.Println("-->webserver closed")
	}
}

func bindWebApi()  {
	ginEngine.GET("/ping", func(c *gin.Context) {
		c.JSON(200,"the server is running")
	})
	ginEngine.POST("/register",register)
	ginEngine.POST("/logIn",logIn)
	auth := ginEngine.Group("/auth",TokenAuthMiddleware())
	{

		auth.POST("/getLinkToken", getLinkToken)
		auth.POST("/postPublicToken", postPublicToken)

		auth.GET("/accounts/get",getBalances)
		auth.GET("/transactions/get", getTransactions)

		auth.GET("/group/accounts/get",getGroupBalances)
		auth.GET("/group/transactions/get", getGroupTransactions)

		auth.POST("/group/addMember",addMember)
		auth.POST("/group/delMember",delMember)
	}


}