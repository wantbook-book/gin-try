package common

import (
	"github.com/dgrijalva/jwt-go"
	"hannibal/gin-try/model"
	"time"
)

var JwtKey = []byte("a_secret_key")

type Claims struct {
	UserId int
	jwt.StandardClaims
}

func ReleaseToken(user model.User) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := Claims{
		UserId: int(user.ID),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "hannibal",
			Subject:   "user token",
		},
	}
	//加密
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//转为token 字符串
	tokenString, err := token.SignedString(JwtKey)
	return tokenString, err
}

func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, err error) {
		return JwtKey, nil
	})
	return token, claims, err
}
