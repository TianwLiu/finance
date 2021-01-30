package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 这里需要根据前端传递token传的方式使用不同的方式获取
		// 放在请求头(Request Header),使用c.Request.Header.Get("[前端设置的key]") 【注意将自定义的key放在header中时，后端是否进行了跨域资源访问的控制】
		// 放在get请求链接中或post请求的参数中，使用c.Query("[前端设置的key]")
		tokenStr := c.Request.Header.Get("api_token")
		/*if tokenStr == "" {
			// 尝试从链接中获取 兼容下部分特殊get请求
			tokenStr = c.Query("api_token")
			if tokenStr == "" {
				c.JSON(http.StatusUnauthorized, gin.H{
					"success": "false",
					"msg":     "用户未登陆",
				})
				c.Abort()
				return
			}
		}*/

		// 解析token
		claims, err := ParseToken(tokenStr,c.ClientIP())
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": "false",
				"msg":     "token 解析失败",
			})
			c.Abort()
			return
		}
		// 注意这里数值类型断言后的类型位float64
		//c.Set("uid", int(claims["uid"].(float64)))
		c.Set("id",claims["id"])

		c.Next()
	}
}
