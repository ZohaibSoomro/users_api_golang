# 👥 Users API in Golang  

A simple **RESTful API** built with **Golang** that manages user records stored in a JSON file.  
This project is designed for learning purposes to practice Go, APIs, and project structure.  

---

## 🚀 Features  
- List all users 📋  
- Get user by email 🔍  
- Update user details ✏️  
- Delete user 🗑️  
- Input validation (including email format ✅)  

---

## 📂 Project Structure  
```bash
.
├── api/          # API layer (handlers & routing)
├── models/       # Data models (User struct)
├── utils/        # Utility functions (DB operations, JSON handling)
├── users.json    # Mock database (users stored here)
├── main.go       # Application entry point
└── go.mod        # Go module file
````

---

## ⚡ API Endpoints

### 1. Get All Users

```http
GET /users
```

### 2. Get User by Email

```http
GET /users/{email}
```

### 3. Update User by Email

```http
PUT /users/update/{email}
Content-Type: application/json

{
  "name": "New Name",
  "email": "new@email.com"
}
```

### 4. Delete User by Email

```http
DELETE /users/delete/{email}
```

---

## 🛠️ Installation & Setup

1. Clone this repo

```bash
git clone https://github.com/zohaibsoomro/users_api_golang.git
cd users_api_golang
```

2. Install dependencies

```bash
go mod tidy
```

3. Run the API server

```bash
go run main.go
```

4. Server will start at:

```
http://localhost:8080
```

---

## 🧪 Testing the API

You can use **Postman** or **VS Code REST Client extension** to test endpoints.

Example request with `curl`:

```bash
curl http://localhost:8080/users
```

---

## 📦 Example users.json

```json
[
  {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com"
  },
  {
    "id": 2,
    "name": "Jane Doe",
    "email": "jane@example.com"
  }
]
```

---

## 🤝 Contributing

Pull requests are welcome! Feel free to fork the repo and submit a PR.

