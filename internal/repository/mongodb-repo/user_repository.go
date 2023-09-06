package mongodbrepo

import (
	"context"
	"errors"
	"senpainikolay/go-internship-smartdata/auth-service/internal/models"
	mongomodels "senpainikolay/go-internship-smartdata/auth-service/internal/models/mongo-models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	// i need some logic on IDs and since mongo has just documents&&ObjIDs, I will return empty gorm model.
	"gorm.io/gorm"
)

type UserRepository struct {
	dbClient *mongo.Client
}

func NewUserRepository(dbClient *mongo.Client) *UserRepository {
	return &UserRepository{
		dbClient: dbClient,
	}
}

func (repo *UserRepository) Register(user *models.UserModel) error {
	collection := repo.dbClient.Database("test_db").Collection("users")

	userModel := mongomodels.User{
		Email:    user.Email,
		Password: user.Password,
		ImgPath:  user.ImgPath,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, insertErr := collection.InsertOne(ctx, userModel)
	if insertErr != nil {
		return insertErr
	}
	return nil
}

func (repo *UserRepository) CheckIfEmailExists(mail string) bool {

	var user mongomodels.User
	collection := repo.dbClient.Database("test_db").Collection("users")
	filter := bson.M{"email": mail}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, filter).Decode(&user)

	return !errors.Is(err, mongo.ErrNoDocuments)
}

func (repo *UserRepository) GetByEmail(email string) (models.UserModel, error) {
	var user mongomodels.User
	collection := repo.dbClient.Database("test_db").Collection("users")
	filter := bson.M{"email": email}
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return models.UserModel{}, err
	}
	return models.UserModel{
		Model: &gorm.Model{
			ID:        uint(1),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: gorm.DeletedAt{},
		},
		Email:    user.Email,
		Password: user.Password,
		ImgPath:  user.ImgPath,
	}, nil
}

func (repo *UserRepository) LogIn(userCredentials *models.UserCredentials) error {
	var user mongomodels.User
	collection := repo.dbClient.Database("test_db").Collection("users")
	filter := bson.M{"email": userCredentials.Email, "password": userCredentials.Password}
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	return err
}

func (repo *UserRepository) DeleteById(id uint) error {
	return errors.New("mongo db repo problem")
}

func (repo *UserRepository) GetById(id uint) (models.UserInfo, error) {

	return models.UserInfo{}, errors.New("mongo db repo problem")

}

func (repo *UserRepository) UpdateImage(id uint, imgPath string) error {
	return errors.New("mongo db repo problem")

}
