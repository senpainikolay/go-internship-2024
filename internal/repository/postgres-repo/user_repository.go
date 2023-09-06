package postgresrepo

import (
	"errors"
	"log"
	"senpainikolay/go-internship-smartdata/auth-service/internal/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	dbClient *gorm.DB
}

func NewUserRepository(dbClient *gorm.DB) *UserRepository {
	return &UserRepository{
		dbClient: dbClient,
	}
}

func (repo *UserRepository) Register(user *models.UserModel) error {
	err := repo.dbClient.Debug().
		Model(models.UserModel{}).
		Create(user).Error
	if err != nil {
		log.Printf("failed to insest user in database: %v\n", err)
		return err
	}
	return nil
}

func (repo *UserRepository) CheckIfEmailExists(mail string) bool {
	var user models.UserModel
	err := repo.dbClient.Debug().Model(models.UserModel{}).Find(&user).Where("email = ?", mail).Error
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func (repo *UserRepository) GetByEmail(email string) (models.UserModel, error) {
	var usr models.UserModel
	err := repo.dbClient.Model(&models.UserModel{}).Where("email = ?", email).First(&usr).Error
	if err != nil {
		return models.UserModel{}, err
	}
	return usr, nil
}

func (repo *UserRepository) LogIn(userCredentials *models.UserCredentials) error {
	// This is logic for an aditional call to the DataBase but can be replaced further by the JWT generation.
	var usr models.UserModel
	err := repo.dbClient.Model(&models.UserModel{}).Where("email = ? and password = ?", userCredentials.Email, userCredentials.Password).First(&usr).Error
	return err
}

func (repo *UserRepository) DeleteById(id uint) error {
	err := repo.dbClient.Model(&models.UserModel{}).Delete(&models.UserModel{}, id).Error
	return err
}

func (repo *UserRepository) GetById(id uint) (models.UserInfo, error) {
	var usr models.UserModel
	err := repo.dbClient.Model(&models.UserModel{}).Where("id = ?", id).First(&usr).Error
	if err != nil {
		return models.UserInfo{}, err
	}
	return models.UserInfo{ID: usr.ID, Email: usr.Email, ImgPath: usr.ImgPath}, nil
}

func (repo *UserRepository) UpdateImage(id uint, imgPath string) error {
	err := repo.dbClient.Model(&models.UserModel{}).Where("id = ?", id).Update("ImgPath", imgPath).Error
	return err
}
