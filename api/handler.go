package api

import (
	"encoding/json"
	"net/http"
	"net/mail"
	"strings"

	"github.com/gorilla/mux"
	"github.com/zohaibsoomro/users_api_golang/models"
	"github.com/zohaibsoomro/users_api_golang/utils"
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

func (api *Api) RegisterHandlers() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", helloWorld)
	router.HandleFunc("/users", getAllUsers).Methods(http.MethodGet)
	router.HandleFunc("/users/add", addUser).Methods(http.MethodPost)
	router.HandleFunc("/users/{email}", getUserWithEmail).Methods(http.MethodGet)
	router.HandleFunc("/users/update/{email}", updateUserWithEmail).Methods(http.MethodPut)
	router.HandleFunc("/users/delete/{email}", deleteUserWithEmail).Methods(http.MethodDelete)
	userDb.LoadAllUsers()
	return router
}

func helloWorld(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte(`<h1>Welcome Back!</h1>`))
}

func getAllUsers(w http.ResponseWriter, req *http.Request) {
	users := userDb.GetAllUsers()

	writeJSON(w, http.StatusOK, users)

}
func getUserWithEmail(w http.ResponseWriter, req *http.Request) {

	// Extract email from path

	email := strings.ToLower(strings.TrimSpace(mux.Vars(req)["email"]))

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

	email := strings.TrimSpace(mux.Vars(req)["email"])

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

	if _, err := mail.ParseAddress(user.Email); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid email in body: "+strings.TrimPrefix(err.Error(), "mail: "))
		return
	}

	if err := userDb.UpdateUserByEmail(email, &user); err != nil {
		writeError(w, http.StatusInternalServerError, "Unable to update user: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "User updated successfully!"})

}

func deleteUserWithEmail(w http.ResponseWriter, req *http.Request) {

	email := strings.TrimSpace(mux.Vars(req)["email"])

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

func addUser(w http.ResponseWriter, req *http.Request) {

	var u models.User
	if err := json.NewDecoder(req.Body).Decode(&u); err != nil {
		writeError(w, http.StatusBadRequest, "Unable to parse body: "+err.Error())
		return
	}
	u.Email = strings.TrimSpace(u.Email)
	if _, err := mail.ParseAddress(u.Email); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid email entered: "+err.Error())
		return
	}

	err := userDb.AddUser(u)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Unable to add user: "+err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "user added successfully."})
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
