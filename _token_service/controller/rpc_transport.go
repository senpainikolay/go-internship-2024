package controller

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	pb "senpainikolay/go-internship-smartdata/auth-service/_token_service/pb"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/grpc"
)

const (
	ACCESS_TOKEN_EXPIRE_TIME  = time.Minute * 15
	REFRESH_TOKEN_EXPIRE_TIME = time.Hour * 2
)

type TokenServer struct {
	pb.UnimplementedTokenServiceServer
}

func Serve(bind string) {
	listener, err := net.Listen("tcp", bind)
	if err != nil {
		log.Fatalf("gRPC server error: failure to bind %v\n", bind)
	}

	grpcServer := grpc.NewServer()

	tokenServer := TokenServer{}

	pb.RegisterTokenServiceServer(grpcServer, &tokenServer)
	log.Printf("gRPC API server listening on %v\n", bind)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("gRPC server error: %v\n", err)
	}
}

func (s *TokenServer) ValidateAccessToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {

	token, err := parseJwt(req.Token, os.Getenv("SECRET_ACCESS_TOKEN"))

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return &pb.ValidateTokenResponse{}, errors.New("token have expired")
		}

		return &pb.ValidateTokenResponse{UsrId: uint64(claims["id"].(float64))}, nil

	} else {
		return &pb.ValidateTokenResponse{}, err
	}
}

func (s *TokenServer) ValidateRefreshToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {

	token, err := parseJwt(req.Token, os.Getenv("SECRET_ACCESS_TOKEN"))

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return &pb.ValidateTokenResponse{}, errors.New("token have expired")
		}

		return &pb.ValidateTokenResponse{UsrId: uint64(claims["id"].(float64))}, nil

	} else {
		return &pb.ValidateTokenResponse{}, err
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

func (s *TokenServer) GenerateTokensPair(ctx context.Context, req *pb.GenerateTokenPairRequest) (*pb.GenerateTokenResponse, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  req.Id,
		"exp": time.Now().Add(ACCESS_TOKEN_EXPIRE_TIME).Unix(),
	})

	accessTk, err := token.SignedString([]byte(os.Getenv("SECRET_ACCESS_TOKEN")))

	if err != nil {
		return &pb.GenerateTokenResponse{}, errors.New("failed to create token")
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["sub"] = req.Id
	rtClaims["exp"] = time.Now().Add(REFRESH_TOKEN_EXPIRE_TIME).Unix()

	rt, err := refreshToken.SignedString([]byte(os.Getenv("SECRET_REFRESH_TOKEN")))
	if err != nil {
		return &pb.GenerateTokenResponse{}, err
	}

	return &pb.GenerateTokenResponse{
		AccessToken:  accessTk,
		RefreshToken: rt,
	}, nil
}
