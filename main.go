package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var db *sql.DB

func connectDB() {
	var err error
	db, err = sql.Open("mysql", "sql12780027:Sfy1gl2QWX@tcp(sql12.freesqldatabase.com:3306)/sql12780027")
	if err != nil {
		log.Fatal(err)
	}
}
func helloHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"message":    "hellow",
		"serverTime": time.Now().Format(time.RFC3339),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
func getUsers(w http.ResponseWriter, r *http.Request) {
	rows, _ := db.Query("SELECT id, name, email FROM users")
	var users []User
	for rows.Next() {
		var u User
		rows.Scan(&u.ID, &u.Name, &u.Email)
		users = append(users, u)
	}
	json.NewEncoder(w).Encode(users)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var u User
	json.NewDecoder(r.Body).Decode(&u)
	_, err := db.Exec("INSERT INTO users (name, email) VALUES (?, ?)", u.Name, u.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	_, err := db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	connectDB()
	defer db.Close()

	http.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getUsers(w, r)
		case "POST":
			createUser(w, r)
		case "DELETE":
			deleteUser(w, r)
		}
	})
	// ✅ Register the hello handler
	http.HandleFunc("/api/hello", helloHandler)
	fmt.Println("Go API running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
