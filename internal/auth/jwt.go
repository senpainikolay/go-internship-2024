package auth

import (
	"errors"
	"fmt"
	"os"
	"senpainikolay/go-internship-smartdata/auth-service/internal/models"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	_ "github.com/joho/godotenv/autoload"
)

const (
	ACCESS_TOKEN_EXPIRE_TIME  = time.Minute * 15
	REFRESH_TOKEN_EXPIRE_TIME = time.Hour * 2
)

func GenerateTokenPair(id uint) (map[string]string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(ACCESS_TOKEN_EXPIRE_TIME).Unix(),
	})

	accessTk, err := token.SignedString([]byte(os.Getenv("SECRET_ACCESS_TOKEN")))

	if err != nil {
		return nil, errors.New("failed to create token")
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["sub"] = id
	rtClaims["exp"] = time.Now().Add(REFRESH_TOKEN_EXPIRE_TIME).Unix()

	rt, err := refreshToken.SignedString([]byte(os.Getenv("SECRET_REFRESH_TOKEN")))
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"access_token":  accessTk,
		"refresh_token": rt,
	}, nil

}

func ValidateAccessToken(tokenStr string) (models.UserJWTInfo, error) {

	token, err := parseJwt(tokenStr, os.Getenv("SECRET_ACCESS_TOKEN"))

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return models.UserJWTInfo{}, errors.New("token have expired")
		}

		return models.UserJWTInfo{ID: uint(claims["id"].(float64))}, nil

	} else {
		return models.UserJWTInfo{}, err
	}
}

func ValidateRefreshToken(tokenStr string) (uint, error) {

	token, err := parseJwt(tokenStr, os.Getenv("SECRET_REFRESH_TOKEN"))
	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return 0, errors.New("token have expired, log in again")
		}

		return uint(claims["sub"].(float64)), nil

	} else {
		return 0, errors.New("something wrong with the token")
	}
}

func parseJwt(tokenStr string, secret_token string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret_token), nil
	})

	return token, err
}
