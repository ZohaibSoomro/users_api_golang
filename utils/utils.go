package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/zohaibsoomro/users_api_golang/models"
)

type UsersDB struct {
	users []models.User
}

func NewUsersDB() *UsersDB {
	return &UsersDB{
		users: make([]models.User, 0),
	}
}

const fileName = "users.json"

func (db *UsersDB) LoadAllUsers() {
	bytes, err := os.ReadFile(fileName)
	handleError("Unable to read file", err)
	err = json.Unmarshal(bytes, &db.users)
	handleError("Unable to parse json", err)
}

func handleError(msg string, err error) {
	if err != nil {
		log.Fatalf("%v: %v", msg, err)
	}
}

func (db *UsersDB) ListAllUsers() {
	for _, u := range db.users {
		fmt.Printf("Id: %v, Name: %v, Email: %v\n", u.Id, u.Name, u.Email)
	}
}

func (db *UsersDB) GetAllUsers() []models.User {
	return db.users
}

func (db *UsersDB) GetUserByEmail(email string) (int, *models.User) {
	email = strings.TrimSpace(email)
	for idx, user := range db.users {
		if strings.EqualFold(email, user.Email) {
			return idx, &user
		}
	}
	return -1, nil
}

func (db *UsersDB) DeleteUserByEmail(email string) error {
	email = strings.TrimSpace(email)
	i, u := db.GetUserByEmail(email)
	if u != nil {
		db.users = append(db.users[:i], db.users[i+1:]...)
		bytes, err := json.MarshalIndent(&db.users, "", "\t")
		handleError("Unable to parse string to Json", err)
		err = os.WriteFile(fileName, bytes, 0644)
		handleError("Unable to write json to file", err)
		fmt.Println("Deleted User\nName:", u.Name, "\nEmail:", u.Email)
		db.LoadAllUsers()
		err = nil

	} else {
		fmt.Println("No user found for", email)
		return errors.New("user not found")
	}
	return nil

}

func (db *UsersDB) UpdateUserByEmail(email string, u *models.User) error {
	email = strings.TrimSpace(email)
	idx, user := db.GetUserByEmail(email)

	if user == nil {
		return errors.New("user not found")
	}

	db.users[idx].Email = u.Email
	db.users[idx].Name = u.Name

	bytes, err := json.MarshalIndent(&db.users, "", "\t")
	handleError("Unable to parse string to Json", err)
	err = os.WriteFile(fileName, bytes, 0644)
	handleError("Unable to write json to file", err)
	fmt.Println("User Updated")

	db.LoadAllUsers()
	return nil
}
