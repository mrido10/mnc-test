package util

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"test2/model"
	"time"
)

type JwtToken interface {
	GenerateToken(userID string, expTime time.Duration) (string, *model.Error)
	ParseToken(jwtToken string) (jwt.MapClaims, *model.Error)
}

type jwtToken struct {
	secretKey []byte
}

func NewJwtToken(secret string) JwtToken {
	return jwtToken{
		secretKey: []byte(secret),
	}
}

func (j jwtToken) GenerateToken(userID string, expTime time.Duration) (string, *model.Error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(expTime).Unix(),
	})

	tokenStr, err := token.SignedString(j.secretKey)
	if err != nil {
		return "", model.NewError(500, "Internal Server Error", err)
	}
	return tokenStr, nil
}

func (j jwtToken) ParseToken(jwtToken string) (jwt.MapClaims, *model.Error) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secretKey, nil
	})
	if err != nil {
		return nil, model.NewError(500, "Internal Server Error", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, model.NewError(401, "Unauthorized", err)
	}
	return claims, nil
}
