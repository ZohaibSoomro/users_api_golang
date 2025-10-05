package utils

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/zohaibsoomro/users_api_golang/models"
	"gorm.io/gorm"
)

type UsersDB struct {
	DB   *gorm.DB
	lock *sync.Mutex
}

func NewUsersDB() *UsersDB {
	db := &UsersDB{
		DB:   NewDb(),
		lock: &sync.Mutex{},
	}

	return db
}

func (db *UsersDB) ListAllUsers() {
	for _, u := range db.GetAllUsers() {
		fmt.Printf("Id: %v, Name: %v, Email: %v\n", u.Id, u.Name, u.Email)
	}
}

func (db *UsersDB) GetAllUsers() []models.User {
	var users []models.User
	dbb := db.DB.Find(&users)
	if dbb.Error != nil {
		log.Fatal("Error getting all users:", dbb.Error)
	}
	return users
}

func (db *UsersDB) GetUserByEmail(email string) *models.User {
	email = strings.TrimSpace(email)
	var user models.User
	result := db.DB.First(&user, "email=?", email)
	if result.Error != nil {
		// if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// 	fmt.Println("User not found")
		// } else {
		// 	fmt.Println("Error:", result.Error)
		// }
		return nil
	}
	return &user
}

func (db *UsersDB) AddUser(user models.User) error {
	db.lock.Lock()
	defer db.lock.Unlock()

	u := db.GetUserByEmail(user.Email)
	if u != nil {
		return errors.New("user already exists with that email")
	}

	db.DB.Create(&user)
	if user.Id < 0 || db.DB.Error != nil {
		return fmt.Errorf("error while adding user: %v", db.DB.Error.Error())
	}

	// fmt.Println("User added.")
	return nil
}

func (db *UsersDB) DeleteUserByEmail(email string) error {
	db.lock.Lock()
	defer db.lock.Unlock()
	email = strings.TrimSpace(email)
	u := db.GetUserByEmail(email)
	if u == nil {
		return fmt.Errorf("unable to find user")
	}

	if err := db.DB.Delete(u, u.Id).Error; err != nil {
		return fmt.Errorf("unable to update user: %v", err.Error())
	}

	return nil

}

func (db *UsersDB) UpdateUserByEmail(email string, u *models.User) error {
	db.lock.Lock()
	defer db.lock.Unlock()
	email = strings.TrimSpace(email)
	user := db.GetUserByEmail(email)

	if user == nil {
		return errors.New("user not found")
	}
	u.Id = user.Id
	if err := db.DB.Save(u).Error; err != nil {
		return fmt.Errorf("unable to update user: %v", err.Error())
	}

	return nil
}
