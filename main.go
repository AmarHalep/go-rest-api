package main

import (
	"fmt"
	"net/http"
	"sync"
)

type User struct {
	Username string
	Password string
}

var (
	users = make(map[string]User)
	mu    sync.Mutex
)

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "ParseForm error", http.StatusBadRequest)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	if username == "" || password == "" {
		http.Error(w, "Username and password required", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()
	if _, exists := users[username]; exists {
		fmt.Fprintf(w, "Username %s already taken", username)
		return
	}
	users[username] = User{Username: username, Password: password}
	fmt.Fprintf(w, "User %s registered successfully", username)

}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "ParseForm error", http.StatusBadRequest)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	if username == "" || password == "" {
		http.Error(w, "Username and password required", http.StatusBadRequest)
		return
	}

	mu.Lock()
	user, exists := users[username]
	mu.Unlock()
	if !exists || user.Password != password {
		fmt.Fprintf(w, "Invalid username or password")
		return
	}
	fmt.Fprintf(w, "Welcome, %s!", username)
}

func showRegisterForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "register.html")

}

func showLoginForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "login.html")
}

// func mainPage(w http.ResponseWriter, r *http.Request) {

// }

func main() {
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/register-form", showRegisterForm)
	http.HandleFunc("/login-form", showLoginForm)
	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
