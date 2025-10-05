package api

import (
	"net/http"
	"net/mail"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zohaibsoomro/users_api_golang/models"
	"github.com/zohaibsoomro/users_api_golang/utils"
)

type Api struct {
	Address string
	Db      *utils.UsersDB
}

const DefaultAddress = "localhost:8080"

func NewApi() *Api {
	return &Api{
		Address: DefaultAddress,
		Db:      utils.NewUsersDB(),
	}
}
func NewApiWithAddress(address string) *Api {
	return &Api{
		Address: address,
		Db:      utils.NewUsersDB(),
	}
}

func (api *Api) RegisterHandlers() *gin.Engine {
	server := gin.Default()

	server.GET("/", api.HelloWorld)
	server.GET("/users", api.GetAllUsers)
	server.POST("/users/add", api.AddUser)
	server.GET("/users/:email", api.GetUserWithEmail)
	server.PUT("/users/update/:email", api.UpdateUserWithEmail)
	server.DELETE("/users/delete/:email", api.DeleteUserWithEmail)
	server.GET("/users/shutdown", func(c *gin.Context) {
		c.JSON(200, "Shutting down server...")
		os.Exit(0)
	})
	return server
}

func (api *Api) HelloWorld(c *gin.Context) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(`<h1>Welcome Back!</h1>`))
}

func (api *Api) GetAllUsers(c *gin.Context) {
	var users *[]models.User
	api.Db.DB.Find(&users)

	c.JSON(http.StatusOK, &users)

}
func (api *Api) GetUserWithEmail(c *gin.Context) {

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
	user := api.Db.GetUserByEmail(email)
	if user == nil {
		writeError(c, http.StatusNotFound, "User not found!")

		return
	}

	c.JSON(http.StatusOK, user)
}

func (api *Api) UpdateUserWithEmail(c *gin.Context) {

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

	if err := api.Db.UpdateUserByEmail(email, &user); err != nil {
		writeError(c, http.StatusInternalServerError, "Unable to update user: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully!"})

}

func (api *Api) DeleteUserWithEmail(c *gin.Context) {

	email := strings.TrimSpace(c.Param("email"))

	if _, err := mail.ParseAddress(email); err != nil {
		writeError(c, http.StatusBadRequest, "Invalid email entered: "+err.Error())
		return
	}

	err := api.Db.DeleteUserByEmail(email)
	if err != nil {
		writeError(c, http.StatusNotAcceptable, "Error while deleting user: "+err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully!"})
}

func (api *Api) AddUser(c *gin.Context) {

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

	err := api.Db.AddUser(u)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "Unable to add user: "+err.Error())
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "user added successfully."})
}

// WriteError writes an error response in JSON.
func writeError(c *gin.Context, status int, msg string) {
	c.JSON(status, gin.H{"error": msg})
}
