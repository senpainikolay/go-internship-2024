package repository

import (
	"database/sql"
	"log"
	db_handler "senpainikolay/go-internship-smartdata/db"
)

type UserRepository struct {
	DBClient string
}

func NewUserRepository(dbClient string) *UserRepository {
	db, err := sql.Open("postgres", dbClient)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db_handler.TryCreate(db)

	return &UserRepository{
		DBClient: dbClient,
	}
}

func (repo *UserRepository) Register(email, pw string) error {
	db, err := sql.Open("postgres", repo.DBClient)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db_handler.CreateUser(db, email, pw)

	return nil
}

func (repo *UserRepository) LogIn(email, pw string) string {
	db, err := sql.Open("postgres", repo.DBClient)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	msg, _ := db_handler.LogInUser(db, email, pw)

	return msg
}

func (repo *UserRepository) DeleteById(id int) error {
	db, err := sql.Open("postgres", repo.DBClient)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db_handler.DeleteUser(db, id)

	return nil
}

func (repo *UserRepository) GetById(id int) string {
	db, err := sql.Open("postgres", repo.DBClient)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	user, err := db_handler.GetUserById(db, id)
	if err != nil {
		log.Println(err.Error())
	}
	return user

}
