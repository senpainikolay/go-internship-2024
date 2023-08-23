package service

import (
	"errors"
	"mime/multipart"
	"senpainikolay/go-internship-smartdata/internal/models"
	utils_hash "senpainikolay/go-internship-smartdata/utils/hash"
	utils_img "senpainikolay/go-internship-smartdata/utils/image_handler"
	"strconv"
)

type IUserRepository interface {
	Register(user *models.UserModel) error
	GetById(uint) (models.UserInfo, error)
	CheckIfEmailExists(mail string) bool
	GetByEmail(string) (models.UserModel, error)
	DeleteById(uint) error
	LogIn(*models.UserCredentials) error
	UpdateImage(uint, string) error
}

type UserService struct {
	userRepo IUserRepository
}

func NewUserService(userRepo IUserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (svc *UserService) GetById(id uint) (models.UserInfo, error) {
	return svc.userRepo.GetById(id)
}

func (svc *UserService) Register(user *models.UserModel) error {
	if svc.userRepo.CheckIfEmailExists(user.Email) {
		return errors.New("user already exists")
	}

	hash, err := utils_hash.GeneratePasswordHash(user.Password)
	if err != nil {
		return errors.New("err can't register user")
	}

	user.Password = hash

	return svc.userRepo.Register(user)
}

func (svc *UserService) LogIn(userCredentials *models.UserCredentials) error {

	usr, err := svc.userRepo.GetByEmail(userCredentials.Email)
	if err != nil {
		return errors.New("something wrong with credentials")
	}
	err = utils_hash.ComparePasswordHash(usr.Password, userCredentials.Password)
	if err != nil {
		return errors.New("something wrong with credentials")
	}
	// This make an aditional call to the DataBase but can be replaced further by the JWT generation.
	userCredentials.Password = usr.Password
	return svc.userRepo.LogIn(userCredentials)
}

func (svc *UserService) DeleteById(id uint) error {
	return svc.userRepo.DeleteById(id)
}

func (svc *UserService) UpdateImage(id uint, file *multipart.File, fileName string) error {

	uniqueImgPath := strconv.Itoa(int(id)) + fileName

	err := utils_img.CreateImageFile(file, uniqueImgPath)
	if err != nil {
		return err
	}

	err = svc.userRepo.UpdateImage(id, uniqueImgPath)
	if err != nil {
		return err
	}
	return nil

}

func (svc *UserService) GetImage(user_id uint) (*[]byte, error) {

	usr, err := svc.userRepo.GetById(user_id)
	if err != nil {
		return nil, err
	}

	binaryImgAddr, err := utils_img.ReadImageFile(usr.ImgPath)
	if err != nil {
		return nil, err
	}

	return binaryImgAddr, nil

}
