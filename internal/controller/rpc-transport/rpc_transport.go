package rpctransport

import (
	"context"
	"log"
	"net"
	"senpainikolay/go-internship-smartdata/internal/models"
	pb "senpainikolay/go-internship-smartdata/internal/pb"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type IUserService interface {
	GetById(uint) (models.UserInfo, error)
	Register(*models.UserModel) error
	LogIn(*models.UserCredentials) (map[string]string, error)
	DeleteById(uint) error
	GetImage(uint) (*[]byte, error)
	UpdateImage(uint, *[]byte, string) error
}

type UserServer struct {
	pb.UnimplementedUserServiceServer
	userSvc IUserService
}

func Serve(usr_service IUserService, bind string) {
	listener, err := net.Listen("tcp", bind)
	if err != nil {
		log.Fatalf("gRPC server error: failure to bind %v\n", bind)
	}

	grpcServer := grpc.NewServer()

	userServer := UserServer{userSvc: usr_service}

	pb.RegisterUserServiceServer(grpcServer, &userServer)
	log.Printf("gRPC API server listening on %v\n", bind)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("gRPC server error: %v\n", err)
	}
}

func (s *UserServer) GetUserById(ctx context.Context, req *pb.GetUserByIdRequest) (*pb.UserResponse, error) {
	usrInfo, err := s.userSvc.GetById(uint(req.Id))
	if err != nil {
		return &pb.UserResponse{}, err
	}
	return &pb.UserResponse{
		Id:      uint64(usrInfo.ID),
		Email:   usrInfo.Email,
		ImgPath: usrInfo.ImgPath,
	}, nil
}

func (s *UserServer) DeleteUserById(ctx context.Context, req *pb.DeleteUserByIdRequest) (*emptypb.Empty, error) {
	err := s.userSvc.DeleteById(uint(req.Id))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (s *UserServer) RegisterUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	usr := &models.UserModel{
		Email:    req.UserEntry.Email,
		Password: req.UserEntry.Password,
	}
	err := s.userSvc.Register(usr)
	if err != nil {
		return &pb.UserResponse{}, err
	}

	return &pb.UserResponse{
		Email: usr.Email,
	}, nil
}

func (s *UserServer) LogInUser(ctx context.Context, req *pb.LogInUserRequest) (*pb.LogInResponse, error) {
	usrCredentials := &models.UserCredentials{
		Email:    req.UserEntry.Email,
		Password: req.UserEntry.Password,
	}
	tokensMap, err := s.userSvc.LogIn(usrCredentials)
	if err != nil {
		return &pb.LogInResponse{}, err
	}

	return &pb.LogInResponse{
		AccessToken:  tokensMap["access_token"],
		RefreshToken: tokensMap["refresh_token"],
	}, nil
}

func (s *UserServer) GetImage(ctx context.Context, req *pb.GetUserImageRequest) (*pb.AvatarResponse, error) {
	binaryImgAddr, err := s.userSvc.GetImage(uint(req.UsrId))
	if err != nil {
		return &pb.AvatarResponse{}, err
	}

	return &pb.AvatarResponse{
		UsrId: req.UsrId,
		Img:   *binaryImgAddr,
	}, nil
}

func (s *UserServer) UploadImage(ctx context.Context, req *pb.UploadUserImageRequest) (*pb.AvatarResponse, error) {

	usr_id_str := strconv.Itoa(int(req.UsrId))
	err := s.userSvc.UpdateImage(uint(req.UsrId), &req.Img, usr_id_str+".jpg")
	if err != nil {
		return &pb.AvatarResponse{}, err
	}
	return &pb.AvatarResponse{
		UsrId: req.UsrId,
		Img:   req.Img,
	}, nil

}
