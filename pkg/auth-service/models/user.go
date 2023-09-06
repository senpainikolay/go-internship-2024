package models

type UserModel struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	ImgPath  string `json:"avatar"`
}

type UserCredentials struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserTokensInfo struct {
	AccesToken   string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserInfo struct {
	ID      uint   `json:"id"`
	Email   string `json:"email"`
	ImgPath string `json:"avatar"`
}
