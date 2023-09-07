package authservice

import (
	"context"
	"log"
	pb "senpainikolay/go-internship-smartdata/gateway/pkg/auth-service/pb"
	"senpainikolay/go-internship-smartdata/gateway/pkg/models"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewUserServiceClient(port string) *UserServiceClient {
	conn, err := grpc.Dial(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return &UserServiceClient{
		pb.NewUserServiceClient(conn),
	}
}

type UserServiceClient struct {
	pb.UserServiceClient
}

func (client *UserServiceClient) GetById(id uint) (*models.UserInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*8)
	defer cancel()

	usrRes, err := client.GetUserById(ctx, &pb.GetUserByIdRequest{Id: uint64(id)})
	if err != nil {
		return &models.UserInfo{}, err
	}

	return &models.UserInfo{
		ID:      uint(usrRes.Id),
		Email:   usrRes.Email,
		ImgPath: usrRes.ImgPath,
	}, nil
}

func (client *UserServiceClient) DeleteById(id uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*8)
	defer cancel()

	_, err := client.DeleteUserById(ctx, &pb.DeleteUserByIdRequest{Id: uint64(id)})

	return err
}

func (client *UserServiceClient) Register(usr *models.UserModel) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*8)
	defer cancel()

	_, err := client.RegisterUser(ctx, &pb.CreateUserRequest{
		UserEntry: &pb.UserEntry{
			Email:    usr.Email,
			Password: usr.Password,
			ImgPath:  &usr.ImgPath,
		}})

	return err
}

func (client *UserServiceClient) LogIn(usrCrd *models.UserCredentials) (models.UserTokensInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*8)
	defer cancel()

	usrTokenMap, err := client.LogInUser(ctx, &pb.LogInUserRequest{
		UserEntry: &pb.UserEntry{
			Email:    usrCrd.Email,
			Password: usrCrd.Password,
		}})

	if err != nil {
		return models.UserTokensInfo{}, err
	}

	return models.UserTokensInfo{
		AccesToken:   usrTokenMap.AccessToken,
		RefreshToken: usrTokenMap.RefreshToken,
	}, nil
}

func (client *UserServiceClient) GetUsrImage(usrId uint64) (*[]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*8)
	defer cancel()

	avatarInfo, err := client.GetImage(ctx, &pb.GetUserImageRequest{UsrId: usrId})
	if err != nil {
		return nil, err
	}

	return &avatarInfo.Img, nil
}

func (client *UserServiceClient) UploadUsrImage(usrId uint64, img *[]byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*8)
	defer cancel()

	_, err := client.UploadImage(ctx, &pb.UploadUserImageRequest{UsrId: usrId, Img: *img})
	if err != nil {
		return err
	}

	return nil
}
