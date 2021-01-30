package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)
var SecretKey string


type JwtCustomClaims struct {
	jwt.StandardClaims // 包中自带的默认属性
	Id string `json:"id"` // 当前访问人uid 自定义添加一些自己需要的元素

}

func jwtReady(key string)  {
	SecretKey=key
}

func GenerateToken(userId string,userIp string) (string, error) {

	claims := JwtCustomClaims{
		jwt.StandardClaims{
			ExpiresAt: int64(time.Now().Add(time.Hour*24*365).Unix()),
			Issuer:    "finance.trytolog.com",
		},
		userId,

	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 这里需要注意当使用jwt.SigningMethodHS256方式生成token串时，SignedString方法的参数应该是[]byte数组，其他方式也对key有着要求
	tokenStr, err := token.SignedString([]byte(SecretKey+userIp))
	if err != nil {

		return "", err
	}

	return tokenStr, nil
}

func ParseToken(tokenStr string,userIp string) (claims jwt.MapClaims, err error) {
	// 创建对象，直接调用parse方法

	token, err := jwt.Parse(tokenStr, func(*jwt.Token) (interface{}, error) {
		// 注意这里也是[]byte数组

		return []byte(SecretKey+userIp),nil
	})
	if err != nil {
		return
	}

	// 获取jwt.Token对象后，获取定义的claims
	var ok bool
	// 注意 生成token串时无论是用自定义结构体的方式还是直接使用jwt.MapClaims，这里断言结果都为jwt.MapClaims（断言为自定义的结构体会失败）
	// 断言获取到的jwt.MapClaims实际上为map[string]interface{},使用key获取值时，需要再做一次类型断言，需要注意类型的转换，数值类型会被转化为float64（设置map["uid"] = 10,获取到的map["uid"]实际为float64类型的10.00000）
	claims, ok = token.Claims.(jwt.MapClaims)
	if !ok {
		return claims, fmt.Errorf("token claims can't convert to JwtCustomerClaims")
	}
	return claims,nil

}



