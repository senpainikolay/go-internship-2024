package tokenvalservice

import (
	"context"
	"log"
	"senpainikolay/go-internship-smartdata/gateway/pkg/models"
	pb "senpainikolay/go-internship-smartdata/gateway/pkg/token-val-service/pb"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TokenValidationServiceClient struct {
	pb.TokenServiceClient
}

func NewUserServiceClient(port string) *TokenValidationServiceClient {
	conn, err := grpc.Dial(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return &TokenValidationServiceClient{
		pb.NewTokenServiceClient(conn),
	}
}

func (client *TokenValidationServiceClient) ValidateUsrAccessToken(token string) (*models.UserJWTInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*8)
	defer cancel()

	usrRes, err := client.ValidateAccessToken(ctx, &pb.ValidateTokenRequest{Token: token})
	if err != nil {
		return &models.UserJWTInfo{}, err
	}

	return &models.UserJWTInfo{
		ID: uint(usrRes.UsrId),
	}, nil
}

func (client *TokenValidationServiceClient) ValidateUsrRefreshToken(token string) (*models.UserJWTInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*8)
	defer cancel()

	usrRes, err := client.ValidateRefreshToken(ctx, &pb.ValidateTokenRequest{Token: token})
	if err != nil {
		return &models.UserJWTInfo{}, err
	}

	return &models.UserJWTInfo{
		ID: uint(usrRes.UsrId),
	}, nil
}

func (client *TokenValidationServiceClient) GenerateUsrTokenPair(id uint64) (*models.UserTokensInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*8)
	defer cancel()

	tokens, err := client.GenerateTokenPair(ctx, &pb.GenerateTokenPairRequest{Id: id})
	if err != nil {
		return &models.UserTokensInfo{}, err
	}
	return &models.UserTokensInfo{
		AccesToken:   tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}
