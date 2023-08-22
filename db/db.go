package db

import (
	"database/sql"
	"log"
)

func TryCreate(db *sql.DB) {
	_, err := db.Exec(`
	 	CREATE TABLE users (
	 		id            SERIAL   PRIMARY KEY,
	 		email         TEXT UNIQUE,
	 		password  TEXT
	 	);
	 `)

	if err != nil {
		log.Println(err.Error())

	}
}

func CreateUser(db *sql.DB, email string, pw string) error {
	_, err := db.Exec(`INSERT INTO
		users(email, password)
		VALUES($1, $2)`, email, pw)

	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func DeleteUser(db *sql.DB, id int) error {
	_, err := db.Exec(` 
	DELETE FROM users WHERE id = $1 
	`, id)

	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func LogInUser(db *sql.DB, email string, pw string) (string, error) {
	row := db.QueryRow("SELECT email FROM users WHERE email=$1 AND password = $2", email, pw)
	var usrEmail string
	err := row.Scan(&usrEmail)
	if err != nil {
		if err == sql.ErrNoRows {
			return "something wrong with credentials!", err
		}
		return "", err
	}
	return usrEmail + "succesfully logged in", nil

}

func GetUserById(db *sql.DB, id int) (string, error) {
	row := db.QueryRow("SELECT email FROM users WHERE id=$1", id)
	var email string

	err := row.Scan(&email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "no user found", err
		}
		return "", err
	}

	return email, nil

}
