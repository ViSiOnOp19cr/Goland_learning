package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", userhandler)
	fmt.Println("Server starting on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

type User struct {
	ID    int    `json:"ID"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var users []User
var nextID = 1

func userhandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		getalluser(w)
	case "POST":
		Createuser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func Createuser(w http.ResponseWriter, r *http.Request) {
	//i should send this to db ideally.
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	user.ID = nextID
	nextID++
	users = append(users, user)
	json.NewEncoder(w).Encode(user)
}

func getalluser(w http.ResponseWriter) {
	json.NewEncoder(w).Encode(users)
	fmt.Println(users)
}

//go can be used to create high performace application.
// go is a fast statically typed, complied language knows for simplicity and efficiency.
// go is mainly used for web based application (server side), cloud-native development.
//why use go???
// go has fast run time and compile time.
// go supports concurrency.
// go has memory management.
