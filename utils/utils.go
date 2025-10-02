package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/zohaibsoomro/users_api_golang/models"
)

type UsersDB struct {
	users    []models.User
	lock     *sync.Mutex
	filename string
}

func NewUsersDB(file string) *UsersDB {
	db := &UsersDB{
		users:    make([]models.User, 0),
		lock:     &sync.Mutex{},
		filename: file,
	}
	db.LoadAllUsers()
	return db
}

func (db *UsersDB) LoadAllUsers() error {
	bytes, err := os.ReadFile(db.filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, &db.users)
	if err != nil {

		return fmt.Errorf("unable to parse json: %w", err)
	}
	return nil
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
			return idx, &db.users[idx]
		}
	}
	return -1, nil
}

func (db *UsersDB) AddUser(user models.User) error {
	db.lock.Lock()
	defer db.lock.Unlock()

	_, u := db.GetUserByEmail(user.Email)
	if u != nil {
		return errors.New("user already exists with that email")
	}
	id := 1
	if len(db.users) > 0 {
		id = db.users[len(db.users)-1].Id + 1
	}
	user.Id = id
	db.users = append(db.users, user)
	bytes, err := json.MarshalIndent(&db.users, "", "\t")
	if err != nil {
		return err
	}
	if err := os.WriteFile(db.filename, bytes, 0644); err != nil {
		return err
	}

	fmt.Println("User added.")
	db.LoadAllUsers()
	return nil
}

func (db *UsersDB) DeleteUserByEmail(email string) error {
	db.lock.Lock()
	defer db.lock.Unlock()
	email = strings.TrimSpace(email)
	i, u := db.GetUserByEmail(email)
	if u != nil {
		db.users = append(db.users[:i], db.users[i+1:]...)
		bytes, err := json.MarshalIndent(&db.users, "", "\t")
		if err != nil {
			return fmt.Errorf("unable to parse string to json: %w", err)
		}
		err = os.WriteFile(db.filename, bytes, 0644)
		if err != nil {
			return fmt.Errorf("unable to write json to file: %w", err)
		}
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
	db.lock.Lock()
	defer db.lock.Unlock()
	email = strings.TrimSpace(email)
	idx, user := db.GetUserByEmail(email)

	if user == nil {
		return errors.New("user not found")
	}

	db.users[idx].Email = u.Email
	db.users[idx].Name = u.Name

	bytes, err := json.MarshalIndent(&db.users, "", "\t")
	if err != nil {
		return fmt.Errorf("unable to parse string to Json: %w", err)
	}
	err = os.WriteFile(db.filename, bytes, 0644)
	if err != nil {
		return fmt.Errorf("unable to write json to file: %w", err)
	}
	fmt.Println("User Updated")

	db.LoadAllUsers()
	return nil
}
