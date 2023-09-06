package authservice

import (
	"context"
	"senpainikolay/go-internship-smartdata/gateway/pkg/auth-service/models"
	pb "senpainikolay/go-internship-smartdata/gateway/pkg/auth-service/pb"
	"time"
)

func NewUserServiceClient() *UserServiceClient {
	return &UserServiceClient{}
}

type UserServiceClient struct {
	pb.UserServiceClient
}

func (client *UserServiceClient) GetById(id int) (*models.UserInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
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

func (client *UserServiceClient) DeleteById(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	_, err := client.DeleteUserById(ctx, &pb.DeleteUserByIdRequest{Id: uint64(id)})

	return err
}

func (client *UserServiceClient) Register(usr models.UserModel) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	_, err := client.RegisterUser(ctx, &pb.CreateUserRequest{
		UserEntry: &pb.UserEntry{
			Email:    usr.Email,
			Password: usr.Password,
			ImgPath:  &usr.ImgPath,
		}})

	return err
}

func (client *UserServiceClient) LogIn(usrCrd models.UserCredentials) (models.UserTokensInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
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

func (client *UserServiceClient) GetUsrImage(usrId int) (*[]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	avatarInfo, err := client.GetImage(ctx, &pb.GetUserImageRequest{UsrId: uint64(usrId)})
	if err != nil {
		return nil, err
	}

	return &avatarInfo.Img, nil
}

func (client *UserServiceClient) UploadUsrImage(usrId int, img []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	_, err := client.UploadImage(ctx, &pb.UploadUserImageRequest{UsrId: uint64(usrId), Img: img})
	if err != nil {
		return err
	}

	return nil
}
