package api

import (
	"encoding/json"
	"net/http"
	"net/mail"
	"strings"

	"github.com/zohaibsoomro/golangpractice/models"
	"github.com/zohaibsoomro/golangpractice/utils"
)

var userDb = utils.NewUsersDB()

type Api struct {
	Address string
}

func NewApi(address string) *Api {
	return &Api{
		Address: address,
	}
}

func (api *Api) RegisterHandlers() {
	http.HandleFunc("/", helloWorld)
	http.HandleFunc("/users", getAllUsers)
	http.HandleFunc("/users/", getUserWithEmail)
	http.HandleFunc("/users/update/", updateUserWithEmail)
	http.HandleFunc("/users/delete/", deleteUserWithEmail)
	userDb.LoadAllUsers()
}

func helloWorld(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte(`<h1>Welcome Back!</h1>`))
}

func getAllUsers(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		users := userDb.GetAllUsers()

		writeJSON(w, http.StatusOK, users)
	} else {
		writeError(w, http.StatusMethodNotAllowed, "Invalid request method.")
	}
}
func getUserWithEmail(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Invalid request method.")
		return
	}

	// Extract email from path
	email := strings.TrimPrefix(req.URL.Path, "/users/")
	email = strings.ToLower(strings.TrimSpace(email))

	if email == "" {
		writeError(w, http.StatusBadRequest, "Email parameter is required.")
		return
	}

	if _, err := mail.ParseAddress(email); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid email entered: "+strings.TrimPrefix(err.Error(), "mail: "))
		return
	}

	// Lookup user
	_, user := userDb.GetUserByEmail(email)
	if user == nil {
		writeError(w, http.StatusNotFound, "User not found.")
		return
	}

	writeJSON(w, http.StatusOK, user)
}

func updateUserWithEmail(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPut {
		writeError(w, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	email := strings.TrimPrefix(req.URL.Path, "/users/update/")
	email = strings.TrimSpace(email)

	if email == "" {
		writeError(w, http.StatusBadRequest, "email param is missing")
		return
	}

	if _, err := mail.ParseAddress(email); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid email entered: "+strings.TrimPrefix(err.Error(), "mail: "))
		return
	}

	var user models.User
	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		writeError(w, http.StatusBadRequest, "Unable to Decode Request Body")
		return
	}

	if err := userDb.UpdateUserByEmail(email, &user); err != nil {
		writeError(w, http.StatusInternalServerError, "Unable to update user: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, "User updated successfully!")

}

func deleteUserWithEmail(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodDelete {
		writeError(w, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}
	email := strings.TrimPrefix(req.URL.Path, "/users/delete/")
	email = strings.TrimSpace(email)

	if _, err := mail.ParseAddress(email); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid email entered: "+err.Error())
		return
	}

	err := userDb.DeleteUserByEmail(email)
	if err != nil {
		writeError(w, http.StatusNotAcceptable, "Error while deleting user: "+err.Error())
		return
	}
	writeJSON(w, http.StatusOK, "User deleted successfully!")
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, `{"error":"failed to encode json"}`, http.StatusInternalServerError)
		}
	}
}

// WriteError writes an error response in JSON.
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
