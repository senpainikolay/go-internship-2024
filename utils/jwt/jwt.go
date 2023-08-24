package jwthelper

import (
	"errors"
	"fmt"
	"os"
	"senpainikolay/go-internship-smartdata/internal/models"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	_ "github.com/joho/godotenv/autoload"
)

const (
	TOKEN_EXPIRE_TIME = time.Hour * 8
)

func GenerateToken(id uint) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(TOKEN_EXPIRE_TIME).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_JWT")))

	if err != nil {
		return "", errors.New("failed to create token")
	}

	return tokenString, nil

}

func ValidateToken(tokenStr string) (models.UserJWTInfo, error) {

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET_JWT")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return models.UserJWTInfo{}, errors.New("token have expired")
		}

		return models.UserJWTInfo{ID: uint(claims["id"].(float64))}, nil

	} else {
		return models.UserJWTInfo{}, err
	}
}
