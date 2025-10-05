package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/zohaibsoomro/users_api_golang/models"
)

func TestGetAllUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	r.GET("/users", getAllUsers)

	req, _ := http.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var users []models.User
	err := json.Unmarshal(w.Body.Bytes(), &users)
	assert.NoError(t, err)

	assert.Greater(t, len(users), 0)
	assert.True(t, strings.Contains(users[0].Email, "@"))
}

func TestGetUserByEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	r.GET("/users/:email", getUserWithEmail)

	req, _ := http.NewRequest(http.MethodGet, "/users/khan@example.comm", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	var user models.User
	err := json.Unmarshal(resp.Body.Bytes(), &user)
	assert.NoError(t, err)

	assert.Equal(t, "khan", user.Name)
}

func TestUpdateUserByEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.PUT("/users/update/:email", updateUserWithEmail)
	user := `
		{
		"name": "khanam",
		"email": "khan@example.com"
	}`

	req, _ := http.NewRequest(http.MethodPut, "/users/update/khan@example.comm", strings.NewReader(user))
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)
	i, _ := userDb.GetUserByEmail("khan@example.comm")
	assert.Equal(t, -1, i)

	_, u := userDb.GetUserByEmail("khan@example.com")

	assert.Equal(t, "khanam", u.Name)

}

func TestAddUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.POST("/users/add", addUser)
	body := `{
		"name": "test",
		"email": "test@example.commm"
	}`

	req, _ := http.NewRequest(http.MethodPost, "/users/add", strings.NewReader(body))
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)

	_, u := userDb.GetUserByEmail("test@example.commm")
	assert.NotEqual(t, nil, u)
	assert.Equal(t, "test", u.Name)

}

func TestDeleteUserWithEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.DELETE("/users/delete/:email", deleteUserWithEmail)

	req, _ := http.NewRequest(http.MethodDelete, "/users/delete/zohaib@gmail.comm", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	i, _ := userDb.GetUserByEmail("zohaib@example.com")
	assert.Equal(t, -1, i)

}
