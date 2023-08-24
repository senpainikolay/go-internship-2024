package service

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
	"senpainikolay/go-internship-smartdata/internal/models"
	utils_hash "senpainikolay/go-internship-smartdata/utils/hash"
	jwthelper "senpainikolay/go-internship-smartdata/utils/jwt"
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

type ICacheRepository interface {
	Get(string) (string, error)
	Set(string, string) error
}

type UserService struct {
	userRepo  IUserRepository
	cacheRepo ICacheRepository
}

func NewUserService(userRepo IUserRepository, cacheRepo ICacheRepository) *UserService {
	return &UserService{
		userRepo:  userRepo,
		cacheRepo: cacheRepo,
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

func (svc *UserService) LogIn(userCredentials *models.UserCredentials) (string, error) {

	usr, err := svc.userRepo.GetByEmail(userCredentials.Email)
	if err != nil {
		return "", errors.New("something wrong with credentials")
	}
	err = utils_hash.ComparePasswordHash(usr.Password, userCredentials.Password)
	if err != nil {
		return "", errors.New("something wrong with credentials")
	}

	token, err := jwthelper.GenerateToken(usr.ID)
	if err != nil {
		return "", err
	}

	idStr := strconv.FormatUint(uint64(usr.ID), 10)
	err = svc.cacheRepo.Set(idStr, token)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (svc *UserService) DeleteById(id uint) error {
	return svc.userRepo.DeleteById(id)
}

func (svc *UserService) UpdateImage(id uint, file *multipart.File, fileName string) error {

	uniqueImgPath := strconv.Itoa(int(id)) + fileName

	f, err := os.OpenFile("./uploads/"+uniqueImgPath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return errors.New("could not create image")
	}
	defer f.Close()

	_, err = io.Copy(f, *file)
	if err != nil {
		return errors.New("could not copy image from request to the server")
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

	file, err := os.OpenFile("./uploads/"+usr.ImgPath, os.O_RDONLY, 0)
	if err != nil {
		return nil, errors.New("could not open the image")
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	binaryData := make([]byte, fileInfo.Size())
	_, err = file.Read(binaryData)
	if err != nil {
		return nil, err
	}

	return &binaryData, nil

}
