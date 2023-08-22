package models

type UserModel struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserEntry struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

type UserCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
