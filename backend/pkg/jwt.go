package pkg

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var (
	jwtSecret  = []byte("1234")
	effectTime = 3 * 24 * time.Hour
)

type Claims struct {
	id string `json:"id"`
	jwt.StandardClaims
}

// 解析token获得claims对象
func ParseToken(token string) (*Claims, error) {
	//用于解析鉴权的声明，方法内部主要是具体的解码和校验的过程，最终返回*Token
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// 生成token
func GenerateToken(claims *Claims) (string, error) {
	//设置有效期
	claims.ExpiresAt = time.Now().Add(effectTime).Unix()
	//生成token
	sign, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtSecret)
	//引入异常
	if err != nil {
		return "", err
	}
	return sign, nil
}
