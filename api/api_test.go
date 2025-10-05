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

func TestAddUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	api := NewApi()
	tx := api.Db.DB.Begin()
	api.Db.DB = tx
	defer tx.Rollback()

	r.POST("/users/add", api.AddUser)
	body := `{
		"name": "test1",
		"email": "test1@example.comm"
	}`

	req, _ := http.NewRequest(http.MethodPost, "/users/add", strings.NewReader(body))
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)

	u := api.Db.GetUserByEmail("test1@example.comm")
	assert.NotEqual(t, nil, u)
	assert.Equal(t, "test1", u.Name)

}

func TestGetAllUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	r.GET("/users", NewApi().GetAllUsers)

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
	r.GET("/users/:email", NewApi().GetUserWithEmail)

	req, _ := http.NewRequest(http.MethodGet, "/users/zohaib@gmail.comm", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	var user models.User
	err := json.Unmarshal(resp.Body.Bytes(), &user)
	assert.NoError(t, err)

	assert.Equal(t, "zohaib", user.Name)
}

func TestUpdateUserByEmail(t *testing.T) {

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	api := NewApi()
	tx := api.Db.DB.Begin()
	api.Db.DB = tx
	defer tx.Rollback()

	r.PUT("/users/update/:email", api.UpdateUserWithEmail)
	user := `
		{
		"name": "aibee006",
		"email": "aibee@example.com"
	}`

	req, _ := http.NewRequest(http.MethodPut, "/users/update/aibee@example.com", strings.NewReader(user))
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	u := api.Db.GetUserByEmail("aibee@example.com")

	assert.Equal(t, "aibee006", u.Name)

}

func TestDeleteUserWithEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	api := NewApi()
	tx := api.Db.DB.Begin()
	api.Db.DB = tx
	defer tx.Rollback()

	r.DELETE("/users/delete/:email", api.DeleteUserWithEmail)

	req, _ := http.NewRequest(http.MethodDelete, "/users/delete/test@example.commm", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	u := api.Db.GetUserByEmail("test@example.commm")
	assert.Nil(t, u)

}
