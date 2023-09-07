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
