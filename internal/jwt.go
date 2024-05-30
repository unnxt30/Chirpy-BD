package database

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (user *User) GenerateToken() string {
	// var expiresAt int
	var expiresAt int = user.ExpirationTime
	claims := &jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiresAt) * time.Second)),
		Subject:   strconv.Itoa(user.ID),
	}
	myToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := os.Getenv("JWT_SECRET")
	// fmt.Println(secretKey)
	returnToken, err := myToken.SignedString([]byte(secretKey))
	if err != nil {
		fmt.Println(err)
		return "error"
	}
	return returnToken
}

func CheckToken(tokenString string) int {
	type parseClaimStruct struct {
		tokenString string
		jwt.RegisteredClaims
	}

	var uid int
	token, err := jwt.ParseWithClaims(tokenString, &parseClaimStruct{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		fmt.Println(err.Error());
		return -1
	}

	claims := token.Claims.(*parseClaimStruct)
	uid, _ = strconv.Atoi(claims.Subject)
	return uid;
}
