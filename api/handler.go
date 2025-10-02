package api

import (
	"net/http"
	"net/mail"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zohaibsoomro/users_api_golang/models"
	"github.com/zohaibsoomro/users_api_golang/utils"
)

var userDb = utils.NewUsersDB("/Users/zohaib/Documents/Golang/GolangPractice/users.json")

type Api struct {
	Address string
}

func NewApi(address string) *Api {
	return &Api{
		Address: address,
	}
}

func (api *Api) RegisterHandlers() *gin.Engine {
	server := gin.Default()

	server.GET("/", helloWorld)
	server.GET("/users", getAllUsers)
	server.POST("/users/add", addUser)
	server.GET("/users/:email", getUserWithEmail)
	server.PUT("/users/update/:email", updateUserWithEmail)
	server.DELETE("/users/delete/:email", deleteUserWithEmail)
	return server
}

func helloWorld(c *gin.Context) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(`<h1>Welcome Back!</h1>`))
}

func getAllUsers(c *gin.Context) {
	users := userDb.GetAllUsers()

	c.JSON(http.StatusOK, users)

}
func getUserWithEmail(c *gin.Context) {

	email := strings.ToLower(strings.TrimSpace(c.Param("email")))

	if email == "" {
		writeError(c, http.StatusBadRequest, "Email parameter is required.")
		return
	}

	if _, err := mail.ParseAddress(email); err != nil {
		writeError(c, http.StatusBadRequest, "Invalid email entered: "+strings.TrimPrefix(err.Error(), "mail: "))
		return
	}

	// Lookup user
	_, user := userDb.GetUserByEmail(email)
	if user == nil {
		writeError(c, http.StatusNotFound, "User not found!")

		return
	}

	c.JSON(http.StatusOK, user)
}

func updateUserWithEmail(c *gin.Context) {

	email := strings.TrimSpace(c.Param("email"))

	if email == "" {
		writeError(c, http.StatusBadRequest, "email param is missing")
		return
	}

	if _, err := mail.ParseAddress(email); err != nil {
		writeError(c, http.StatusBadRequest, "Invalid email entered: "+strings.TrimPrefix(err.Error(), "mail: "))
		return
	}

	var user models.User
	if err := c.ShouldBindBodyWithJSON(&user); err != nil {
		writeError(c, http.StatusBadRequest, "Unable to Decode Request Body")
		return
	}

	if _, err := mail.ParseAddress(user.Email); err != nil {
		writeError(c, http.StatusBadRequest, "Invalid email in body: "+strings.TrimPrefix(err.Error(), "mail: "))
		return
	}

	if err := userDb.UpdateUserByEmail(email, &user); err != nil {
		writeError(c, http.StatusInternalServerError, "Unable to update user: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully!"})

}

func deleteUserWithEmail(c *gin.Context) {

	email := strings.TrimSpace(c.Param("email"))

	if _, err := mail.ParseAddress(email); err != nil {
		writeError(c, http.StatusBadRequest, "Invalid email entered: "+err.Error())
		return
	}

	err := userDb.DeleteUserByEmail(email)
	if err != nil {
		writeError(c, http.StatusNotAcceptable, "Error while deleting user: "+err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully!"})
}

func addUser(c *gin.Context) {

	var u models.User
	if err := c.ShouldBindBodyWithJSON(&u); err != nil {
		writeError(c, http.StatusBadRequest, "Unable to parse body: "+err.Error())
		return
	}
	u.Email = strings.TrimSpace(u.Email)
	if _, err := mail.ParseAddress(u.Email); err != nil {
		writeError(c, http.StatusBadRequest, "Invalid email entered: "+err.Error())
		return
	}

	err := userDb.AddUser(u)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "Unable to add user: "+err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user added successfully."})
}

// WriteError writes an error response in JSON.
func writeError(c *gin.Context, status int, msg string) {
	c.JSON(status, gin.H{"error": msg})
}
