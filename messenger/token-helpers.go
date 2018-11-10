package main

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// TOKEN helpers

func (u *User) generateToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
		"user_id": u.ID,
	})
	tokenString, _ := token.SignedString(mySigningKey)
	return tokenString
}

func getUserFromToken(t string) *User {
	token, _ := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})
	claims := token.Claims.(jwt.MapClaims)
	uid := claims["user_id"].(string)
	return findUserByID(uid)
}
