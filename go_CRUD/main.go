package main

import (
	"encoding/json"
	"log"
	"net/http"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID       uint `gorm:"primaryKey"`
	Name     string
	Email    string `gorm:"unique"`
	Age      int
	Verified bool
}

var db *gorm.DB

func connectDB() {
	dsn := "host=localhost user=postgres password=yourpassword dbname=gormdb port=5432 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatal("Migration failed:", err)
	}
}

func userhandler(w http.ResponseWriter, r *http.Request) {

	//first to get the json.
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusBadRequest)
		return
	}
	//send reqeust to database.
	result := db.Create(&user)
	if result.Error != nil {
		http.Error(w, "wrong in schema", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func getuserhandler(w http.ResponseWriter, r *http.Request) {
	var users []User
	result := db.Find(&users)

	// Handle DB error
	if result.Error != nil {
		http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(users)
}
func updatehandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusBadRequest)
		return
	}
	result := db.Model(&user).Where("id = ?", user.ID).Updates(user)
	if result.Error != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
func deletehandler(w http.ResponseWriter, r *http.Request){
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusBadRequest)
		return
	}
	result := db.Delete(&user)
	if result.Error != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}
}
func main() {
	connectDB()
	http.HandleFunc("/", userhandler)
	http.HandleFunc("/users", getuserhandler)
	http.HandleFunc("/update", updatehandler)
	http.HandleFunc("/delete", deletehandler)
	http.ListenAndServe(":8080", nil)
}
