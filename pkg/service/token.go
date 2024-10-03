package service

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId string `json:"user_id"`
	Role   string `json:"role"`
}

func CreateToken(userId, role string) (token string) {
	params := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(12 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserId: userId,
		Role:   role,
	})

	token, _ = params.SignedString([]byte(signingKey))
	return token
}
