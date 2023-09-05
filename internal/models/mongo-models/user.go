package mongomodels

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Email    string             `bson:"email"`
	Password string             `bson:"password"`
	ImgPath  string             `bson:"avatar"`
}

type UserInfo struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Email   string             `bson:"email"`
	ImgPath string             `bson:"avatar"`
}

type UserJWTInfo struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
}

type UserCredentials struct {
	Email    string `bson:"email"`
	Password string `bson:"password"`
}
