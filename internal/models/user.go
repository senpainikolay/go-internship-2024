package models

import "gorm.io/gorm"

type UserModel struct {
	*gorm.Model
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	ImgPath  string `json:"avatar"`
}

type UserInfo struct {
	ID      uint   `json:"id"`
	Email   string `json:"email"`
	ImgPath string `json:"avatar"`
}

type UserJWTInfo struct {
	ID uint `json:"id"`
}

type UserCredentials struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
